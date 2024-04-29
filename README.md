# Description

The purpose of this project is to track the specification for the Kobold General Purpose Messaging Protocol (KGPMP) throughout it's development. This spec will describe the transfer of data between `nodes` and `clients`. the KGPMP is meant for the the Kobold Distributed Messaging System (KDMS) under active development. Kobold is an in memory messaging system that uses a client/server model where messages are sent from a `client` to a `node` where they are then routed to other `clients`.

## Functionality

`kgpmp` is meant to support the following functionalities to facilitation communication between systems. These systems can be running on a single machine on an isolated network or a cluster of machines working together across the world.

- unidirectional `client` -> `node`
- bidirectional `client` <-> `node`
- unidirectional `node` -> `node`
- bidirectional `node` <-> `node`

## Purpose

The purpose of KGPMP is to standardize message passing between `clients` and `nodes`. The protocol is meant to simplify and reduce overhead. It is meant to be easily parsed and turned into useable data structures within most modern programming languages. There are common patterns that keep coming up over and over again that can be solved simply and efficiently when communicating between distributed systems across multiple languages.

**Core patterns**

| Pattern           | Description                                                                    |
| ----------------- | ------------------------------------------------------------------------------ |
| Request/Reply     | a client sends a transaction to be acknowledged and returned by another client |
| Publish/Subscribe | `n` message producers sending messages to `x` consumers without a buffer       |

**Non Core patterns**

| Pattern       | Description                                                                         |
| ------------- | ----------------------------------------------------------------------------------- |
| Message Queue | `n` producers enqueue messages and `x` consumers dequeue messages at their own pace |
| Data Store    | ability to CRUD a key/value from multiple clients                                   |
| Stream        | send a stream of bytes over a persistent connection                                 |

## Message Types

Messages that are sent from clients to nodes

| Message Type | Description                         |
| ------------ | ----------------------------------- |
| Publish      | (push) publish a message to a topic |
| Request      | send a request to a service topic.  |
| Reply        | send a reply to a request.          |

| Request Type | Description                                        |
| ------------ | -------------------------------------------------- |
| Forward      | forward the request to a client advertised service |
| Authenticate | authenticate the connection                        |
| Subscribe    | subscribe to a topic                               |
| Unsubscribe  | unsubscribe from a topic                           |
| Advertise    | advertise a service topic                          |
| Undvertise   | unadvertise a service topic                        |
| NodeInfo     | get information about a node                       |

| Topic Keywords | Description                               |
| -------------- | ----------------------------------------- |
| $node          | node the client is currently connected to |

```go
type Publish<T> struct {
    id: string
    topic: string
    metadata: Metadata // information about the client/connection
    message: T
}
---
// client creates a request
// client dispatches requests to node w/ tx_id
// node receives request
// node enqueues request in service request queue
// service dequeues next request
// service handles request
// service dispatches reply to node w/ tx_id
// node enqueues reply in client reply queue
// client handles reply
type Request<T> struct {
    id: string
    request_type: RequestType,
    topic: string
    tx_id: string
    metadata: Metadata // information about the client/connection
    message: T
}
---
type Reply<T> struct {
    id: string
    topic: string
    tx_id: string
    metadata: Metadata // information about the client/connection
    message: T
}
```

## Large Messages

A message must be able to pass from a client -> node -> `n` nodes -> `x` `clients`. So forwarding a message while maintaining it's integrity is essential. So in order to support larger payloads, we will have to cut messages into chunks. The simplest way I can think to do this is to add a 2 fields to the metadata struct. `part`, `total_parts`. The receiving client can simply track all parts of the message with `id` and reconstruct the larger message.

**Client based handlers**

So very simply, the receiving client can receive any chunk of a message and know the total number of parts of the message. If a part was not received in the correct amount of time, the client could ask for a retransmission of just that part. The receiver client would just need to retrace the message back to it's origin and request that specific part. Although I think the functionality of this is necessary, I think it is quite cumbersome and could be done better.

**Node based handlers**

Another easy way to get data verification is for clients to upload the entire message to the `node` and have the node validate that all parts are received. This puts the data validation part on the `node` and removes the complexity of retracing a message back to the originating client. It does put more demand on available disk or memory of the node.

