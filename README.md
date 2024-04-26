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
| Proxy             | a client can connect to another client using the node as a proxy               |

**Non Core patterns**

| Pattern       | Description                                                                         |
| ------------- | ----------------------------------------------------------------------------------- |
| Message Queue | `n` producers enqueue messages and `x` consumers dequeue messages at their own pace |
| Data Store    | ability to CRUD a key/value from multiple clients                                   |
| Byte Stream   | send a stream of bytes over a persistent connection                                 |

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
type KoboldProto struct {
        // Client scoped unique indentifier for this message
	ID string `cbor:"id"`
	// Operation to be performed
	Operation KoboldOperation `cbor:"op"`
	// the topic for which this message is to be forwarded to
	Topic string `cbor:"topic"`
	// data describing the message and it's origins
	Metadata KoboldMetadata `cbor:"metadata,omitempty"`
	// globally unique id used to tie the request and reply together to the same client/connection
	TransactionID string `cbor:"tx_id,omitempty"`
}
```

Each message is to be prefixed by a 4 byte integer describing the length of the message. This includes all delimiters included.

```
<prefix><proto>\r\n<content\r\n
```

Examples

```
# request
<prefix><BINARYENCODEDPROTO>\r\nmycoolcontent\r\n
```
