package network

import (
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

type Message struct {
	Kind int
	Addr string `json:"-"` // skal funke
	Floor int
	Button int
	Cost int
}

type UdpConnection struct { //should this be in udp.go
	Addr  string
	Timer *time.Timer
}

// Consider visibility of these three:
var sendChan = make (chan udpMessage)
var ReceiveChan = make (chan udpMessage)
var ConnectionTimer	 = make(chan UdpConnection)

func Init (){
	err := Udp_init(20001, 20058, 1024, sendChan, ReceiveChan)	

	if (err != nil){
		fmt.Print("err = %s \n", err)
	}
}

func Send(message Message) {
	printMessage(message)
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		// worry
	} else {
		sendChan <- udpMessage{raddr: "broadcast", data: msg, length: len(msg)}
		time.Sleep(500*time.Millisecond)	// What's this for?
	}
}

func ReceiveMsg(){ // bad abstraction! doesn't just receive msg. GIVE THIS NEW NAME!
	connectionMap := make(map[string] UdpConnection)
	for {
		select{
		case rcvMsg := <- ReceiveChan:
			Print_udp_message(rcvMsg)
			
			//keep track of witch connections that exist
			if connection, exist := connectionMap[rcvMsg.raddr]; exist {
				connection.Timer.Reset(1*time.Second)
				fmt.Println("timer reset for IP: ")
				fmt.Println(rcvMsg.raddr)
			} else {
				newConnection := UdpConnection{rcvMsg.raddr, time.NewTimer(1*time.Second)}
				connectionMap[rcvMsg.raddr] = newConnection
				fmt.Println("New connection, with IP: ")
				fmt.Println(rcvMsg.raddr)
				go connectionTimer(&newConnection)
			}
		//deletes connection when timer goes out
		case connection := <- ConnectionTimer:
			fmt.Println(connection.Addr, "is dead")
			delete(connectionMap, connection.Addr)
			for key, _ := range connectionMap {
				fmt.Println(key)
			}
		}
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
		log.Println("Invalid message type!\n")
	}
}

// --------------- PRIVATE: ---------------

func connectionTimer(connection *UdpConnection) {
	for {
		<- connection.Timer.C
		ConnectionTimer <- *connection
	}
}
