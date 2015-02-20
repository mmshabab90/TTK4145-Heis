package main

import(
	"log"
	"net"
	"time"
	"encoding/binary"
	"os/exec"
)

func primary(start int, udpBroadcast *net.UDPConn){
	
	msg := make([]byte, 1)
	
	for i := start;; i++{
		log.Println(i)
		msg[0] = byte(i);
		udpBroadcast.Write(msg)
		time.Sleep(100*time.Millisecond)
	}
}

func backup(udpListen *net.UDPConn) int{
	time.Sleep(1*time.Second)
	listenChan := make(chan int, 1); 
	backupvalue := 0
	//primaryDead := false
	//primaryDeadChan := make(chan bool, 1)
	buffer := make([]byte, 1024)
	for {
		select {
			case <-time.After(1*time.Second):
				log.Println("The primary is dead, long live the primary")
				return backupvalue
			case backupvalue = <-listenChan:
				break
		}
		_, _, err := udpListen.ReadFromUDP(buffer[:])
		if err != nil {log.Fatal(err)} 
		
		listenChan <- int(binary.LittleEndian.Uint64(buffer)) //convert an bytearray to int
		time.Sleep(100*time.Millisecond)
	}
	
	
}

/*func listen(listenChan chan int, udpListen *net.UDPConn, primaryDeadChan chan bool) {

	buffer := make([]byte, 1024)

	for {
		select{
			case <-primaryDeadChan:
				return
			default:
				break
		}
		_, _, err := udpListen.ReadFromUDP(buffer[:])
		if err != nil {log.Fatal(err)} 
		
		listenChan <- int(binary.LittleEndian.Uint64(buffer)) //convert an bytearray to int
		time.Sleep(100*time.Millisecond)
	}
	
}*/

func main() {
	
	udpAddr, err := net.ResolveUDPAddr("udp", ":20014")
	if err != nil {log.Fatal(err)}

	udpListen, err := net.ListenUDP("udp", udpAddr)
	if err != nil {log.Fatal(err)}
	
	backupvalue := backup(udpListen)
	
	udpListen.Close()
	
	newBackup := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run backup.go")
	err = newBackup.Run()
	if err != nil {log.Fatal(err)}
	
	udpAddr, err = net.ResolveUDPAddr("udp","129.241.187.142:20014")
	if err != nil {log.Fatal(err)}

	udpBroadcast, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {log.Fatal(err)}
	
	
	primary(backupvalue, udpBroadcast)
	
	
	
	
}
