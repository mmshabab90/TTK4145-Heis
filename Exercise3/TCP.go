package main

import(
	"net"
	"log"
	"time"
)

func main() {

	raddr, err := net.ResolveTCPAddr("tcp", "129.241.187.136:33546")
	if err != nil {log.Fatal(err)}

	socket, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {log.Fatal(err)}

	var buffer[64] byte
	var msg = "I am the walrus\x00"

	for {
		_, err := socket.Read(buffer[:])
		if err != nil {log.Fatal(err)}
		log.Println(string(buffer[:]))
		
		_, err = socket.Write([]byte(msg))
		if err != nil {log.Fatal(err)}

		time.Sleep(500*time.Millisecond)
	}
}
