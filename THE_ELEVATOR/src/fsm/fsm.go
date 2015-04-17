package fsm

import (
	"../elev"
	"../queue"
	"../network"
	"fmt"
	"log"
	"time"
)

type stateType int
const (
	idle stateType = iota
	moving
	doorOpen
)

var state stateType
var floor int
var direction elev.DirnType
var departDirection elev.DirnType

var DoorTimeout = make(chan bool)
var doorReset = make(chan bool)
const doorOpenTime = 3 * time.Second

func Init() {
	log.Println("FSM Init")
	queue.Init()
	runTimer()
	state = idle
	direction = elev.DirnStop
	floor = elev.GetFloor()
	departDirection = elev.DirnDown
	syncLights()
}

func runTimer() {
	timer := time.NewTimer(0)
	timer.Stop()

	go func() {
		for {
			select {
			case <-doorReset:
				timer.Reset(doorOpenTime)
			case <-timer.C:
				DoorTimeout <- true
				timer.Stop()
			}
		}
	}()
}


func EventButtonPressed(buttonFloor int, buttonType elev.ButtonType) {
	fmt.Print("Event button pressed in state ")
	switch state {
	case idle:
		fmt.Println("idle")
		queue.AddOrder(buttonFloor, buttonType)
		direction = queue.ChooseDirection(floor, direction)
		if direction == elev.DirnStop {
			elev.SetDoorOpenLamp(true)
			queue.RemoveOrdersAt(floor)
			doorReset <- true
			state = doorOpen
		} else {
			elev.SetMotorDirection(direction)
			departDirection = direction
			state = moving
		}
	case doorOpen:
		fmt.Println("door open")
		if floor == buttonFloor {
			doorReset <- true
		} else {
			queue.AddOrder(buttonFloor, buttonType)
		}
	case moving:
		fmt.Println("moving")
		queue.AddOrder(buttonFloor, buttonType)
	default:
		log.Fatalf("State %d is invalid!\n", state)
	}
	syncLights()
}

func EventFloorReached(newFloor int) {
	fmt.Print("Event floor reached in state ")
	floor = newFloor
	elev.SetFloorIndicator(floor)
	switch state {
	case moving:
		fmt.Println("moving")
		if queue.ShouldStop(floor, direction) {
			elev.SetMotorDirection(elev.DirnStop)
			elev.SetDoorOpenLamp(true)
			queue.RemoveOrdersAt(floor)
			doorReset <- true
			state = doorOpen
		} else {
			departDirection = direction
		}
	default:
		log.Fatalf("Makes no sense to arrive at a floor in state %d", state)
	}
	syncLights()
}

func EventTimerOut() {
	fmt.Print("Event timer out in state ")
	switch state {
	case doorOpen:
		fmt.Println("door open")
		direction = queue.ChooseDirection(floor, direction)
		elev.SetDoorOpenLamp(false)
		elev.SetMotorDirection(direction)
		if direction == elev.DirnStop {
			state = idle
		} else {
			state = moving
			departDirection = direction
		}
	default:
		log.Fatalf("Makes no sense to time out when not in doorOpen\n")
	}
	syncLights()
}

func GetDirection() elev.DirnType {
	return direction
}

func GetFloor() int {
	return floor
}

func syncLights() {
	var b elev.ButtonType
	for f := 0; f < elev.NumFloors; f++ {
		for b = 0; b < elev.NumButtons; b++ {
			if (b == elev.ButtonCallUp && f == elev.NumFloors-1) ||
				(b == elev.ButtonCallDown && f == 0) {
				continue
			} else {
				elev.SetButtonLamp(f, b, queue.IsOrder(f, b))
			}
		}
	}
}
