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
	
	buttonChan := make(chan temp.Keypress) // does this need to be buffered to handle many keypresses happening "at once"?
	floorChan := make(chan int)

	go temp.PollKeypresses(buttonChan)
	go temp.PollFloor(floorChan)

	for {
		select {
		case myKeypress := <-buttonChan:
			log.Println("Got keypress")
			fsm.EventButtonPressed(myKeypress.Floor, myKeypress.Button)
		case floor := <-floorChan:
			log.Println("Got floor")
			fsm.EventFloorReached(floor)
		case <- timer.TimerChan:
			log.Println("Got timeout")
			fsm.EventTimerOut()
		}
	}
}

