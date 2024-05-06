package protocol

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/fxamacker/cbor/v2"
	"github.com/vmihailenco/msgpack/v5"
)

const MAX_MSG_SIZE = 1024 * 1024

type MessageType uint8

const (
	Unsupported MessageType = iota
	Request                 // Send a request to a service topic
	Reply                   // Send a reply from a service topic
	Advertise               // Initiate a service topic
	Unadvertise             // Close a service topic
	Publish                 // Publish a message to a topic
	Subscribe               // Subscribe to messages on a topic
	Unsubscribe             // Unsubscribe from a topic
)

type ErrorCode uint8

const (
	CodeNoError ErrorCode = iota
	CodeServiceTopicNotFound
	CodeCouldNotHandleMessage
	CodeMalformedMessage
	CodeUnauthorized
)

var (
	ErrorServiceTopicNotFound  = errors.New("service topic not found")
	ErrorCouldNotHandleMessage = errors.New("could not handle message")
	ErrorMalformedMessage      = errors.New("malformed message")
	ErrorUnauthorized          = errors.New("unauthorized")
)

// type Message struct {
// 	ID       string          `msgpack:"id" json:"id" cbor:"id"`
// 	Op       KoboldOperation `msgpack:"op" json:"op" cbor:"op"`
// 	Topic    string          `msgpack:"topic" json:"topic" cbor:"topic"`
// 	Metadata KoboldMetadata  `msgpack:"Metadata" json:"metadata" cbor:"metadata,omitempty"`
// 	TxID     string          `msgpack:"tx_id" json:"tx_id" cbor:"tx_id,omitempty"`
// 	Content  []byte          `msgpack:"content" json:"content" cbor:"content,omitempty"`
// }

type Message struct {
	Id          string      `cbor:"id"`
	MessageType MessageType `cbor:"message_type"`
	Topic       string      `cbor:"topic"`
	TxId        string      `cbor:"tx_id,omitempty"`
	Headers     Headers     `cbor:"headers,omitempty"`
	Content     []byte      `cbor:"content,omitempty"`
	Errors      []Error     `cbor:"errors,omitempty"`
	Timestamp   int64       `cbor:"timestamp,omitempty"`
}

type Error struct {
	Message string    `cbor:"message,omitempty"`
	Code    ErrorCode `cbor:"code,omitempty"`
}

type Headers struct {
	ClientId  string `cbor:"client_id,omitempty"`
	ConnId    string `cbor:"conn_id,omitempty"`
	AuthToken string `cbor:"auth_token,omitempty"`
}

func PrefixWithLength(payload []byte) ([]byte, error) {
	var buf bytes.Buffer
	// Check if payload exceeds maximum message size
	if len(payload) > MAX_MSG_SIZE {
		return nil, errors.New("message is too large")
	}

	// Write payload length prefix to the buffer
	bytesBigEndian := make([]byte, 4)
	binary.BigEndian.PutUint32(bytesBigEndian, uint32(len(payload)))
	_, err := buf.Write(bytesBigEndian)
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(payload)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil

}

func SerializeCBOR(msg Message) ([]byte, error) {
	var payload []byte
	var err error

	payload, err = cbor.Marshal(msg)
	if err != nil {
		return nil, err
	}

	// Check if payload exceeds maximum message size
	if len(payload) > MAX_MSG_SIZE {
		return nil, errors.New("message is too large")
	}

	return PrefixWithLength(payload)
}

func DeserializeCBOR(data []byte, m *Message) error {
	// in this case we assume that we already have chopped off the first 4 bytes
	// as part of the parsing step. we now just need to Unmarshal cbor
	return cbor.Unmarshal(data, m)
}

func SerializeMsgpack(msg Message) ([]byte, error) {
	var payload []byte
	var err error

	payload, err = msgpack.Marshal(msg)
	if err != nil {
		return nil, err
	}

	// Check if payload exceeds maximum message size
	if len(payload) > MAX_MSG_SIZE {
		return nil, errors.New("message is too large")
	}

	return PrefixWithLength(payload)
}

func DeserializeMsgpack(data []byte, m *Message) error {
	// in this case we assume that we already have chopped off the first 4 bytes
	// as part of the parsing step. we now just need to Unmarshal cbor
	return msgpack.Unmarshal(data, m)
}

func SerializeJSON(msg Message) ([]byte, error) {
	var payload []byte
	var err error

	payload, err = json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	// Check if payload exceeds maximum message size
	if len(payload) > MAX_MSG_SIZE {
		return nil, errors.New("message is too large")
	}

	// Write payload length prefix to the buffer
	return PrefixWithLength(payload)
}

func DeserializeJSON(data []byte, m *Message) error {
	// in this case we assume that we already have chopped off the first 4 bytes
	// as part of the parsing step. we now just need to Unmarshal json
	return json.Unmarshal(data, m)
}
