package main

import (
	"./network"
)

func main() {
	network.Init()
	network.ReceiveMsg()
}
