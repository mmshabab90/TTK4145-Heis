package main

import (
	"./src/hw"
	"./src/fsm"
	"./src/poller"
	"./src/network"
	"log"
)

func main() {
	if err := hw.Init(); err != nil {
		log.Fatal(err)
	}

	fsm.Init()
	network.Init() //run this from somewhere else?
	poller.Run()
}
