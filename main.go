package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"capnproto.org/go/capnp/v3"
	"github.com/bahodge/kgpmp-prototype/pkg/protocol"
	"github.com/bahodge/kgpmp-prototype/protos"
)

func SendMessage(conn net.Conn, message []byte) error {
	n, err := conn.Write(message)
	if err != nil {
		fmt.Println("unable to send message", err)
		return err
	}

	if len(message) != n {
		fmt.Println("did not send full payload", n)
		return errors.New("did not send full payload")
	}

	return nil
}

func RunCapn(iterations int) {
	var sendBuf bytes.Buffer
	sendBuf.Grow(1024 * 1024)

	for i := 0; i < iterations; i++ {
		arena := capnp.SingleSegment(nil)
		msg, seg, err := capnp.NewMessage(arena)
		if err != nil {
			log.Fatal("could not create new message", err)
		}

		kmsg, err := protos.NewRootKoboldMessage(seg)
		if err != nil {
			log.Fatal("could not spawn new root message", err)
		}

		err = kmsg.SetId(fmt.Sprintf("%d", i))
		err = kmsg.SetTopic("/hello/world")
		err = kmsg.SetTxId(fmt.Sprintf("sometxid%d", i))
		err = kmsg.SetContent([]byte("hello there"))
		if err != nil {
			log.Fatal("could not set err", err)
		}

		payload, err := msg.Marshal()
		if err != nil {
			log.Fatal("could not marshal message!", err)
		}

		b, err := protocol.PrefixWithLength(payload)
		if err != nil {
			log.Fatal("could not prefix", err)
		}

		_, err = sendBuf.Write(b)
		if err != nil {
			log.Fatal("could not write bytes to buffer")
		}

	}

	// parsingStart := time.Now()
	var rawMessages [][]byte
	parser := protocol.NewMessageParser()

	for {
		// Read a chunk of data from the buffer
		chunk := make([]byte, 1024*1024)
		n, err := sendBuf.Read(chunk)
		if err != nil {
			if err == io.EOF {
				// End of buffer reached, exit the loop
				break
			}
			log.Fatal("could not read", err)
		}

		// Process the chunk only if it contains data
		if n > 0 {
			// Parse the chunk to extract complete messages
			messages, err := parser.Parse(chunk[:n]) // Pass only the portion of the chunk that contains valid data
			if err != nil {
				log.Fatal("unable to parse data", err)
			}

			// Update counters and append parsed messages
			rawMessages = append(rawMessages, messages...)
		}
	}

	deserializedMessages := []protos.KoboldMessage{}
	for _, payload := range rawMessages {
		msg, err := capnp.Unmarshal(payload)
		if err != nil {
			log.Fatal("could not unmarshal payload", err)
		}

		kmsg, err := protos.ReadRootKoboldMessage(msg)
		if err != nil {
			log.Fatal("could not read root message", err)
		}

		deserializedMessages = append(deserializedMessages, kmsg)
	}
}

func RunCBOR(iterations int) {
	var sendBuf bytes.Buffer
	sendBuf.Grow(1024 * 1024)

	serializeCount := 0
	for i := 0; i < iterations; i++ {
		m := protocol.KoboldMessage{
			ID:      fmt.Sprintf("%d", i),
			Op:      protocol.Reply,
			Topic:   "/hello/world",
			TxID:    fmt.Sprintf("sometxid - %d", i),
			Content: []byte("hello world"),
		}

		s, err := protocol.SerializeCBOR(m)
		if err != nil {
			log.Fatal("could not serialize message", err)
		}

		_, err = sendBuf.Write(s)
		if err != nil {
			log.Fatal("could not write to buffer")
		}

		serializeCount++
	}

	parser := protocol.NewMessageParser()

	var rawMessages [][]byte

	for {
		// Read a chunk of data from the buffer
		chunk := make([]byte, 1024*1024)
		n, err := sendBuf.Read(chunk)
		if err != nil {
			if err == io.EOF {
				// End of buffer reached, exit the loop
				break
			}
			log.Fatal("could not read", err)
		}

		// Process the chunk only if it contains data
		if n > 0 {
			// Parse the chunk to extract complete messages
			messages, err := parser.Parse(chunk[:n]) // Pass only the portion of the chunk that contains valid data
			if err != nil {
				log.Fatal("unable to parse data", err)
			}

			// Update counters and append parsed messages
			// parsedCount += len(messages)
			rawMessages = append(rawMessages, messages...)
		}
	}

	deserializedMessages := []protocol.KoboldMessage{}
	for _, msg := range rawMessages {
		var deserializedMessage protocol.KoboldMessage
		err := protocol.DeserializeCBOR(msg, &deserializedMessage)
		if err != nil {
			log.Fatal("could not deserialize message", err)
		}

		deserializedMessages = append(deserializedMessages, deserializedMessage)
	}
}

func RunJSON(iterations int) {
	var sendBuf bytes.Buffer
	sendBuf.Grow(1024 * 1024)

	for i := 0; i < iterations; i++ {
		m := protocol.KoboldMessage{
			ID:      fmt.Sprintf("%d", i),
			Op:      protocol.Reply,
			Topic:   "/hello/world",
			TxID:    fmt.Sprintf("sometxid - %d", i),
			Content: []byte("hello world"),
		}

		s, err := protocol.SerializeJSON(m)
		if err != nil {
			log.Fatal("could not serialize message", err)
		}

		_, err = sendBuf.Write(s)
		if err != nil {
			log.Fatal("could not write to buffer")
		}

	}

	parser := protocol.NewMessageParser()

	var rawMessages [][]byte

	for {
		// Read a chunk of data from the buffer
		chunk := make([]byte, 1024*1024)
		n, err := sendBuf.Read(chunk)
		if err != nil {
			if err == io.EOF {
				// End of buffer reached, exit the loop
				break
			}
			log.Fatal("could not read", err)
		}

		// Process the chunk only if it contains data
		if n > 0 {
			// Parse the chunk to extract complete messages
			messages, err := parser.Parse(chunk[:n]) // Pass only the portion of the chunk that contains valid data
			if err != nil {
				log.Fatal("unable to parse data", err)
			}

			// Update counters and append parsed messages
			rawMessages = append(rawMessages, messages...)
		}
	}

	deserializedMessages := []protocol.KoboldMessage{}
	for _, msg := range rawMessages {
		var deserializedMessage protocol.KoboldMessage
		err := protocol.DeserializeJSON(msg, &deserializedMessage)
		if err != nil {
			log.Fatal("could not deserialize message", err)
		}

		deserializedMessages = append(deserializedMessages, deserializedMessage)
	}

}

func main() {
	const ITERATIONS = 1_000_000
	RunCapn(ITERATIONS)
	RunCBOR(ITERATIONS)
	RunJSON(ITERATIONS)
}