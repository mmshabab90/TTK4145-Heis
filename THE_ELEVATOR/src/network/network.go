package network

import (
	"../defs"
	"encoding/json"
	"fmt"
	"time"
)

// --------------- PUBLIC: ---------------

const (
	Alive int = iota
	NewOrder
	CompleteOrder
	Cost
)

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
}

func Send(message *defs.Message) { //now takes a pointer, does it still work over the network?
	//printMessage(*message)
	jsonMessage, err := json.Marshal(message) //is json good? can it take a pointer?
	if err != nil {
		// worry
	} else {
		sendChan <- udpMessage{raddr: "broadcast", data: jsonMessage, length: len(jsonMessage)}
	}
}

func ParseMessage(udpMessage udpMessage) defs.Message {
	var message defs.Message
	err := json.Unmarshal(udpMessage.data, &message)
	if err != nil {
		// handle
	}
	message.Addr = udpMessage.raddr
	return message
}

// --------------- PRIVATE: ---------------

var sendChan = make(chan udpMessage)

func aliveSpammer() {
	const spamInterval = 500 * time.Millisecond
	message := &defs.Message{Kind: Alive}
	for {
		Send(message)
		time.Sleep(spamInterval)
	}
}

func printMessage(msg defs.Message) {
	fmt.Println("Message")
	fmt.Println("---------------------------")
	switch msg.Kind {
	case Alive:
		fmt.Println("I'm alive\n")
	case NewOrder:
		fmt.Println("New order:")
		fmt.Printf("Floor: %d\n", msg.Floor)
		fmt.Printf("Button: %d\n\n", msg.Button)
	case CompleteOrder:
		fmt.Println("Complete order:")
		fmt.Printf("Floor: %d\n", msg.Floor)
		fmt.Printf("Button: %d\n\n", msg.Button)
	case Cost:
		fmt.Printf("Cost: %d\n\n", msg.Cost)
	default:
		fmt.Println("Invalid message type!\n")
	}
}