## Terms

| term   | definition                                                                       |
| ------ | -------------------------------------------------------------------------------- |
| prefix | 4 byte integer describing the number of bytes included in this message           |
| proto  | an encoded struct containing the required fields to route and handle the message |

## Protocol

Each message is to be prefixed by a 4 byte integer describing the length of the message. This includes all delimiters included.

```
<prefix><proto>
```

Examples

```
# request
<prefix><encoded message>
```

## Encoding Benchmarks

TLDR; `cbor` seems to be the best starting point for encoding that I can come up with.

although my tests are shitty and primitive, they show that json is slow as fuck and that Cap'n Proto and CBOR are pretty similar. There are great benefits to using either.

This is all running on my local machine and therefore all tests are uncontrolled. This is all biased, I know. I'm not interested in being really thorough. There are some obvious optimizations that can happen. I'm sure I could write a faster or more efficient encoding system but I doubt it will be as robust as any of these others. `cbor` wins for now.

| Encoding    | Iterations | Serialization(ms) | Parsing(ms) | Deserialization(ms) | Total Time(ms) |
| ----------- | ---------- | ----------------- | ----------- | ------------------- | -------------- |
| Cap'n Proto | 1_000_000  | 1250.15           | 96.87       | 334.78              | 1681.90        |
| CBOR        | 1_000_000  | 725.97            | 86.38       | 672.57              | 1484.96        |
| Msgpack     | 1_000_000  | untested          | untested    | untested            | untested       |
| JSON        | 1_000_000  | 1095.81           | 122.90      | 1595.00             | 2813.74        |

```
goos: linux
goarch: amd64
pkg: github.com/bahodge/kgpmp-prototype
cpu: 13th Gen Intel(R) Core(TM) i9-13900K
BenchmarkRunCapn1
BenchmarkRunCapn1-24             	    2774	    418643 ns/op
BenchmarkRunCapn100
BenchmarkRunCapn100-24           	    2532	    514270 ns/op
BenchmarkRunCapn10000
BenchmarkRunCapn10000-24         	      69	  17829592 ns/op
BenchmarkRunCapn100000
BenchmarkRunCapn100000-24        	       6	 169682814 ns/op
BenchmarkRunCapn1000000
BenchmarkRunCapn1000000-24       	       1	1751991028 ns/op
BenchmarkRunCBOR1
BenchmarkRunCBOR1-24             	    3051	    373979 ns/op
BenchmarkRunCBOR100
BenchmarkRunCBOR100-24           	    1687	    598904 ns/op
BenchmarkRunCBOR10000
BenchmarkRunCBOR10000-24         	      72	  16247588 ns/op
BenchmarkRunCBOR100000
BenchmarkRunCBOR100000-24        	       6	 177246694 ns/op
BenchmarkRunCBOR1000000
BenchmarkRunCBOR1000000-24       	       1	1685516807 ns/op
BenchmarkRunMsgpack1
BenchmarkRunMsgpack1-24          	    4377	    401017 ns/op
BenchmarkRunMsgpack100
BenchmarkRunMsgpack100-24        	    1807	    676553 ns/op
BenchmarkRunMsgpack10000
BenchmarkRunMsgpack10000-24      	      63	  18400306 ns/op
BenchmarkRunMsgpack100000
BenchmarkRunMsgpack100000-24     	       7	 179656233 ns/op
BenchmarkRunMsgpack1000000
BenchmarkRunMsgpack1000000-24    	       1	1929324637 ns/op
BenchmarkRunJSON1
BenchmarkRunJSON1-24             	    3318	    457479 ns/op
BenchmarkRunJSON100
BenchmarkRunJSON100-24           	    1335	    763268 ns/op
BenchmarkRunJSON10000
BenchmarkRunJSON10000-24         	      44	  26974104 ns/op
BenchmarkRunJSON100000
BenchmarkRunJSON100000-24        	       5	 235920948 ns/op
BenchmarkRunJSON1000000
BenchmarkRunJSON1000000-24       	       1	2910120434 ns/op
```
