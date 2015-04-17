package main

import (
	"./src/hw"
	"./src/fsm"
	"./src/poller"
	"log"
)

func main() {
	if !hw.Init() {
		log.Fatalln("Io_init() failed!")
	}
	fsm.Init()

	// Move to defined state:
	hw.SetMotorDirection(hw.DirnDown)
	floor := hw.GetFloor()
	for floor == -1 {
		floor = hw.GetFloor()
	}
	hw.SetFloorIndicator(floor)
	hw.SetMotorDirection(hw.DirnStop)

	poller.Run()
}
