package poller

import (
	"../elev"
	"../fsm"
	"../timer"
	"log"
	"time"
	"fmt"
)

var _ = log.Println

type keypress struct {
	button elev.ButtonType
	floor  int
}

func Run() {
	buttonChan := pollButtons()
	floorChan := pollFloors()

	for {
		select {
		case keypress := <-buttonChan:
			fsm.EventButtonPressed(keypress.floor, keypress.button)
			fmt.Println("Cost: %d", cost.CalculateCost(keypress.floor, keypress.button))
		case floor := <-floorChan:
			fsm.EventFloorReached(floor)
		case <-timer.TimerOut:
			fsm.EventTimerOut()
		}
	}
}

func pollButtons() <-chan keypress {
	c := make(chan keypress)

	go func() {
		var buttonState [elev.NumFloors][elev.NumButtons]bool

		var b elev.ButtonType
		for {
			for f := 0; f < elev.NumFloors; f++ {
				for b = 0; b < elev.NumButtons; b++ {
					if (f == 0 && b == elev.ButtonCallDown) ||
						(f == elev.NumFloors-1 && b == elev.ButtonCallUp) {
						continue
					}
					if elev.GetButton(f, b) {
						if !buttonState[f][b] {
							c <- keypress{button: b, floor: f}
						}
						buttonState[f][b] = true
					} else {
						buttonState[f][b] = false
					}
				}
			}
			time.Sleep(time.Millisecond * 5)
		}
	}()

	return c
}

func pollFloors() <-chan int {
	c := make(chan int)

	go func() {
		oldFloor := elev.GetFloor()

		for {
			newFloor := elev.GetFloor()
			if newFloor != oldFloor && newFloor != -1 {
				c <- newFloor
			}
			oldFloor = newFloor
			time.Sleep(time.Millisecond * 5)
		}
	}()

	return c
}
