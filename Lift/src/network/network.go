package network

import (
	"../defs"
	"encoding/json"
	"fmt"
	"time"
)

var ReceiveChan = make(chan udpMessage)

func Init() {
	const localListenPort = 20057
	const broadcastListenPort = 20058
	const messageSize = 1024

	err := UdpInit(localListenPort, broadcastListenPort, messageSize, sendChan, ReceiveChan)
	if err != nil {
		fmt.Print("Udp_init() error: %s \n", err)
	}

	go aliveSpammer()
	go pollMessages()
}

func pollMessages() { // change name to pollOutgoing or something
	for {
		msg := <-defs.MessageChan

		PrintMessage(msg)

		var i int
		jsonMsg, err := json.Marshal(msg)

		for i = 0; err != nil && i < 10; i++ {
			fmt.Printf("json.Marshal error: %v\n", err)
			jsonMsg, err = json.Marshal(msg)
		}
		if i < 10 {
			sendChan <- udpMessage{raddr: "broadcast", data: jsonMsg, length: len(jsonMsg)}
		}

		time.Sleep(time.Millisecond)
	}
}

func ParseMessage(udpMessage udpMessage) defs.Message {
	fmt.Printf("before parse: %s from %s\n", string(udpMessage.data), udpMessage.raddr)

	var message defs.Message
	if err := json.Unmarshal(udpMessage.data[:udpMessage.length], &message); err != nil {
		fmt.Printf("json.Unmarshal error: %s\n", err)
	}

	message.Addr = udpMessage.raddr
	fmt.Printf("ಠ_ಠ after Unmarshal:       %s\n", message.Addr)
	return message
}

// --------------- PRIVATE: ---------------

var sendChan = make(chan udpMessage)

func aliveSpammer() {
	alive := defs.Message{Kind: defs.Alive, Floor: -1, Button: -1, Cost: -1}
	for {
		defs.MessageChan <- alive
		time.Sleep(defs.SpamInterval)
	}
}

func PrintMessage(msg defs.Message) {
	fmt.Printf("\n-----Message start-----\n")
	switch msg.Kind {
	case defs.Alive:
		fmt.Println("I'm alive")
	case defs.NewOrder:
		fmt.Println("New order")
	case defs.CompleteOrder:
		fmt.Println("Complete order")
	case defs.Cost:
		fmt.Println("Cost:")
	default:
		fmt.Println("Invalid message type!\n")
	}
	fmt.Printf("Floor:  %d\n", msg.Floor)
	fmt.Printf("Button: %d\n", msg.Button)
	fmt.Printf("Cost:   %d\n", msg.Cost)
	fmt.Println("-----Message end-------\n")
}
