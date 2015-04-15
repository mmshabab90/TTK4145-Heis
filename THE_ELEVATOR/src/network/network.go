package network

import (
	"fmt"
	"./udp"
	"time"
)

type UdpConnection struct {
	Addr  string
	Timer *time.Timer
}

//must these be global?
var Send_ch = make (chan udp.Udp_message)
var Receive_ch = make (chan udp.Udp_message)
var ConnectionTimerChan = make(chan UdpConnection)

func Init (){
	err := udp.Udp_init(20001, 20058, 1024, Send_ch, Receive_ch)	

	if (err != nil){
		fmt.Print("err = %s \n", err)
	}
}

func SendMsg(msg string){
	sndMsg := udp.Udp_message{Raddr:"broadcast", Data:msg, Length:len(msg)}
	Send_ch <- sndMsg
	fmt.Println("Msg sent")
	udp.Print_udp_message(sndMsg)
	time.Sleep(500*time.Millisecond)
}

func ReceiveMsg(){
	var connectionMap map[string] UdpConnection
	for {
		rcvMsg := <- Receive_ch
		fmt.Println("Msg received")
		udp.Print_udp_message(rcvMsg)
		
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
	}
}

func connectionTimer(connection *UdpConnection) {
	for {
		select {
		case <- connection.Timer.C:
			ConnectionTimerChan <- *connection
		}
	}
}
