package main

import (
	"./defs"
	"./network"
	"time"
)

func main() {
	network.Init()
	message1 := defs.Message{
		Kind:   defs.Cost,
		Storey: 0,
		Button: 0,
		Cost:   2}

	message2 := defs.Message{Kind: defs.Alive}
	for {

		network.Send(message1)
		time.Sleep(time.Second)

		network.Send(message2)
		time.Sleep(time.Second)
	}
}
