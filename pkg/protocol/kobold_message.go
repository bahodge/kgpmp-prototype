package protocol

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/fxamacker/cbor/v2"
)

type KoboldOperation uint8

const MAX_MSG_SIZE = 1024 * 1024

const (
	Unsupported KoboldOperation = iota
	Subscribe
	Unsubscribe
	Publish
	Advertise
	Unadvertise
	Request
	Reply
)

type KoboldMetadata struct {
	ClientID     string `cbor:"client_id"`
	ConnectionID string `cbor:"conn_id"`
	Token        string `cbor:"token,omitempty"`
}

// TODO: I think that most of the metadata can be put in the main message and omitted as neccessary

type KoboldMessage struct {
	ID       string          `cbor:"id"`
	Op       KoboldOperation `cbor:"op"`
	Topic    string          `cbor:"topic"`
	Metadata KoboldMetadata  `cbor:"metadata,omitempty"`
	TxID     string          `cbor:"tx_id,omitempty"`
	Content  []byte          `cbor:"content,omitempty"`
}

func SerializeCBOR(msg KoboldMessage) ([]byte, error) {
	var payload []byte
	var err error
	var buf bytes.Buffer

	payload, err = cbor.Marshal(msg)
	if err != nil {
		return nil, err
	}

	// Check if payload exceeds maximum message size
	if len(payload) > MAX_MSG_SIZE {
		return nil, errors.New("message is too large")
	}

	// Write payload length prefix to the buffer
	bytesBigEndian := make([]byte, 4)
	binary.BigEndian.PutUint32(bytesBigEndian, uint32(len(payload)))
	_, err = buf.Write(bytesBigEndian)
	if err != nil {
		return nil, err
	}

	// Write payload to the buffer
	_, err = buf.Write(payload)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func DeserializeCBOR(data []byte, m *KoboldMessage) error {
	// in this case we assume that we already have chopped off the first 4 bytes
	// as part of the parsing step. we now just need to Unmarshal cbor
	return cbor.Unmarshal(data, m)
}

func SerializeJSON(msg KoboldMessage) ([]byte, error) {
	var payload []byte
	var err error
	var buf bytes.Buffer

	payload, err = json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	// Check if payload exceeds maximum message size
	if len(payload) > MAX_MSG_SIZE {
		return nil, errors.New("message is too large")
	}

	// Write payload length prefix to the buffer
	bytesBigEndian := make([]byte, 4)
	binary.BigEndian.PutUint32(bytesBigEndian, uint32(len(payload)))
	_, err = buf.Write(bytesBigEndian)
	if err != nil {
		return nil, err
	}

	// Write payload to the buffer
	_, err = buf.Write(payload)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func DeserializeJSON(data []byte, m *KoboldMessage) error {
	// in this case we assume that we already have chopped off the first 4 bytes
	// as part of the parsing step. we now just need to Unmarshal cbor
	return json.Unmarshal(data, m)
}
