package main

import (
	"./network"
)

func main() {
	network.NetworkInit()
	network.ReceiveMsg()
}
