package main

import (
	"fmt"
	"./network"
)

func main() {
	network.NetworkInit()
	for {
		select {
			case rcvMsg := <- network.Receive_ch:
				fmt.Println("Msg received")
				network.Print_udp_message(rcvMsg)
		}
	}
}
