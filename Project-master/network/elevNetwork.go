package main

import(
	"log"
	"net"
	"time"
)

func Broadcast(addr string, msg string){
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {log.Fatal(err)}

	udpBroadcast, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {log.Fatal(err)}

	defer udpBroadcast.Close()

	for{
		udpBroadcast.Write([]byte(msg))
		time.Sleep(100*time.Millisecond)
	}
}

func Listen(port string){
	udpAddr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {log.Fatal(err)}

	udpListen, err := net.ListenUDP("udp", udpAddr)
	if err != nil {log.Fatal(err)}

	defer udpListen.Close()

	ipList := make([]string, 0)
	var buffer[1024]byte

	for{
		_, ipAddr, err := udpListen.ReadFromUDP(buffer[:])
		if err != nil {log.Fatal(err)} 

		if(!ipInList(ipAddr.String(), ipList)){
			ipList = append(ipList, ipAddr.String())
			//go timer(), the timer most be connected to the spesific ipAddr somehow
			
		} else {
			//resetTimer, this most be connected to the right timer function somehow
		}	
		time.Sleep(100*time.Millisecond)
	}
}

func ipInList(ipAddr string, ipList []string) bool {
    for _, b := range ipList {
        if b == ipAddr {
            return true
        }
    }
    return false
}



