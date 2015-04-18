/*
 *  This file is a golang port of elev.c from the hand out driver
 *  https://github.com/TTK4145/Project
 */

package hw

import (
	"log"
	"errors"
	"../defs"
)

var lampChannelMatrix = [NumFloors][NumButtons]int{
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}
var buttonChannelMatrix = [NumFloors][NumButtons]int{
	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}

func Init() error {
	// Init hardware
	if !ioInit() {
		return errors.New("Hardware driver: ioInit() failed!")
	}

	// Zero all floor button lamps
	for f := 0; f < NumFloors; f++ {
		if f != 0 {
			SetButtonLamp(f, ButtonCallDown, false)
		}
		if f != NumFloors-1 {
			SetButtonLamp(f, ButtonCallUp, false)
		}
		SetButtonLamp(f, ButtonCommand, false)
	}

	SetStopLamp(false)
	SetDoorOpenLamp(false)

	MoveToDefinedState()

	return nil
}

func SetMotorDirection(dirn int) {
	if dirn == 0 {
		ioWriteAnalog(MOTOR, 0)
	} else if dirn > 0 {
		ioClearBit(MOTORDIR)
		ioWriteAnalog(MOTOR, 2800)
	} else if dirn < 0 {
		ioSetBit(MOTORDIR)
		ioWriteAnalog(MOTOR, 2800)
	}
}

func SetDoorOpenLamp(value bool) {
	if value {
		ioSetBit(LIGHT_DOOR_OPEN)
	} else {
		ioClearBit(LIGHT_DOOR_OPEN)
	}
}

func GetFloor() int {
	if ioReadBit(SENSOR_FLOOR1) {
		return 0
	} else if ioReadBit(SENSOR_FLOOR2) {
		return 1
	} else if ioReadBit(SENSOR_FLOOR3) {
		return 2
	} else if ioReadBit(SENSOR_FLOOR4) {
		return 3
	} else {
		return -1
	}
}

func SetFloorLamp(floor int) {
	if floor < 0 || floor >= NumFloors {
		log.Printf("Error: Floor %d out of range!\n", floor)
		log.Println("No floor indicator will be set.")
		return
	}

	// Binary encoding. One light must always be on.
	if floor&0x02 > 0 {
		ioSetBit(LIGHT_FLOOR_IND1)
	} else {
		ioClearBit(LIGHT_FLOOR_IND1)
	}

	if floor&0x01 > 0 {
		ioSetBit(LIGHT_FLOOR_IND2)
	} else {
		ioClearBit(LIGHT_FLOOR_IND2)
	}
}

func ReadButton(floor int, button int) bool {
	if floor < 0 || floor >= NumFloors {
		log.Printf("Error: Floor %d out of range!\n", floor)
		return false
	}
	if button < 0 || button >= NumButtons {
		log.Printf("Error: Button %d out of range!\n", button)
		return false
	}
	if button == ButtonCallUp && floor == NumFloors-1 {
		log.Println("Button up from top floor does not exist!")
		return false
	}
	if button == ButtonCallDown && floor == 0 {
		log.Println("Button down from ground floor does not exist!")
		return false
	}

	if ioReadBit(buttonChannelMatrix[floor][button]) {
		return true
	} else {
		return false
	}
}

func SetButtonLamp(floor int, button int, value bool) {
	if floor < 0 || floor >= NumFloors {
		log.Printf("Error: Floor %d out of range!\n", floor)
		return
	}
	if button == ButtonCallUp && floor == NumFloors-1 {
		log.Println("Button up from top floor does not exist!")
		return
	}
	if button == ButtonCallDown && floor == 0 {
		log.Println("Button down from ground floor does not exist!")
		return
	}
	if button != ButtonCallUp &&
		button != ButtonCallDown &&
		button != ButtonCommand {
		log.Printf("Invalid button %d\n", button)
		return
	}

	if value {
		ioSetBit(lampChannelMatrix[floor][button])
	} else {
		ioClearBit(lampChannelMatrix[floor][button])
	}
}

func MoveToDefinedState() int {
	SetMotorDirection(DirnDown)
	floor := GetFloor()
	for floor == -1 {
		floor = GetFloor()
	}
	SetMotorDirection(DirnStop)
	SetFloorLamp(floor)
	return floor
}

// Not used:

func SetStopLamp(value bool) {
	if value {
		ioSetBit(LIGHT_STOP)
	} else {
		ioClearBit(LIGHT_STOP)
	}
}

func GetObstructionSignal() bool {
	return ioReadBit(OBSTRUCTION)
}

func GetStopSignal() bool {
	return ioReadBit(STOP)
}
