package main

import (
	"fmt"
	"./network"
	"time"
	"./network/udp"
)

func main() {
	network.NetworkInit()
	for {
		sndMsg := udp.Udp_message{Raddr:"broadcast", Data:"Hello World", Length:11}
		network.Send_ch <- sndMsg
		fmt.Println("Msg sent")
		network.Print_udp_message(sndMsg)
		time.Sleep(1*time.Second)
	}
}
