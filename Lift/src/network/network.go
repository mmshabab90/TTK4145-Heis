package network

import (
	def "config"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

var receiveChan = make(chan udpMessage, 10)
var sendChan = make(chan udpMessage)

func Init(outgoingMsg, incomingMsg chan def.Message) {
	// Ports randomly chosen to reduce likelihood of port collision.
	const localListenPort = 37103
	const broadcastListenPort = 37104

	const messageSize = 1024

	err := udpInit(localListenPort, broadcastListenPort, messageSize, sendChan, receiveChan)
	if err != nil {
		fmt.Print("UdpInit() error: %v \n", err)
	}

	go aliveSpammer(outgoingMsg)
	go forwardOutgoing(outgoingMsg)
	go forwardIncoming(incomingMsg)

	log.Println(def.ClrG, "Network initialised.", def.ClrN)
}

// aliveSpammer  sends messages on the network to periodically notify
// all lifts that this lift is still online ("alive").
func aliveSpammer(outgoingMsg chan<- def.Message) {
	const spamInterval = 400 * time.Millisecond
	alive := def.Message{Category: def.Alive, Floor: -1, Button: -1, Cost: -1}
	for {
		outgoingMsg <- alive
		time.Sleep(spamInterval)
	}
}

// forwardOutgoing continuosly checks for messages to be sent on the network
// by reading the OutgoingMsg channel. Each message read is sent to the udp file
// as JSON.
func forwardOutgoing(outgoingMsg chan def.Message) { //todo: change name to pollOutgoing or something
	for {
		msg := <-outgoingMsg

		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			fmt.Printf("json.Marshal error: %v\n", err)
		}

		sendChan <- udpMessage{raddr: "broadcast", data: jsonMsg, length: len(jsonMsg)}
	}
}

func forwardIncoming(incomingMsg chan def.Message) {
	for {
		udpMessage := <-receiveChan
		var message def.Message

		if err := json.Unmarshal(udpMessage.data[:udpMessage.length], &message); err != nil {
			fmt.Printf("json.Unmarshal error: %s\n", err)
		}

		message.Addr = udpMessage.raddr

		incomingMsg <- message
	}
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
