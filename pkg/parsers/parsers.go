package parsers

import "encoding/binary"

type MessageParser struct {
	buffer []byte
}

func NewMessageParser() *MessageParser {
	return &MessageParser{
		buffer: make([]byte, 0),
	}
}

// Extracts raw message bytes from a byte slice.
func (p *MessageParser) ParseLengthPrefixedMessage(data []byte) ([][]byte, error) {
	var messages [][]byte

	// Append incoming data to the buffer
	p.buffer = append(p.buffer, data...)

	// Parse complete messages from the buffer
	for len(p.buffer) >= 4 {
		// Read the length prefix
		messageLength := binary.BigEndian.Uint32(p.buffer[:4])

		// Check if the buffer contains the complete message
		if len(p.buffer) >= int(messageLength)+4 {
			// Slice the buffer to extract message content
			message := p.buffer[4 : 4+messageLength]

			// Append the message to the list of parsed messages
			messages = append(messages, message)

			// Remove the parsed message from the buffer
			p.buffer = p.buffer[4+messageLength:]
		} else {
			// Incomplete message in the buffer, wait for more data
			break
		}
	}

	return messages, nil
}
