package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/bahodge/kgpmp-prototype/pkg/protocol"
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

func main() {
	// Connect to the server
	// conn, err := net.Dial("tcp", "localhost:8080")
	// if err != nil {
	// 	fmt.Println("Error connecting:", err)
	// 	return
	// }
	// defer conn.Close()

	start := time.Now()
	m := protocol.KoboldMessage{
		// ID:    fmt.Sprintf("%d", i),
		Op:      protocol.Reply,
		Topic:   "test",
		TxID:    "",
		Content: []byte(""),
	}

	serializationStart := time.Now()

	chunks := [][]byte{}

	serializeCount := 0
	for i := 0; i < 1_000_000; i++ {
		m.ID = fmt.Sprintf("%d", i)

		s, err := protocol.Serialize(m)
		if err != nil {
			log.Fatal("could not serialize message", err)
		}

		chunks = append(chunks, s)

		serializeCount++
	}

	fmt.Println("chunks count", len(chunks))

	serializationEnd := time.Now()
	parsingStart := time.Now()
	parser := protocol.NewMessageParser()

	var rawMessages [][]byte
	parsedCount := 0
	for _, chunk := range chunks {
		messages, err := parser.Parse(chunk)
		if err != nil {
			log.Fatal("unable to parse data", err)
		}

		parsedCount += len(messages)
		rawMessages = append(rawMessages, messages...)
	}

	parsingEnd := time.Now()

	deserializationStart := time.Now()
	deserializedCount := 0
	deserializedMessages := []protocol.KoboldMessage{}
	for _, msg := range rawMessages {
		var deserializedMessage protocol.KoboldMessage
		err := protocol.Deserialize(msg, &deserializedMessage)
		if err != nil {
			log.Fatal("could not deserialize message", err)
		}

		deserializedMessages = append(deserializedMessages, deserializedMessage)
		deserializedCount++
	}

	deserializationEnd := time.Now()

	fmt.Println("serialization time", serializationEnd.Sub(serializationStart))
	fmt.Println("parsing time", parsingEnd.Sub(parsingStart))
	fmt.Println("deserialization time", deserializationEnd.Sub(deserializationStart))
	fmt.Println("serialized items", serializeCount)
	fmt.Println("parsedCount items", parsedCount)
	fmt.Println("deserialized items", deserializedCount)
	fmt.Println("total time", time.Since(start))

	// fmt.Printf("s %#v\n", string(s))

	// fmt.Printf("dec %#v\n", dec)

	// cbormsg := Message{
	// 	ID:            "",
	// 	MessageType:   1,
	// 	Topic:         "/hello/world",
	// 	TransactionID: "hereisthesenderid",
	// 	Content:       []byte("hello from message!The Parse function you provided parses a byte slice data into individual messages based on a length prefix encoding scheme. Each message is prefixed with a 4-byte length indicating the size of the message payload. The parsed result's structure is a slice of byte slices ([][]byte), where each inner byte slice represents a parsed message. Here's a breakdown of the parsed result's structure: The messages variable is a slice of byte slices ([][]byte), where each element represents a parsed message. Inside the loop, each complete message extracted from the buffer (message) is appended to the messages slice. After processing all complete messages in the buffer, the function returns the messages slice containing all parsed messages."),
	// 	Timestamp:     time.Now().Unix(),
	// }
	//
	// totalStart := time.Now()
	// rtts := []time.Duration{}
	// requests := 0
	// replies := 0
	// msgsOut := 0
	// bytesOut := 0
	//
	// msg, _ := cbor.Marshal(cbormsg)
	// // Convert int32 to bytes (big-endian)
	// bytesBigEndian := make([]byte, 4)
	// binary.BigEndian.PutUint32(bytesBigEndian, uint32(len(msg)))
	// msg = append(bytesBigEndian, msg...)
	// var sendBuffer []byte
	// maxBufferSize := 1024 * 1024 // 1mb
	//
	// iters := 10_000 * 1000
	// // In total we are sending 8.5GB over the network to the server.
	// for i := 0; i < iters; i++ {
	// 	start := time.Now()
	//
	// 	sendBuffer = append(sendBuffer, msg...)
	// 	msgsOut++
	//
	// 	if len(sendBuffer) > maxBufferSize || i+1 == iters {
	// 		requests++
	// 		if err := SendMessage(conn, sendBuffer); err != nil {
	// 			fmt.Println("Error sending message:", err)
	// 			fmt.Println("sent", i)
	// 			return
	// 		}
	//
	// 		// clear the send buffer
	// 		bytesOut += len(sendBuffer)
	// 		sendBuffer = sendBuffer[:0]
	// 	}
	//
	// 	rtts = append(rtts, time.Since(start))
	// }
	//
	// var sum int64
	// var mininum int64
	// var maximum int64
	// for _, dur := range rtts {
	// 	d := int64(dur)
	// 	if d > maximum {
	// 		maximum = d
	// 	}
	// 	if mininum == 0 {
	// 		mininum = d
	// 	} else if mininum > d {
	// 		mininum = d
	// 	}
	//
	// 	sum += d
	//
	// }
	//
	// avg := float64(sum) / float64(len(rtts))
	// fmt.Println("total roundtrips", len(rtts))
	// fmt.Println("total requests sent", requests)
	// fmt.Println("total bytes sent", bytesOut)
	// fmt.Println("total messages sent", msgsOut)
	// fmt.Println("total replies received", replies)
	// fmt.Println("total messages sent/received", replies+requests)
	// fmt.Println("min iter time", float64(mininum)/float64(time.Microsecond), "microseconds")
	// fmt.Println("max iter time", float64(maximum)/float64(time.Microsecond), "microseconds")
	// fmt.Println("average iter time", avg/float64(time.Microsecond), "microseconds")
	// fmt.Println("total command time", time.Since(totalStart))
	// fmt.Println("msgs per second", int64(float64(time.Second)/avg))
	//
	// // time.Sleep(time.Second)

}
