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
	Token        string `cbor:"token,omitempty"`
}

// TODO: I think that most of the metadata can be put in the main message and omitted as neccessary

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
	// Used for proxy connection where bytes are sent to and fro
	SenderConnID   string `cbor:"sender_conn_id,omitempty"`
	ReceiverConnID string `cbor:"receiver_conn_id,omitempty"`
}
