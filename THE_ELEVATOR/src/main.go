package main

import (
	"./elev"
	"./fsm"
	"./queue"
	"log"
)

type Keypress struct {
	button Elev_button_type_t
	floor int
}

func main() {
	init()

	buttonChan = make(chan Keypress) // does this need to be buffered to handle many keypresses happening "at once"?
	floorChan = make(chan int)

	go pollKeypresses(buttonChan)
	go pollFloor(floorChan)

	for {
		select {
		case keypress := <- buttonChan:
			fsm.EventButtonPressed(keypress.floor, keypress.button)
		case floor := <- floorChan:
			fsm.EventFloorReached(floor)
		}
	}
}

func init() {
	if !elev.Init() {
		log.Fatalln("Io_init() failed!")
	}

	elev.SetMotorDirection(DirnDown)
	for elev.GetFloor() == -1 {}
	elev.SetMotorDirection(DirnStop)

	fsm.Init()
	// Add some error handling here.
}

func pollKeypresses(c chan Keypress) {
	var buttonState = [nFloors][nButtons]bool{
		{false, false, false},
		{false, false, false},
		{false, false, false},
		{false, false, false} // find a better way to do this
	}

	for {
		for f := 0; f < nFloors; f++ {
			for b := 0; b < nButtons; b++ {
				if (f == 0 && b == ButtonCallDown) ||
				(f == nFloors-1 && b == ButtonCallUp) {
					continue
				}
				if elev.GetButton(f, b) {
					if !buttonState[f][b] {
						c <- Keypress{button = b, floor = f}
					}
					buttonState[f][b] = true
				} else {
					buttonState[f][b] = false
				}
			}
		}
	}
}

func pollFloor(c chan int) {
	oldFloor := elev.GetFloor()

	for {
		newFloor := elev.GetFloor()
		if newFloor != oldFloor {
			if newFloor != -1 {
				c <- newFloor
			}
			oldFloor = newFloor
		}
	}
}