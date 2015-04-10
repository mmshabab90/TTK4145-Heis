package network

import (
	"fmt"
	"./udp"
)

var Send_ch = make (chan udp.Udp_message)
var Receive_ch = make (chan udp.Udp_message)

func Print_udp_message(msg udp.Udp_message){
	fmt.Printf("msg:  \n \t raddr = %s \n \t data = %s \n \t length = %v \n", msg.Raddr, msg.Data, msg.Length)
}

func NetworkInit (){
	err := udp.Udp_init(20001, 20014, 1024, Send_ch, Receive_ch)	

	if (err != nil){
		fmt.Print("err = %s \n", err)
	}
}


