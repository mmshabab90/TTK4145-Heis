package main

import (
	"./elev_driver"
)

func main(){
	elev_driver.Elev_init()
	elev_driver.Elev_set_motor_direction(elev_driver.DIRN_UP)
}