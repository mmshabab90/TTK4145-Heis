package main

import(
	"./network"
	"log"
)

func main() {
	network.NetworkInit()
	go  network.ListenForLiveElevators()
	
	for {
		select {
		case connection := <- network.ConnectionTimerChan:
			log.Println(connection.Addr, "is dead")
		}
	}
}
