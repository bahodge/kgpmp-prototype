# Description

The purpose of this project is to track the specification for the Kobold General Purpose Messaging Protocol (KGPMP) throughout it's development. This spec will describe the transfer of data between `nodes` and `clients`. the KGPMP is meant for the the Kobold Distributed Messaging System (KDMS) under active development. Kobold uses a client/server model where messages are sent from a `client` to a `node` where they are then routed to other `clients`.

## Purpose

The purpose of KGPMP is to standardize message passing between `clients` and `nodes`. The protocol is meant to simplify and reduce overhead. It is meant to be easily parsed and turned into useable data structures within most modern programming languages. There are common patterns that keep coming up over and over again that can be solved simply and efficiently when communicating between distributed systems across multiple languages. 

**Core patterns**

| Pattern | Description |
|--------|--------|
| Publish/Subscribe | `n` message producers sending messages to `x` consumers without a buffer |
| Request/Reply | a client sends a transaction to be acknowledged and returned by another client |

**Non Core patterns**

| Pattern | Description |
|--------|--------|
| Message Queue | `n` producers enqueue messages and `x` consumers dequeue messages at their own pace |
| Data Store | ability to CRUD a key/value from multiple clients | 
| Byte Stream | send `n` bytes to a specific client | 

## Terms

| term | definition |
|--------|--------|
| prefix | 4 byte integer describing the number of bytes included in this message |
| op | code describing the type of message |
| metadata | key value pairs separated by `;` |
| topic | an endpoint TODO improve this definition |
| content | message content | 

## Protocols

### Protocol 1

Each message is to be prefixed by a 4 byte integer describing the length of the message. This includes all delimiters included. 

```
<prefix><op>\r\n<metadata>\r\n<topic>\r\n<content>\r\n
```

Examples
```
51PUB\r\ntoken=secret;\r\n/hello/world\r\nhello\r\n
```

### Protocol 2

Each message is prefixed by a 4 byte integer describing the length of the message. The content of the body is encoded using cbor using a predefined message structure. Although there are benefits in the simplicity of this route, it does include a requirement to deserialize the content in order to determine where to route a message. This means that the server is doing a full deserialization of the message, looking at a subset of fields and then forwarding that messaging to the target client(s). It does make everything much simpler. Although the structure is much more concise, it does require a core dependency to an encoding dependency.

```
<prefix><content>
```

Examples
```
19binaryencodeddata
```

### Protocol 3

This protocol minimizes the required encoding/decoding to just focus on handling the meta data. The `op` and `topic` are plain text whereas the `metadata` translated and received in encoded using binary. The content is received as a byte slice meaning the actual content of the message is never required to be known or even understood by the server. This pattern follows more closely to the style of the `HTTP` protocol but with required encoding for metadata.

```
<prefix><op> <topic>\r\n<metadata>\r\ncontent\r\n
```

Examples
```
59PUB /hello/world\r\nbinaryencodeddata\r\nmysecretdata\r\n
```
