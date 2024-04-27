# Description

The purpose of this project is to track the specification for the Kobold General Purpose Messaging Protocol (KGPMP) throughout it's development. This spec will describe the transfer of data between `nodes` and `clients`. the KGPMP is meant for the the Kobold Distributed Messaging System (KDMS) under active development. Kobold uses a client/server model where messages are sent from a `client` to a `node` where they are then routed to other `clients`.

Although enforcing types across modern applications is very nice, it isn't essential. To me this means that if we wanted to enforce some sort of type safety, we could do it on the client side. A client could parse the `kgpmp` message and then validate the content of it as a type.

## Purpose

The purpose of KGPMP is to standardize message passing between `clients` and `nodes`. The protocol is meant to simplify and reduce overhead. It is meant to be easily parsed and turned into useable data structures within most modern programming languages. There are common patterns that keep coming up over and over again that can be solved simply and efficiently when communicating between distributed systems across multiple languages.

**Core patterns**

| Pattern           | Description                                                                    |
| ----------------- | ------------------------------------------------------------------------------ |
| Publish/Subscribe | `n` message producers sending messages to `x` consumers without a buffer       |
| Request/Reply     | a client sends a transaction to be acknowledged and returned by another client |

**Non Core patterns**

| Pattern       | Description                                                                         |
| ------------- | ----------------------------------------------------------------------------------- |
| Message Queue | `n` producers enqueue messages and `x` consumers dequeue messages at their own pace |
| Data Store    | ability to CRUD a key/value from multiple clients                                   |
| Stream        | send a stream of bytes over a persistent connection                                 |

## Terms

| term    | definition                                                                       |
| ------- | -------------------------------------------------------------------------------- |
| prefix  | 4 byte integer describing the number of bytes included in this message           |
| proto   | an encoded struct containing the required fields to route and handle the message |
| content | message content                                                                  |

## Protocol

I think that it's pretty important that the `node` no as little information as to the content is transporting between clients. By removing the encoding of the content from the `<proto>` payload, you can completely separate the actual message content.

## The Proto struct

```go
type KoboldMessage struct {
	// Client scoped unique indentifier for this message
	ID string `cbor:"id"`
	// Operation to be performed
	Op KoboldOperation `cbor:"op"`
	// the topic for which this message is to be forwarded to
	Topic string `cbor:"topic"`
	// data describing the message and it's origins
	Metadata KoboldMetadata `cbor:"metadata,omitempty"`
	// Used for request/reply to tie the request and reply together to the same client/connection
	TxID string `cbor:"tx_id,omitempty"`

	Content []byte `cbor:"content,omitempty"`
}
```

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
