package main

import(
	"log"
	"net"
	"time"
	"encoding/binary"
)

func primary(){
	udpAddr, err := net.ResolveUDPAddr("udp","129.241.187.255:20014")
	if err != nil {log.Fatal(err)}

	udpBroadcast, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {log.Fatal(err)}

	defer udpBroadcast.Close()
	
	msg := make([]byte, 1)
	
	for i := 1;; i++{
		log.Println(i)
		msg[0] = byte(i);
		udpBroadcast.Write(msg)
		time.Sleep(200*time.Millisecond)
	}
}

func backup(listenChan chan int){
	var backupvalue int
	primaryDead := false
	go listen(listenChan)
	for {
		select {
		case backupvalue = <-listenChan:
			break
		case <-time.After(1*time.Second):
			primaryDead = true
		}
		if primaryDead {
			log.Println("The primary is  dead, long live the primary!")
			break
		}	
	}
	log.Println(backupvalue);
}

func listen(listenChan chan int) bool {
	udpAddr, err := net.ResolveUDPAddr("udp", ":20014")
	if err != nil {log.Fatal(err)}

	udpListen, err := net.ListenUDP("udp", udpAddr)
	if err != nil {log.Fatal(err)}

	defer udpListen.Close()

	buffer := make([]byte, 1024)

	for {
		_, _, err := udpListen.ReadFromUDP(buffer[:])
		if err != nil {log.Fatal(err)} 
		
		listenChan <- int(binary.LittleEndian.Uint64(buffer))
		time.Sleep(100*time.Millisecond)
	}
}

func main() {
	listenChan := make(chan int, 1);
	backup(listenChan)
}
