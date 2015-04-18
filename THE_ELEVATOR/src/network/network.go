package network

import (
	"fmt"
	"time"
	"encoding/json"
)

// --------------- PUBLIC: ---------------

const (
	Alive int = iota
	NewOrder
	CompleteOrder
	Cost
)

type Message struct {
	Kind int
	Addr string `json:"-"` // skal funke
	Floor int
	Button int
	Cost int
}

type UdpConnection struct { //should this be in udp.go? or in poller.go?
	Addr  string
	Timer *time.Timer
}

// Consider visibility of these three:
var sendChan = make (chan udpMessage)
var ReceiveChan = make (chan udpMessage)


func Init (){
	const localListenPort = 20001
	const broadcastListenPort = 20058
	const messageSize = 1024
	err := Udp_init(localListenPort, broadcastListenPort, messageSize, sendChan, ReceiveChan)		

	if (err != nil){
		fmt.Print("err = %s \n", err)
	}
	
	go aliveSpammer()
}

func Send(message *Message) { // should take a pointer instead
	printMessage(*message)
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		// worry
	} else {
		sendChan <- udpMessage{raddr: "broadcast", data: jsonMessage, length: len(jsonMessage)}
	}
}

func ParseMessage(udpMessage udpMessage) Message {
	var message Message
	err := json.Unmarshal(udpMessage.data, &message)
	if err != nil {
		// handle
	}
	message.Addr = udpMessage.raddr
	return message
}

// --------------- PRIVATE: ---------------

func aliveSpammer() {
	const spamInterval = 500*time.Millisecond
	message := &Message{Kind: Alive}
	for {
		Send(message)
		time.Sleep(spamInterval)
	}
}

func printMessage(msg Message) {
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


