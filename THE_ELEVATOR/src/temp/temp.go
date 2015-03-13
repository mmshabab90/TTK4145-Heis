package temp

import (
	"elev"
	"log"
	"fsm"
	"time"
)

type Keypress struct {
	Button elev.ButtonType
	Floor int
}

func Init() {
	if !elev.Init() {
		log.Fatalln("Io_init() failed!")
	}

	floor := elev.GetFloor()
	elev.SetMotorDirection(elev.DirnDown)
	for floor == -1 {floor = elev.GetFloor()}
	elev.SetFloorIndicator(floor)
	elev.SetMotorDirection(elev.DirnStop)

	fsm.Init()
	// Add some error handling here.
}

func PollKeypresses() <- chan Keypress {
	c := make(chan Keypress)

	go func() {
		var buttonState [elev.NumFloors][elev.NumButtons] bool
		
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
							c <- Keypress{Button: b, Floor: f}
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

func PollFloor() <-chan int {
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
