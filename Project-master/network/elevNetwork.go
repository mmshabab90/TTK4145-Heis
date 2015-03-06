package main

import (
	"log"
	"net"
	"time"
)

type UdpConnection struct {
	Addr  string
	Timer *time.Timer
}

func Broadcast(addr string, msg string) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	udpBroadcast, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	defer udpBroadcast.Close()

	udpBroadcast.Write([]byte(msg))
}

func StillAliveBroadcast() {
	for {
		Broadcast("129.241.187.255:20014", "I'm alive!")
		time.Sleep(100*time.Millisecond)
	}
}

func Listen(port string, timeoutChan chan UdpConnection, msgChan chan string) {
	udpAddr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		log.Fatal(err)
	}

	udpListen, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	defer udpListen.Close()

	var buffer [1024]byte
	var connectionMap map[string] UdpConnection 

	for {
		_, ipAddr, err := udpListen.ReadFromUDP(buffer[:])
		if err != nil {
			log.Fatal(err)
		}

		if connection, exist := connectionMap[ipAddr.String()]; exist {
			connection.Timer.Reset(1*time.Second)
		} else {
			newConnection := UdpConnection{ipAddr.String(), time.NewTimer(1*time.Second)}
			connectionMap[ipAddr.String()] = newConnection
			go connectionTimer(&newConnection, timeoutChan)
		}
		
		if string(buffer[:]) != "I'm alive" {
			msgChan <- string(buffer[:])  
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func connectionTimer(connection *UdpConnection, timeoutChan chan UdpConnection) {
	for {
		select {
		case <- connection.Timer.C:
			timeoutChan <- *connection
		}
	}
}

func main() {
	msgChan := make(chan string)
	timeoutChan := make(chan UdpConnection)
	
	go StillAliveBroadcast()
	go Listen(":20014", timeoutChan, msgChan)

	for {
		select {
		case connection := <- timeoutChan:
			log.Println(connection.Addr, "is dead")
		case msg := <- msgChan:
			log.Println(msg)
		}
	}
}






