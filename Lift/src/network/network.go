package network

import (
	def "../config"
	"encoding/json"
	"fmt"
	"time"
)

var ReceiveChan = make(chan udpMessage, 10) //buffered with 10 slots

func Init() {
	const localListenPort = 37103
	const broadcastListenPort = 37104
	const messageSize = 1024

	err := UdpInit(localListenPort, broadcastListenPort, messageSize, sendChan, ReceiveChan)
	if err != nil {
		fmt.Print("UdpInit() error: %s \n", err)
	}

	go aliveSpammer()
	go pollMessages()
	fmt.Println("Network initialized")
}

func pollMessages() { // change name to pollOutgoing or something
	for {
		msg := <-def.MessageChan //crashes here for some reason with the orderTimer
		//PrintMessage(msg)

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

func ParseMessage(udpMessage udpMessage) def.Message {
	//fmt.Printf("before parse: %s from %s\n", string(udpMessage.data), udpMessage.raddr)

	var message def.Message
	if err := json.Unmarshal(udpMessage.data[:udpMessage.length], &message); err != nil {
		fmt.Printf("json.Unmarshal error: %s\n", err)
	}

	message.Addr = udpMessage.raddr
	//fmt.Printf("ಠ_ಠ after Unmarshal:       %s\n", message.Addr)
	return message
}

// --------------- PRIVATE: ---------------

var sendChan = make(chan udpMessage)

func aliveSpammer() {
	alive := def.Message{Kind: def.Alive, Floor: -1, Button: -1, Cost: -1}
	for {
		def.MessageChan <- alive
		time.Sleep(def.SpamInterval)
	}
}

func PrintMessage(msg def.Message) {

	if msg.Kind == def.Alive {
		return
	}

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
