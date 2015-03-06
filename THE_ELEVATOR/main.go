package main

import (
	"elev"
	"fsm"
	"temp"
	"log"
)

var _ = elev.Init
var _ = log.Println

func main() {
	temp.Init()
	
	buttonChan := make(chan keypress) // does this need to be buffered to handle many keypresses happening "at once"?
	floorChan := make(chan int)

	go PollKeypresses(buttonChan)
	go PollFloor(floorChan)

	for {
		select {
		case myKeypress := <-buttonChan:
			fsm.EventButtonPressed(myKeypress.Floor, myKeypress.Button)
		case floor := <-floorChan:
			fsm.EventFloorReached(floor)
		}
	}
}

