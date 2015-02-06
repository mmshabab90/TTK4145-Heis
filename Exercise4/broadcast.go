package main

import(
	"log"
	"net"
	"time"
)

func broadcastUdp(addr string){
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {log.Fatal(err)}

	udpBroadcast, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {log.Fatal(err)}

	defer udpBroadcast.Close()

	for{
		udpBroadcast.Write([]byte("Not master"))
		time.Sleep(100*time.Millisecond)
	}
}

func listenUdp(port string, ipListChannel chan []string){
	udpAddr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {log.Fatal(err)}

	udpListen, err := net.ListenUDP("udp", udpAddr)
	if err != nil {log.Fatal(err)}

	defer udpListen.Close()

	ipList := make([]string ,0)
	var buffer[1024]byte

	timer := make(chan bool, 1)
	timeout := false

	go timerout(timer)

	for{
		_, ipAddr, err := udpListen.ReadFromUDP(buffer[:])
		if err != nil {log.Fatal(err)} 

		if(!ipInList(ipAddr.String(), ipList)){
			ipList = append(ipList, ipAddr.String())
		}
		//log.Println(ipList)

		select{
		case <-timer:
			timeout = true;
		default:
			break;
		}
		if timeout{break;}
			
		time.Sleep(100*time.Millisecond)
	}
	ipListChannel <- ipList
}

func ipInList(ipAddr string, ipList []string) bool {
    for _, b := range ipList {
        if b == ipAddr {
            return true
        }
    }
    return false
}

func timerout(timer chan bool){
	time.Sleep(5*time.Second)
	timer <- true
}

func main(){
	doneChannel := make(chan bool, 1);
	ipListChannel := make(chan []string, 1)

	port := ":20014"
	broadcastAddr := "129.241.187.255:20014"

	go broadcastUdp(broadcastAddr)
	go listenUdp(port, ipListChannel)

	log.Println(<-ipListChannel)
	<-doneChannel
}