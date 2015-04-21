/*
 *  This file is a golang port of elev.c from the hand out driver
 *  https://github.com/TTK4145/Project
 */

package hw

import (
	"../defs"
	"errors"
	"log"
)

var lampChannelMatrix = [defs.NumStoreys][defs.NumButtons]int{
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}
var buttonChannelMatrix = [defs.NumStoreys][defs.NumButtons]int{
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

	// Zero all storey button lamps
	for f := 0; f < defs.NumStoreys; f++ {
		if f != 0 {
			SetButtonLamp(f, defs.ButtonDown, false)
		}
		if f != defs.NumStoreys-1 {
			SetButtonLamp(f, defs.ButtonUp, false)
		}
		SetButtonLamp(f, defs.ButtonCommand, false)
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

func Storey() int {
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

func SetStoreyLamp(storey int) {
	if storey < 0 || storey >= defs.NumStoreys {
		log.Printf("Error: Storey %d out of range!\n", storey)
		log.Println("No storey indicator will be set.")
		return
	}

	// Binary encoding. One light must always be on.
	if storey&0x02 > 0 {
		ioSetBit(LIGHT_FLOOR_IND1)
	} else {
		ioClearBit(LIGHT_FLOOR_IND1)
	}

	if storey&0x01 > 0 {
		ioSetBit(LIGHT_FLOOR_IND2)
	} else {
		ioClearBit(LIGHT_FLOOR_IND2)
	}
}

func ReadButton(storey int, button int) bool {
	if storey < 0 || storey >= defs.NumStoreys {
		log.Printf("Error: Storey %d out of range!\n", storey)
		return false
	}
	if button < 0 || button >= defs.NumButtons {
		log.Printf("Error: Button %d out of range!\n", button)
		return false
	}
	if button == defs.ButtonUp && storey == defs.NumStoreys-1 {
		log.Println("Button up from top storey does not exist!")
		return false
	}
	if button == defs.ButtonDown && storey == 0 {
		log.Println("Button down from ground storey does not exist!")
		return false
	}

	if ioReadBit(buttonChannelMatrix[storey][button]) {
		return true
	} else {
		return false
	}
}

func SetButtonLamp(storey int, button int, value bool) {
	if storey < 0 || storey >= defs.NumStoreys {
		log.Printf("Error: Storey %d out of range!\n", storey)
		return
	}
	if button == defs.ButtonUp && storey == defs.NumStoreys-1 {
		log.Println("Button up from top storey does not exist!")
		return
	}
	if button == defs.ButtonDown && storey == 0 {
		log.Println("Button down from ground storey does not exist!")
		return
	}
	if button != defs.ButtonUp &&
		button != defs.ButtonDown &&
		button != defs.ButtonCommand {
		log.Printf("Invalid button %d\n", button)
		return
	}

	if value {
		ioSetBit(lampChannelMatrix[storey][button])
	} else {
		ioClearBit(lampChannelMatrix[storey][button])
	}
}

func MoveToDefinedState() int {
	SetMotorDirection(defs.DirDown)
	storey := Storey()
	for storey == -1 {
		storey = Storey()
	}
	SetMotorDirection(defs.DirStop)
	SetStoreyLamp(storey)
	return storey
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
