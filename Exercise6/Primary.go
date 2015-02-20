package main 

import(
	"log"
	"net"
	"time"
)

func broadcastUdp() {
	udpAddr, err := net.ResolveUDPAddr("udp","129.241.187.255:20014")
	if err != nil {log.Fatal(err)}

	udpBroadcast, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {log.Fatal(err)}

	defer udpBroadcast.Close()
	
	msg := make([]byte, 1)
	
	for i := 1;; i++{
		log.Println(i)
		msg = byte(i);
		udpBroadcast.Write(msg)
		time.Sleep(50*time.Millisecond)
	}
}

func main() {
	broadcastUdp()
}


