package main

import (
	"./network"
)

func main() {
	network.NetworkInit()
	for {
		network.SendMsg("Hello Morten")
	}
}
