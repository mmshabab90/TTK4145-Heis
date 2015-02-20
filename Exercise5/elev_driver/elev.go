package elev_driver

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
	DIRN_DOWN Elev_motor_direction_t = -1
	DIRN_STOP Elev_motor_direction_t  = 0
	DIRN_UP Elev_motor_direction_t = 1
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

func Elev_init() int {
	// Init hardware
	if !Io_init() {
		return 0
	}

	// Zero all floor button lamps
	for i := 0; i < N_FLOORS; i++ {
		if i != 0 {
			Elev_set_button_lamp(BUTTON_CALL_DOWN, i, false)
		}
		if i != N_FLOORS -1 {
			Elev_set_button_lamp(BUTTON_CALL_UP, i, false)
		}
		Elev_set_button_lamp(BUTTON_COMMAND, i, false)
	}

	// Clear stop lamp, door open lamp,
	// and set floor indicator to ground floor.
	Elev_set_stop_lamp(false)
	Elev_set_door_open_lamp(false)
	Elev_set_floor_indicator(0)

	// Return success.
	return 1
}

func Elev_set_motor_direction(dirn Elev_motor_direction_t) {
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

func Elev_set_door_open_lamp(value bool) {
	if value {
		Io_set_bit(LIGHT_DOOR_OPEN)
	} else {
		Io_clear_bit(LIGHT_DOOR_OPEN)
	}
}

func Elev_get_obstruction_signal() bool {
	return Io_read_bit(OBSTRUCTION)
}

func Elev_get_stop_signal() bool {
	return Io_read_bit(STOP)
}

func Elev_set_stop_lamp(value bool) {
	if value {
		Io_set_bit(LIGHT_STOP)
	} else {
		Io_clear_bit(LIGHT_STOP)
	}
}

func Elev_get_floor_sensor_signal() int {
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

func Elev_set_floor_indicator(floor int) {
	if floor < 0 {log.Fatalf("Floor number %s is negative!", floor)}
	if floor >= N_FLOORS {log.Fatalf("Floor number %s is above top floor!", floor)}

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

func Elev_get_button_signal(button Elev_button_type_t, floor int) bool {
	if floor < 0 {log.Fatalf("Floor number %s is negative!", floor)}
	if floor >= N_FLOORS {log.Fatalf("Floor number %s is above top floor!", floor)}
	if button == BUTTON_CALL_UP && floor == N_FLOORS - 1 {log.Fatal("Button up from top floor does not exist!")}
	if button == BUTTON_CALL_DOWN && floor == 0 {log.Fatal("Button down from ground floor does not exist!")}
	if button != BUTTON_CALL_UP && button != BUTTON_CALL_DOWN && button != BUTTON_COMMAND {log.Fatalf("Invalid button %s", button)}

	if (Io_read_bit(button_channel_matrix[floor][button])) {
		return true
	} else {
		return false
	}
}

func Elev_set_button_lamp(button Elev_button_type_t, floor int, value bool) {
	if floor < 0 {log.Fatalf("Floor number %s is negative!", floor)}
	if floor >= N_FLOORS {log.Fatalf("Floor number %s is above top floor!", floor)}
	if button == BUTTON_CALL_UP && floor == N_FLOORS - 1 {log.Fatal("Button up from top floor does not exist!")}
	if button == BUTTON_CALL_DOWN && floor == 0 {log.Fatal("Button down from ground floor does not exist!")}
	if button != BUTTON_CALL_UP && button != BUTTON_CALL_DOWN && button != BUTTON_COMMAND {log.Fatalf("Invalid button %s", button)}

	if value {
		Io_set_bit(lamp_channel_matrix[floor][button])
	} else {
		Io_clear_bit(lamp_channel_matrix[floor][button])
	}
}
