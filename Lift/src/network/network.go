package network

import (
	def "../config"
	"encoding/json"
	"fmt"
	"time"
)

// Generic network message. No other messages are ever sent on the network.
const (
	Alive int = iota + 1
	NewOrder
	CompleteOrder
	Cost
)

type Message struct {
	Kind   int
	Floor  int
	Button int
	Cost   int
	Addr   string
}

var receiveChan = make(chan udpMessage)

// Move these out of here:
const spamInterval = 30 * time.Second
const resetTime = 120 * time.Second // rename

func Init() {
	const localListenPort = 37103
	const broadcastListenPort = 37104
	const messageSize = 1024

	err := UdpInit(localListenPort, broadcastListenPort, messageSize, sendChan, receiveChan)
	if err != nil {
		fmt.Print("UdpInit() error: %s \n", err)
	}

	go aliveSpammer()
	go pollIncoming()
	go pollOutgoing()
}

func pollIncoming() { // merge with pollOutgoing?
	for {
		udpMsg <- receiveChan
		msg := new(Message)
		json.Unmarshal(udpMsg.data[:udpMsg.lenght], &msg)
		// acceptance test msg here!
		msg.Addr = udpMsg.raddr
		incoming <- msg
	}
}

func pollOutgoing() {
	for {
		msg := <-def.Outgoing

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

// --------------- PRIVATE: ---------------

var sendChan = make(chan udpMessage)

func aliveSpammer() {
	alive := def.Message{Kind: def.Alive, Floor: -1, Button: -1, Cost: -1}
	for {
		def.Outgoing <- alive
		time.Sleep(def.SpamInterval)
	}
}

func PrintMessage(msg def.Message) {
	fmt.Printf("\n-----Message start-----\n")
	switch msg.Kind {
	case def.Alive:
		fmt.Println("I'm alive")
	case def.NewOrder:
		fmt.Println("New order")
	case def.CompleteOrder:
		fmt.Println("Complete order")
	case def.Cost:
		fmt.Println("Cost:")
	default:
		fmt.Println("Invalid message type!\n")
	}
	fmt.Printf("Floor: %d\n", msg.Floor)
	fmt.Printf("Button: %d\n", msg.Button)
	fmt.Printf("Cost:   %d\n", msg.Cost)
	fmt.Println("-----Message end-------\n")
}
