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

type UdpConnection struct { //should this be in udp.go?
	Addr  string
	Timer *time.Timer
}

// Consider visibility of these three:
var sendChan = make (chan udpMessage)
var ReceiveChan = make (chan udpMessage)


func Init (){
	err := Udp_init(20001, 20058, 1024, sendChan, ReceiveChan)		

	if (err != nil){
		fmt.Print("err = %s \n", err)
	}
	
	go aliveSpammer()
}

func Send(message Message) {
	printMessage(message)
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		// worry
	} else {
		sendChan <- udpMessage{raddr: "broadcast", data: jsonMessage, length: len(jsonMessage)}
	}
}

func ParseMessage(udpMessage udpMessage) Message { // work this into network package!
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
	message := Message{Kind: Alive}
	for {
		Send(message)
		time.Sleep(500*time.Millisecond)
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


