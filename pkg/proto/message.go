package proto

type KoboldOperation uint16

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
}

type KoboldMessage struct {
	ID            string          `cbor:"id"`
	Operation     KoboldOperation `cbor:"op"`
	Topic         string          `cbor:"topic"`
	Metadata      KoboldMetadata  `cbor:"metadata,omitempty"`
	Headers       map[string]any  `cbor:"headers,omitempty"`
	TransactionID string          `cbor:"tx_id,omitempty"`
}
