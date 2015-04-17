package network

import (
	"fmt"
	"time"
)

// --------------- PUBLIC: ---------------

type UdpConnection struct { //should this be in udp.go
	Addr  string
	Timer *time.Timer
}

//must these be global?
var Send_ch = make (chan Udp_message)
var Receive_ch = make (chan Udp_message)
//this must be global (i think)
var ConnectionTimer	 = make(chan UdpConnection)

func Init (){
	err := Udp_init(20001, 20058, 1024, Send_ch, Receive_ch)	

	if (err != nil){
		fmt.Print("err = %s \n", err)
	}
}

func SendMsg(msg []byte){
	sndMsg := Udp_message{Raddr:"broadcast", Data:msg, Length:len(msg)}
	Send_ch <- sndMsg
	time.Sleep(500*time.Millisecond)
}

func ReceiveMsg(){ //bad abstraction? doesn't just receive msg. GIVE THIS NEW NAME!
	connectionMap := make(map[string] UdpConnection)
	for {
		select{
		case rcvMsg := <- Receive_ch:
			Print_udp_message(rcvMsg)
			
			//keep track of witch connections that exist
			if connection, exist := connectionMap[rcvMsg.Raddr]; exist {
				connection.Timer.Reset(1*time.Second)
				fmt.Println("timer reset for IP: ")
				fmt.Println(rcvMsg.Raddr)
			} else {
				newConnection := UdpConnection{rcvMsg.Raddr, time.NewTimer(1*time.Second)}
				connectionMap[rcvMsg.Raddr] = newConnection
				fmt.Println("New connection, with IP: ")
				fmt.Println(rcvMsg.Raddr)
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

// --------------- PRIVATE: ---------------

func connectionTimer(connection *UdpConnection) {
	for {
		<- connection.Timer.C
		ConnectionTimer <- *connection
	}
}
