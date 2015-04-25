package network

import (
	def "config"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

var ReceiveChan = make(chan udpMessage, 10) //buffered with 10 slots
var sendChan = make(chan udpMessage)

// --------------- PUBLIC: ---------------

func Init() {
	// Ports randomly chosen to reduce likelihood of collision.
	const localListenPort = 37103
	const broadcastListenPort = 37104

	const messageSize = 1024

	err := UdpInit(localListenPort, broadcastListenPort, messageSize, sendChan, ReceiveChan)
	if err != nil {
		fmt.Print("UdpInit() error: %v \n", err)
	}

	go aliveSpammer()
	go forwardOutgoing()

	log.Println("Network initialised.")
}

// aliveSpammer  sends messages on the network to periodically notify
// all lifts that this lift is still online ("alive").
func aliveSpammer() {
	const spamInterval = 400 * time.Millisecond
	alive := def.Message{Category: def.Alive, Floor: -1, Button: -1, Cost: -1}
	for {
		def.OutgoingMsg <- alive
		time.Sleep(spamInterval)
	}
}

// forwardOutgoing continuosly checks for messages to be sent on the network
// by reading the OutgoingMsg channel. Each message read is sent to the udp file
// as JSON.
func forwardOutgoing() { //todo: change name to pollOutgoing or something
	for {
		msg := <-def.OutgoingMsg

		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			fmt.Printf("json.Marshal error: %v\n", err)
		}

		sendChan <- udpMessage{raddr: "broadcast", data: jsonMsg, length: len(jsonMsg)}
	}
}

func ParseMessage(udpMessage udpMessage) def.Message {
	var message def.Message
	if err := json.Unmarshal(udpMessage.data[:udpMessage.length], &message); err != nil {
		fmt.Printf("json.Unmarshal error: %s\n", err)
	}

	message.Addr = udpMessage.raddr

	return message
}

func PrintMessage(msg def.Message) {

	if msg.Category == def.Alive {
		return
	}

	fmt.Printf("\n-----Message start-----\n")
	switch msg.Category {
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
