package main

import (
	def "./src/config"
	"./src/fsm"
	"./src/hw"
	"./src/network"
	"./src/queue"
	"errors"
	"fmt"
	"log"
)

const debugPrint = false

var _ = log.Println
var _ = fmt.Println
var _ = errors.New

type keypress struct {
	button int
	floor  int
}

var orderTimeoutChan = make(chan order)

func main() {
	if err := hw.Init(); err != nil {
		log.Fatal(err)
	}

	eventNewOrder := make(chan bool)
	eventFloorReached := make(chan int)

	fsm.Init(eventNewOrder, eventFloorReached)
	network.Init()

	run(eventNewOrder)
}

func run(eventNewOrder <-chan bool,
	eventFloorReached <-chan int) {
	buttonChan := pollButtons()
	floorChan := pollFloors()

	for {
		select {
		case key := <-buttonChan:
			queue.AddKeypressOrder(key.floor, key.button)
			eventNewOrder <- true
		case floor := <-floorChan:
			eventFloorReached <- floor
		}
	}
}

func pollButtons() <-chan keypress {
	c := make(chan keypress)

	go func() {
		var buttonState [def.NumFloors][def.NumButtons]bool

		for {
			for f := 0; f < def.NumFloors; f++ {
				for b := 0; b < def.NumButtons; b++ {
					if (f == 0 && b == def.ButtonDown) ||
						(f == def.NumFloors-1 && b == def.ButtonUp) {
						continue
					}
					if hw.ReadButton(f, b) {
						if !buttonState[f][b] {
							c <- keypress{button: b, floor: f}
						}
						buttonState[f][b] = true
					} else {
						buttonState[f][b] = false
					}
				}
			}
			time.Sleep(time.Millisecond)
		}
	}()

	return c
}

func pollFloors() <-chan int {
	c := make(chan int)
	go func() {
		oldFloor := hw.Floor()

		for {
			newFloor := hw.Floor()
			if newFloor != oldFloor && newFloor != -1 {
				c <- newFloor
			}
			oldFloor = newFloor
			time.Sleep(time.Millisecond)
		}
	}()
	return c
}
