package main

import (
	"./network"
)

func main() {
	network.Init()
	for {
		message := <- network.ReceiveChan
		network.PrintMessage(network.ParseMessage(message))
	}
}
