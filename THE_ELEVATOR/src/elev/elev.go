/*
 *  This file is a golang port of elev.c
 */

package elev

import (
	"log"
)

var _ = log.Fatal // For debugging; delete when done.

const nButtons = 3
const nFloors = 4

type Elev_button_type_t int

const (
	ButtonCallUp Elev_button_type_t = iota
	ButtonCallDown
	ButtonCommand
)

type Elev_motor_direction_t int

const (
	DirnDown Elev_motor_direction_t = iota - 1
	DirnStop
	DirnUp
)

var lamp_channel_matrix = [nFloors][nButtons]int{
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var button_channel_matrix = [nFloors][nButtons]int{
	{BUTTON_UP1, BUTTON_DOWN1, ButtonCommand1},
	{BUTTON_UP2, BUTTON_DOWN2, ButtonCommand2},
	{BUTTON_UP3, BUTTON_DOWN3, ButtonCommand3},
	{BUTTON_UP4, BUTTON_DOWN4, ButtonCommand4},
}

func Init() int {
	// Init hardware
	if !Io_init() {
		return 0
	}

	// Zero all floor button lamps
	for f := 0; f < nFloors; f++ {
		if f != 0 {
			SetButtonLamp(f, ButtonCallDown, false)
		}
		if f != nFloors-1 {
			SetButtonLamp(f, ButtonCallUp, false)
		}
		SetButtonLamp(i, ButtonCommand, false)
	}

	// Clear stop lamp, door open lamp,
	// and set floor indicator to ground floor.
	SetStopLamp(false)
	SetDoorOpenLamp(false)
	SetFloorIndicator(0)

	// Return success.
	return 1
}

func SetMotorDirection(dirn Elev_motor_direction_t) {
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
	if floor < 0 || floor >= nFloors {
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

func GetButton(floor int, button Elev_button_type_t) bool {
	if floor < 0 || floor >= nFloors {
		log.Printf("Error: Floor %d out of range!\n", floor)
		return false
	}
	if button == ButtonCallUp && floor == nFloors-1 {
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

func SetButtonLamp(floor int, value bool, button Elev_button_type_t) {
	if floor < 0 || floor >= nFloors {
		log.Printf("Error: Floor %d out of range!\n", floor)
		return false
	}
	if button == ButtonCallUp && floor == nFloors-1 {
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

	if value {
		Io_set_bit(lamp_channel_matrix[floor][button])
	} else {
		Io_clear_bit(lamp_channel_matrix[floor][button])
	}
}
