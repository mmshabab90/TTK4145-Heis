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
	
	//buttonChan := make(chan temp.Keypress) // does this need to be buffered to handle many keypresses happening "at once"?
	//floorChan := make(chan int)

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
