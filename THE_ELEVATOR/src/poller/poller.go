package poller

import (
	"../elev"
	"../fsm"
	"../timer"
	"../cost"
	"../network"
	"log"
	"time"
	"fmt"
	"encoding/json"
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
			fmt.Printf("Cost: %d\n", cost.CalculateCost(keypress.floor, keypress.button))
		case floor := <-floorChan:
			fsm.EventFloorReached(floor)
		case <-timer.TimerOut:
			fsm.EventTimerOut()
		case udpMessage := <-network.ReceiveChan:
			handleMessage(parseMessage(udpMessage))
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

func parseMessage(udpMessage network.Udp_message) network.Message { // work this into network package!
	var message network.Message
	message = json.Unmarshal(udpMessage.Data)
	message.Addr = udpMessage.Raddr
	return message
}

func handleMessage(message network.Message) {
	switch message.Kind {
		case queue.Alive:
			// reset lift timer (not door timer lol)
		case queue.NewOrder:
			costMessage := queue.Message{
				Kind: queue.Cost,
				Floor: message.Floor,
				Button: message.Button,
				Cost: cost.CalculateCost(message.Floor, message.Button)}
			
		case queue.CompleteOrder:
			// remove from queues
		case queue.Cost:
			// notify assignment routine
	}
}
