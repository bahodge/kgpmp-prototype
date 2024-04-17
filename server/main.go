package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

type Message struct {
	ID            string `cbor:"id"`
	MessageType   uint32 `cbor:"message_type"`
	Topic         string `cbor:"topic"`
	TransactionID string `cbor:"transaction_id"`
	Content       []byte `cbor:"content"`
	Timestamp     int64  `cbor:"timestamp"`
}

type MessageParser struct {
	buffer []byte
}

func NewMessageParser() *MessageParser {
	return &MessageParser{
		buffer: make([]byte, 0),
	}
}

func (p *MessageParser) Parse(data []byte) ([][]byte, error) {
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

func handleClient(conn net.Conn) {
	msgCount := 0
	bytesCount := 0
	start := time.Now()
	defer func() {
		fmt.Println("message count", msgCount)
		fmt.Println("bytes count", bytesCount)
		fmt.Println("time to receive messages", time.Since(start))

		conn.Close()

	}()

	// Create a new message parser for each client connection
	parser := NewMessageParser()

	// Buffer to store incoming data from the client
	buffer := make([]byte, 1024*1024)

	for {
		// Read data from the client
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from client:", err)
			return
		}

		// Parse complete messages from the received data
		messages, err := parser.Parse(buffer[:n])
		if err != nil {
			fmt.Println("Error parsing messages:", err)
			return
		}
		// Process the parsed messages
		for _, message := range messages {
			// var msg Message
			// err := cbor.Unmarshal(message, &msg)
			// if err != nil {
			// 	log.Fatal("could not parse message", msg.ID)
			// }

			msgCount += 1
			bytesCount += len(message)
		}
	}
}

func main() {
	// Start TCP server
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server is listening on port 8080")

	// Accept incoming connections
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Handle client connections concurrently
		go handleClient(conn)
	}
}
