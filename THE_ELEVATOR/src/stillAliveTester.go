package main

import (
	"./network"
)

func main() {
	network.NetworkInit()
	go network.StillAliveBroadcast()
	
	neverReturn := make (chan int)
	<-neverReturn
}
