package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/bahodge/kgpmp-prototype/pkg/protocol"
)

func handleRequest(conn net.Conn) {
	defer conn.Close()
	// Create a new message parser for each client connection
	parser := protocol.NewMessageParser()

	// Buffer to store incoming data from the client
	chunk := make([]byte, 1024*1024)

	for {
		// Read data from the client
		n, err := conn.Read(chunk)
		if err != nil {
			fmt.Println("Error reading from client:", err)
			return
		}

		// Parse complete messages from the received data
		messages, err := parser.Parse(chunk[:n])
		if err != nil {
			fmt.Println("Error parsing messages:", err)
			return
		}

		// Process the parsed messages
		for _, message := range messages {
			var dec protocol.KoboldMessage
			err := protocol.DeserializeCBOR(message, &dec)
			if err != nil {
				fmt.Fprintln(os.Stdout, []any{"could not deserialize messagge", err}...)
				return
			}

		}
	}
}

func RunNode(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Printf("node listening on %s\n", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("new connection")

		go handleRequest(conn)
	}

}

func RunPub(url string, topic string, payload string) {

}

func main() {
	if len(os.Args) > 2 && os.Args[1] == "node" {
		RunNode(os.Args[2])
		os.Exit(0)
	}
	if len(os.Args) > 3 && os.Args[1] == "pub" {
		RunPub(os.Args[2], os.Args[3], os.Args[4])
		os.Exit(0)
	}
	fmt.Fprintf(os.Stderr, "Usage: pubsub server|client <URL> <TOPIC> <CONTENT>\n")
	os.Exit(1)
}
