package main

import (
	"fmt"
	"github.com/bahodge/kgpmp-prototype/pkg/protocol"
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

	// Create a new message p for each client connection
	p := protocol.NewMessageParser()

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
		messages, err := p.Parse(buffer[:n])
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
