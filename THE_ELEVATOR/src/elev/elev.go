/*
 *  This file is a golang port of elev.c
 */

package elev

import (
	"log"
)

var _ = log.Fatal // For debugging; delete when done.

const NumButtons = 3
const NumFloors = 4

type ButtonType int
type MotorDirnType int

const (
	ButtonCallUp int = iota
	ButtonCallDown
	ButtonCommand
)

const (
	DirnDown MotorDirnType = iota - 1
	DirnStop
	DirnUp
)

var lamp_channel_matrix = [NumFloors][NumButtons]int{
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var button_channel_matrix = [NumFloors][NumButtons]int{
	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}

func Init() bool {
	// Init hardware
	if !Io_init() {
		return false
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

	// Clear stop lamp, door open lamp,
	// and set floor indicator to ground floor.
	SetStopLamp(false)
	SetDoorOpenLamp(false)
	SetFloorIndicator(0)

	// Return success.
	return true
}

func SetMotorDirection(dirn MotorDirnType) {
	if dirn == 0 {
		Io_write_analog(MOTOR, 0)
	} else if dirn > 0 {
		Io_clear_bit(MOTORDIR)
		Io_write_analog(MOTOR, 2800)
	} else if dirn < 0 {
		Io_set_bit(MOTORDIR)
		Io_write_analog(MOTOR, 2800)
	}
}

func SetDoorOpenLamp(value bool) {
	if value {
		Io_set_bit(LIGHT_DOOR_OPEN)
	} else {
		Io_clear_bit(LIGHT_DOOR_OPEN)
	}
}

func GetObstructionSignal() bool {
	return Io_read_bit(OBSTRUCTION)
}

func GetStopSignal() bool {
	return Io_read_bit(STOP)
}

func SetStopLamp(value bool) {
	if value {
		Io_set_bit(LIGHT_STOP)
	} else {
		Io_clear_bit(LIGHT_STOP)
	}
}

func GetFloor() int {
	if Io_read_bit(SENSOR_FLOOR1) {
		return 0
	} else if Io_read_bit(SENSOR_FLOOR2) {
		return 1
	} else if Io_read_bit(SENSOR_FLOOR3) {
		return 2
	} else if Io_read_bit(SENSOR_FLOOR4) {
		return 3
	} else {
		return -1
	}
}

func SetFloorIndicator(floor int) {
	if floor < 0 || floor >= NumFloors {
		log.Printf("Error: Floor %d out of range!\n", floor)
		log.Println("No floor indicator will be set.")
		return
	}

	// Binary encoding. One light must always be on.
	if floor&0x02 > 0 {
		Io_set_bit(LIGHT_FLOOR_IND1)
	} else {
		Io_clear_bit(LIGHT_FLOOR_IND1)
	}

	if floor&0x01 > 0 {
		Io_set_bit(LIGHT_FLOOR_IND2)
	} else {
		Io_clear_bit(LIGHT_FLOOR_IND2)
	}
}

func GetButton(floor int, button int) bool {
	if floor < 0 || floor >= NumFloors {
		log.Printf("Error: Floor %d out of range!\n", floor)
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
	if button != ButtonCallUp &&
		button != ButtonCallDown &&
		button != ButtonCommand {
		log.Printf("Invalid button %d\n", button)
		return false
	}

	if Io_read_bit(button_channel_matrix[floor][button]) {
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
		Io_set_bit(lamp_channel_matrix[floor][button])
	} else {
		Io_clear_bit(lamp_channel_matrix[floor][button])
	}
}
