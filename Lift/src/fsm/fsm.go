package fsm

import (
	"../defs"
	"../hw"
	"../queue"
	"fmt"
	"log"
	"time"
)

type stateType int // kill this

const (
	idle stateType = iota
	moving
	doorOpen
)

var state stateType
var storey int
var direction int
var departDirection int

var doorReset = make(chan bool)

const doorOpenTime = 1 * time.Second

// --------------- PUBLIC: ---------------

var DoorTimeoutChan = make(chan bool)

func Init() {
	log.Println("FSM Init")
	go startTimer()
	state = idle
	direction = defs.DirStop
	storey = hw.Storey()
	if storey == -1 {
		storey = hw.MoveToDefinedState()
	}
	departDirection = defs.DirDown
	go syncLights()
}

func EventInternalButtonPressed(buttonStorey int, buttonType int) {
	fmt.Printf("\n\n   ☺      Event internal button (storey %d %s) pressed in state %s\n",
		buttonStorey, buttonString(buttonType), stateString(state))
	queue.Print()
	switch state {
	case idle:
		queue.AddLocalOrder(buttonStorey, buttonType)
		switch direction = queue.ChooseDirection(storey, direction); direction {
		case defs.DirStop:
			hw.SetDoorOpenLamp(true)
			queue.RemoveOrdersAt(storey)
			doorReset <- true
			state = doorOpen
		case defs.DirUp, defs.DirDown:
			hw.SetMotorDirection(direction)
			departDirection = direction
			state = moving
		}
	case doorOpen:
		if storey == buttonStorey {
			doorReset <- true
		} else {
			queue.AddLocalOrder(buttonStorey, buttonType)
		}
	case moving:
		queue.AddLocalOrder(buttonStorey, buttonType)
	default:
		log.Fatalf("State %d is invalid!\n", state)
	}

	defs.SyncLightsChan <- true
}

func EventExternalButtonPressed(buttonStorey int, buttonType int) {
	fmt.Printf("\n\n   ☺      Event external button (storey %d %s) pressed in state %s\n",
		buttonStorey, buttonString(buttonType), stateString(state))
	queue.Print()
	switch state {
	case idle, doorOpen, moving:
		// send order on network
		message := defs.Message{Kind: defs.NewOrder, Storey: buttonStorey, Button: buttonType, Cost: -1}
		defs.MessageChan <- message
	default:
		//
	}

	defs.SyncLightsChan <- true
}

func EventExternalOrderGivenToMe() {
	fmt.Printf("\n\n   ☺      Event external order given to me.\n")
	queue.Print()

	if queue.IsLocalEmpty() {
		// strange
	}
	switch state {
	case idle:
		switch direction = queue.ChooseDirection(storey, direction); direction {
		case defs.DirStop:
			hw.SetDoorOpenLamp(true)
			fmt.Println("in eventExternalOrderGivenToMe dir:stop")
			queue.RemoveOrdersAt(storey) //her tror jeg buggen ligger!
			doorReset <- true
			state = doorOpen
		case defs.DirUp, defs.DirDown:
			hw.SetMotorDirection(direction)
			departDirection = direction
			state = moving
		}
	default:
		fmt.Println("   ☺      EventExternalOrderGivenToMe(): Not in idle, will ignore.")
	}
	defs.SyncLightsChan <- true
}

func EventStoreyReached(newStorey int) {
	fmt.Printf("\n\n   ☺      Event storey %d reached in state %s\n", newStorey, stateString(state))
	queue.Print()
	storey = newStorey
	hw.SetStoreyLamp(storey)
	switch state {
	case moving:
		if queue.ShouldStop(storey, direction) {
			hw.SetMotorDirection(defs.DirStop)
			hw.SetDoorOpenLamp(true)
			queue.RemoveOrdersAt(storey)
			go queue.SendOrderCompleteMessage(storey)
			doorReset <- true
			state = doorOpen
		} else {
			departDirection = direction
		}
	default:
		log.Printf("Makes no sense to arrive at a storey in state %s.\n", stateString(state))
	}
	defs.SyncLightsChan <- true
}

func EventDoorTimeout() {
	fmt.Printf("\n\n   ☺      Event door timeout in state %s\n", stateString(state))
	queue.Print()
	switch state {
	case doorOpen:
		direction = queue.ChooseDirection(storey, direction)
		hw.SetDoorOpenLamp(false)
		hw.SetMotorDirection(direction)
		if direction == defs.DirStop {
			state = idle
		} else {
			state = moving
			departDirection = direction
		}
	default:
		log.Fatalf("Makes no sense to time out when not in state door open\n")
	}
	defs.SyncLightsChan <- true
}

func Direction() int {
	return direction
}

func DepartDirection() int {
	return departDirection
}

func Storey() int {
	return storey
}

func startTimer() {
	timer := time.NewTimer(0)
	timer.Stop()
	for {
		select {
		case <-doorReset:
			timer.Reset(doorOpenTime)
		case <-timer.C:
			timer.Stop()
			EventDoorTimeout()
		}
	}
}

func syncLights() {
	for {
		<-defs.SyncLightsChan

		for f := 0; f < defs.NumStoreys; f++ {
			for b := 0; b < defs.NumButtons; b++ {
				if (b == defs.ButtonUp && f == defs.NumStoreys-1) ||
					(b == defs.ButtonDown && f == 0) {
					continue
				} else {
					hw.SetButtonLamp(f, b, queue.IsOrder(f, b))
				}
			}
		}
		time.Sleep(time.Millisecond)
	}
}

func stateString(state stateType) string {
	switch state {
	case idle:
		return "idle"
	case moving:
		return "moving"
	case doorOpen:
		return "door open"
	default:
		return "error: bad state"
	}
}

func buttonString(button int) string {
	switch button {
	case defs.ButtonUp:
		return "up"
	case defs.ButtonDown:
		return "down"
	case defs.ButtonCommand:
		return "command"
	default:
		return "error: bad button"
	}
}
