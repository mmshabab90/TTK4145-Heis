//testprogram for elev.go

package main

import (
	"./driver"
	"log"
)

func main(){
	driver.Elev_init()
	driver.Elev_set_motor_direction(driver.DIRN_UP)
	
	log.Println(driver.DIRN_STOP)
	
	for{
		floor := driver.Elev_get_floor_sensor_signal()
		stop := driver.Elev_get_stop_signal()
		obstr := driver.Elev_get_obstruction_signal()
		if floor == 3 {
			driver.Elev_set_floor_indicator(floor)
			driver.Elev_set_motor_direction(driver.DIRN_DOWN)
			for j:= driver.BUTTON_CALL_DOWN; j < 3; j++ {
				driver.Elev_set_button_lamp(j, floor, false)
			}
		} else if floor == 0 {
			driver.Elev_set_floor_indicator(floor)
			driver.Elev_set_motor_direction(driver.DIRN_UP)
			for j:= driver.BUTTON_CALL_UP; j < 3; j++ {
				if j == driver.BUTTON_CALL_DOWN{continue}
				driver.Elev_set_button_lamp(j, floor, false)
			}
		} else if floor != -1 {
			driver.Elev_set_floor_indicator(floor)
			for j:= driver.BUTTON_CALL_UP; j < 3; j++ {
				driver.Elev_set_button_lamp(j, floor, false)
			}
		}
		for j:= driver.BUTTON_CALL_UP; j < 3; j++ {
			for i:=0; i < 4; i++{
				if i == 3 && j == driver.BUTTON_CALL_UP{continue}
				if i == 0 && j == driver.BUTTON_CALL_DOWN{continue}
				button := driver.Elev_get_button_signal(j, i)
				if button {
					driver.Elev_set_button_lamp(j, i, button)
				}
			}
		}
		driver.Elev_set_stop_lamp(stop)
		driver.Elev_set_door_open_lamp(obstr)
		if(stop){
			driver.Elev_set_motor_direction(driver.DIRN_STOP)
		}
	}
}
