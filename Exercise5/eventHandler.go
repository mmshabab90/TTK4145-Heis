package driver

import (
	"log"
)

type State_t int

const(
	IDLE State_t = iota
	RUNNING
	DOOROPEN
)

func EventHandler(){

}

func init(){
	Elev_init()
	Elev_set_speed(-300)
	for{
		floor := Elev_get_floor_sensor_signal()
		if floor != -1{
		
		}
	}
}