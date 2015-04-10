package network

import (
	"fmt"
	"./udp"
	"time"
)

type UdpConnection struct { //should be in udp-package
	Addr  string
	Timer *time.Timer //does these have to be global?
}

var Send_ch = make (chan udp.Udp_message)
var Receive_ch = make (chan udp.Udp_message)
var ConnectionTimerChan = make(chan UdpConnection)

func Print_udp_message(msg udp.Udp_message){
	fmt.Printf("msg:  \n \t raddr = %s \n \t data = %s \n \t length = %v \n", msg.Raddr, msg.Data, msg.Length)
}

func NetworkInit (){
	err := udp.Udp_init(20001, 20014, 1024, Send_ch, Receive_ch)	

	if (err != nil){
		fmt.Print("err = %s \n", err)
	}
}

func StillAliveBroadcast(){ //global?
	for {
		sndMsg := udp.Udp_message{Raddr:"broadcast", Data:"I'm alive", Length:9}
		Send_ch <- sndMsg
		fmt.Println("Msg sent")
		Print_udp_message(sndMsg)
		time.Sleep(500*time.Second)
	} 
}

func ListenForLiveElevators(){ //global?
	var connectionMap map[string] UdpConnection
	for {
		rcvMsg := <- Receive_ch
		
		if connection, exist := connectionMap[rcvMsg.Raddr]; exist {
			connection.Timer.Reset(1*time.Second)
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
