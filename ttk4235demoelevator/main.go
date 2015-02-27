//testprogram for elev.go

package main

import (
	"./elev_driver"
	"log"
)

func main(){
	elev_driver.Elev_init()
	elev_driver.Elev_set_motor_direction(elev_driver.DIRN_UP)
	
	log.Println(elev_driver.DIRN_STOP)
	
	for{
		floor := elev_driver.Elev_get_floor_sensor_signal()
		stop := elev_driver.Elev_get_stop_signal()
		obstr := elev_driver.Elev_get_obstruction_signal()
		if floor == 3 {
			elev_driver.Elev_set_floor_indicator(floor)
			elev_driver.Elev_set_motor_direction(elev_driver.DIRN_DOWN)
			for j:= elev_driver.BUTTON_CALL_DOWN; j < 3; j++ {
				elev_driver.Elev_set_button_lamp(j, floor, false)
			}
		} else if floor == 0 {
			elev_driver.Elev_set_floor_indicator(floor)
			elev_driver.Elev_set_motor_direction(elev_driver.DIRN_UP)
			for j:= elev_driver.BUTTON_CALL_UP; j < 3; j++ {
				if j == elev_driver.BUTTON_CALL_DOWN{continue}
				elev_driver.Elev_set_button_lamp(j, floor, false)
			}
		} else if floor != -1 {
			elev_driver.Elev_set_floor_indicator(floor)
			for j:= elev_driver.BUTTON_CALL_UP; j < 3; j++ {
				elev_driver.Elev_set_button_lamp(j, floor, false)
			}
		}
		for j:= elev_driver.BUTTON_CALL_UP; j < 3; j++ {
			for i:=0; i < 4; i++{
				if i == 3 && j == elev_driver.BUTTON_CALL_UP{continue}
				if i == 0 && j == elev_driver.BUTTON_CALL_DOWN{continue}
				button := elev_driver.Elev_get_button_signal(j, i)
				if button {
					elev_driver.Elev_set_button_lamp(j, i, button)
				}
			}
		}
		elev_driver.Elev_set_stop_lamp(stop)
		elev_driver.Elev_set_door_open_lamp(obstr)
		if(stop){
			elev_driver.Elev_set_motor_direction(elev_driver.DIRN_STOP)
		}
	}
}
