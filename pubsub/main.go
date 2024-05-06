package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/bahodge/kgpmp-prototype/pkg/protocol"
)

func handleRequest(conn net.Conn) {
	totalMessages := 0
	defer conn.Close()
	defer func() {
		fmt.Println("received total messages", totalMessages)

	}()
	// Create a new message parser for each client connection
	parser := protocol.NewMessageParser()

	// Buffer to store incoming data from the client
	chunk := make([]byte, 1024*1024)

	for {
		// Read data from the client
		n, err := conn.Read(chunk)
		if err != nil {
			if err == io.EOF {
				fmt.Println("client closed connection")
				return
			}
			fmt.Println("Error reading from client:", err)
			continue
		}

		// Parse complete messages from the received data
		messages, err := parser.Parse(chunk[:n])
		if err != nil {
			fmt.Println("Error parsing messages:", err)
			return
		}

		// Process the parsed messages
		for _, message := range messages {
			var dec protocol.Message
			err := protocol.DeserializeCBOR(message, &dec)
			if err != nil {
				fmt.Fprintln(os.Stdout, []any{"could not deserialize messagge", err}...)
				return
			}

			totalMessages++
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

var msgId int

func RunPub(addr string, topic string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal("could not reach addr", err)
	}
	defer conn.Close()

	totalStart := time.Now()
	rtts := []time.Duration{}
	requests := 0
	replies := 0
	msgsOut := 0
	bytesOut := 0

	for i := 0; i < 50_000; i++ {
		start := time.Now()
		m := protocol.Message{
			Id:          fmt.Sprintf("%d", i),
			MessageType: protocol.Publish,
			Topic:       topic,
			TxId:        "",
			Headers:     protocol.Headers{},
			Content:     []byte("hello world"),
			Errors:      []protocol.Error{},
			Timestamp:   time.Now().UnixMicro(),
		}

		s, err := protocol.SerializeCBOR(m)
		if err != nil {
			log.Fatal("could not serialize message", err)

		}

		// fmt.Println("s", s)

		// Note there there is no chunking, we are just sending 1 message at a time.
		// this does not perform optimally
		n, err := conn.Write(s)
		if err != nil {
			log.Fatal("could not write buf", err)
		}

		bytesOut += n
		rtts = append(rtts, time.Since(start))
	}

	var sum int64
	var mininum int64
	var maximum int64
	for _, dur := range rtts {
		d := int64(dur)
		if d > maximum {
			maximum = d
		}
		if mininum == 0 {
			mininum = d
		} else if mininum > d {
			mininum = d
		}

		sum += d

	}

	avg := float64(sum) / float64(len(rtts))
	fmt.Println("total roundtrips", len(rtts))
	fmt.Println("total requests sent", requests)
	fmt.Println("total bytes sent", bytesOut)
	fmt.Println("total messages sent", msgsOut)
	fmt.Println("total replies received", replies)
	fmt.Println("total messages sent/received", replies+requests)
	fmt.Println("min iter time", float64(mininum)/float64(time.Microsecond), "microseconds")
	fmt.Println("max iter time", float64(maximum)/float64(time.Microsecond), "microseconds")
	fmt.Println("average iter time", avg/float64(time.Microsecond), "microseconds")
	fmt.Println("total command time", time.Since(totalStart))
	fmt.Println("msgs per second", int64(float64(time.Second)/avg))
}

func RunSub(addr string, topic string, clientId string) {
	// TODO
}

func main() {
	if len(os.Args) > 2 && os.Args[1] == "node" {
		RunNode(os.Args[2])
		os.Exit(0)
	}
	if len(os.Args) > 3 && os.Args[1] == "pub" {
		RunPub(os.Args[2], os.Args[3])
		os.Exit(0)
	}
	if len(os.Args) > 3 && os.Args[1] == "sub" {
		RunSub(os.Args[2], os.Args[3], os.Args[4])
		os.Exit(0)
	}
	fmt.Fprintf(os.Stderr, "Usage: pubsub server|client <URL> <TOPIC> <CONTENT>\n")
	os.Exit(1)
}
