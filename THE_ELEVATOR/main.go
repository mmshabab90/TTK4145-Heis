package main

import (
	"elev"
	"fsm"
	"timer"
	"temp"
	"time"
	"log"
	"fmt"
)

var _ = elev.Init
var _ = log.Println
var _ = fmt.Println
var _ = time.Sleep

func main() {
	temp.Init()
	
	buttonChan := temp.PollKeypresses()
	floorChan := temp.PollFloor()

	for {
		select {
		case myKeypress := <-buttonChan:
			fsm.EventButtonPressed(myKeypress.Floor, myKeypress.Button)
		case floor := <-floorChan:
			fsm.EventFloorReached(floor)
		case <- timer.TimerOut:
			fsm.EventTimerOut()
		}
	}
}

