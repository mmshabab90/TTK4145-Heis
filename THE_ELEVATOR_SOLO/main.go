package main

import (
	"./src/elev"
	"fmt"
	"./src/fsm"
	"log"
	"./src/poller"
	"time"
	"./src/timer"
)

var _ = elev.Init
var _ = log.Println
var _ = fmt.Println
var _ = time.Sleep

func main() {
	if !elev.Init() {
		log.Fatalln("Io_init() failed!")
	}
	fsm.Init()
	timer.Init()

	// Move to defined state:
	elev.SetMotorDirection(elev.DirnDown)
	floor := elev.GetFloor()
	for floor == -1 {
		floor = elev.GetFloor()
	}
	elev.SetFloorIndicator(floor)
	elev.SetMotorDirection(elev.DirnStop)

	// Run elevator:
	poller.Run()
}
