package main

import (
	"./src/hw"
	"./src/fsm"
	"./src/poller"
	"log"
)

func main() {
	if err := hw.Init(); err != nil {
		log.Fatal(err)
	}

	fsm.Init()
	poller.Run()
}
