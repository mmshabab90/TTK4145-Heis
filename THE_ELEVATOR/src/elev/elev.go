package elev

import(
	"log"
)

/*
 *  This file is a port of elev.c
 */

var _ = log.Fatal // For debugging; delete when done.

const N_BUTTONS = 3
const N_FLOORS  = 4

type Elev_button_type_t int
const (
	BUTTON_CALL_UP Elev_button_type_t = iota
	BUTTON_CALL_DOWN
	BUTTON_COMMAND
)

type Elev_motor_direction_t int
const (
	DIRN_DOWN Elev_motor_direction_t = iota - 1
	DIRN_STOP
	DIRN_UP
)

var lamp_channel_matrix = [N_FLOORS][N_BUTTONS] int {
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var button_channel_matrix = [N_FLOORS][N_BUTTONS] int {
	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}

func Init() int {
	// Init hardware
	if !Io_init() {
		return 0
	}

	// Zero all floor button lamps
	for i := 0; i < N_FLOORS; i++ {
		if i != 0 {
			SetButtonLamp(BUTTON_CALL_DOWN, i, false)
		}
		if i != N_FLOORS -1 {
			SetButtonLamp(BUTTON_CALL_UP, i, false)
		}
		SetButtonLamp(BUTTON_COMMAND, i, false)
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

func GetFloorSensorSignal() int {
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
	if floor < 0 || floor >= N_FLOORS {
		log.Printf("Error: Floor %d out of range!\n", floor)
		log.Println("No floor indicator will be set.")
		return
	}

	// Binary encoding. One light must always be on.
	if floor & 0x02 > 0 {
		Io_set_bit(LIGHT_FLOOR_IND1)
	} else {
		Io_clear_bit(LIGHT_FLOOR_IND1)
	}

	if floor & 0x01 > 0 {
		Io_set_bit(LIGHT_FLOOR_IND2)
	} else {
		Io_clear_bit(LIGHT_FLOOR_IND2)
	}
}

func GetButtonSignal(button Elev_button_type_t, floor int) bool {
	if floor < 0 || floor >= N_FLOORS {
		log.Printf("Error: Floor %d out of range!\n", floor)
		return false
	}
	if button == BUTTON_CALL_UP && floor == N_FLOORS - 1 {
		log.Println("Button up from top floor does not exist!")
		return false
	}
	if button == BUTTON_CALL_DOWN && floor == 0 {
		log.Println("Button down from ground floor does not exist!")
		return false
	}
	if button != BUTTON_CALL_UP &&
	button != BUTTON_CALL_DOWN &&
	button != BUTTON_COMMAND {
		log.Printf("Invalid button %d\n", button)
		return false
	}

	if (Io_read_bit(button_channel_matrix[floor][button])) {
		return true
	} else {
		return false
	}
}

func SetButtonLamp(button Elev_button_type_t, floor int, value bool) {
	if floor < 0 || floor >= N_FLOORS {
		log.Printf("Error: Floor %d out of range!\n", floor)
		return false
	}
	if button == BUTTON_CALL_UP && floor == N_FLOORS - 1 {
		log.Println("Button up from top floor does not exist!")
		return false
	}
	if button == BUTTON_CALL_DOWN && floor == 0 {
		log.Println("Button down from ground floor does not exist!")
		return false
	}
	if button != BUTTON_CALL_UP &&
	button != BUTTON_CALL_DOWN &&
	button != BUTTON_COMMAND {
		log.Printf("Invalid button %d\n", button)
		return false
	}

	if value {
		Io_set_bit(lamp_channel_matrix[floor][button])
	} else {
		Io_clear_bit(lamp_channel_matrix[floor][button])
	}
}
