package main

import (
	"./src/elev"
	"./src/fsm"
	"./src/poller"
	"./src/timer"
	"log"
)

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

	poller.Run()
}
