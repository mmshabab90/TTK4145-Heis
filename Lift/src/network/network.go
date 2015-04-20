package network

import (
	"../defs"
	"encoding/json"
	"fmt"
	"time"
)

// --------------- PUBLIC: ---------------

type UdpConnection struct { //should this be in udp.go? or in poller.go?
	Addr  string
	Timer *time.Timer
}

var ReceiveChan = make(chan udpMessage, 10) //this is now buffered with 10 slots, does this lead to fuckup?

func Init() {
	const localListenPort = 20001
	const broadcastListenPort = 20058
	const messageSize = 1024
	err := Udp_init(localListenPort, broadcastListenPort, messageSize, sendChan, ReceiveChan)

	if err != nil {
		fmt.Print("err = %s \n", err)
	}

	go aliveSpammer()
	go pollMessages()
}

func pollMessages() {
	var msg defs.Message
	for {
		msg = <- defs.MessageChan
		//PrintMessage(msg)
		if jsonMessage, err := json.Marshal(msg); err != nil { 
			fmt.Printf("json.Unmarshal error: %s\n", err)
		} else {
			sendChan <- udpMessage{raddr: "broadcast", data: jsonMessage, length: len(jsonMessage)}
		}
		time.Sleep(time.Millisecond)
	}
}

func ParseMessage(udpMessage udpMessage) defs.Message {
	var message defs.Message
	/*fmt.Println("in parsemessage")
	PrintMessage(message)
	PrintMessage(message)*/
	if err := json.Unmarshal(udpMessage.data[:udpMessage.length], &message); err != nil {
		fmt.Printf("json.Unmarshal error: %s\n", err)
	}

	message.Addr = udpMessage.raddr
	return message
}

// --------------- PRIVATE: ---------------

var sendChan = make(chan udpMessage)

func aliveSpammer() {
	const spamInterval = 100 * time.Millisecond
	message := defs.Message{Kind: defs.Alive}
	for {
		defs.MessageChan <- message
		time.Sleep(spamInterval)
	}
}

func PrintMessage(msg defs.Message) {
	fmt.Printf("\n-----Message start-----\n")
	switch msg.Kind {
	case defs.Alive:
		fmt.Println("I'm alive")
	case defs.NewOrder:
		fmt.Println("New order:")
		fmt.Printf("Floor: %d\n", msg.Floor)
		fmt.Printf("Button: %d\n", msg.Button)
	case defs.CompleteOrder:
		fmt.Println("Complete order:")
		fmt.Printf("Floor: %d\n", msg.Floor)
		fmt.Printf("Button: %d\n", msg.Button)
	case defs.Cost:
		fmt.Printf("Cost: %d\n", msg.Cost)
	default:
		fmt.Println("Invalid message type!\n")
	}
	fmt.Println("-----Message end-------\n")
}
