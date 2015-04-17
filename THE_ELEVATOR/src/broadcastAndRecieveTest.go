package main

//this was useless

import (
	"./network"
	"fmt"
)

func main() {
	network.Init()
	go send()
	go network.ReceiveMsg()
	for {
		connection := <- network.ConnectionTimer
		fmt.Println(connection.Addr, "is dead")
	}
}

func send() {
	for {
		network.SendMsg("Hello Somebody")
	}
}
