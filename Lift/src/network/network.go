package network

import (
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
var incoming = make(chan Message)
var Outgoing = make(chan Message)

// Move these out of here:
const spamInterval = 30 * time.Second
const resetTime = 120 * time.Second // rename

func Init(
	floorCompleted <- chan int,
	deathChan chan<- string,
	addRemoteOrder chan<- RemoteOrder,
	costMessage <- chan Message,
	onlineLifts map[string]*time.Timer) {
	
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
	go liftAssigner(addRemoteOrder, costMessage, onlineLifts)
}

func floorCompleteForwarder(floorCompleted <-chan int) {
	for {
		floor := <- floorCompleted
		Outgoing <- Message{
			Kind: CompleteOrder,
			Floor: floor}
	}
}

func pollIncoming() { // merge with pollOutgoing?
	for {
		udpMsg := <- receiveChan
		var msg Message
		json.Unmarshal(udpMsg.data[:udpMsg.length], &msg)
		// acceptance test msg here!
		msg.Addr = udpMsg.raddr
		incoming <- msg
	}
}

func pollOutgoing() {
	for {
		msg := <-Outgoing

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
	alive := Message{
		Kind: Alive,
		Floor: -1,
		Button: -1,
		Cost: -1}
	for {
		Outgoing <- alive
		time.Sleep(spamInterval)
	}
}

func PrintMessage(msg Message) {
	fmt.Printf("\n-----Message start-----\n")
	switch msg.Kind {
	case Alive:
		fmt.Println("I'm alive")
	case NewOrder:
		fmt.Println("New order")
	case CompleteOrder:
		fmt.Println("Complete order")
	case Cost:
		fmt.Println("Cost:")
	default:
		fmt.Println("Invalid message type!\n")
	}
	fmt.Printf("Floor: %d\n", msg.Floor)
	fmt.Printf("Button: %d\n", msg.Button)
	fmt.Printf("Cost:   %d\n", msg.Cost)
	fmt.Println("-----Message end-------\n")
}
