package main

import (
	"./network"
)

func main() {
	network.Init()
	for {
		network.SendMsg("Hello Morten")
	}
}
