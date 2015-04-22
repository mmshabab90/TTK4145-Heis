package network

import (
	"encoding/json"
	"fmt"
	"time"
)

// Generic network message. No other messages are ever sent on the network.
const (
	alive int = iota + 1
	newOrder
	completeOrder
	cost
)

type message struct {
	kind   int
	floor  int
	button int
	cost   int
	addr   string
}

var receiveChan = make(chan udpMessage)
var incoming = make(chan message)
var outgoing = make(chan message)

// Move these out of here:
const spamInterval = 30 * time.Second
const resetTime = 120 * time.Second // rename

func Init(floorCompleted <- chan int) {
	const localListenPort = 37103
	const broadcastListenPort = 37104
	const messageSize = 1024

	err := UdpInit(localListenPort, broadcastListenPort, messageSize, sendChan, receiveChan)
	if err != nil {
		fmt.Print("UdpInit() error: %s \n", err)
	}

	go floorCompleteForwarder(floorCompleted)
	go aliveSpammer()
	go pollIncoming()
	go pollOutgoing()
}

func floorCompleteForwarder(floorCompleted <-chan int) {
	for {
		floor := <- floorCompleted
		outgoing <- message{
			kind: completeOrder,
			floor: floor}
	}
}

func pollIncoming() { // merge with pollOutgoing?
	for {
		udpMsg := <- receiveChan
		var msg message
		json.Unmarshal(udpMsg.data[:udpMsg.length], &msg)
		// acceptance test msg here!
		msg.addr = udpMsg.raddr
		incoming <- msg
	}
}

func pollOutgoing() {
	for {
		msg := <-outgoing

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
	alive := message{
		kind: alive,
		floor: -1,
		button: -1,
		cost: -1}
	for {
		outgoing <- alive
		time.Sleep(spamInterval)
	}
}

func PrintMessage(msg message) {
	fmt.Printf("\n-----Message start-----\n")
	switch msg.kind {
	case alive:
		fmt.Println("I'm alive")
	case newOrder:
		fmt.Println("New order")
	case completeOrder:
		fmt.Println("Complete order")
	case cost:
		fmt.Println("Cost:")
	default:
		fmt.Println("Invalid message type!\n")
	}
	fmt.Printf("Floor: %d\n", msg.floor)
	fmt.Printf("Button: %d\n", msg.button)
	fmt.Printf("Cost:   %d\n", msg.cost)
	fmt.Println("-----Message end-------\n")
}
