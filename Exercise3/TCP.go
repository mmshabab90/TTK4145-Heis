package main

import(
	"net"
	"log"
	"time"
)

//forferdelig funksjonsnavn
func connect_to_server(socket *net.TCPConn, doneChannel chan bool) {
	var buffer[1024] byte
	var msg = "I am the walrus\x00"
	for {
		//read from server
		_, err := socket.Read(buffer[:])
		if err != nil {log.Fatal(err)}
		log.Println(string(buffer[:]))
		
		//write to server
		_, err = socket.Write([]byte(msg))
		if err != nil {log.Fatal(err)}

		//delay
		time.Sleep(2000*time.Millisecond)
	}
	doneChannel <- true
	
}

//forferdelig funksjonsnavn
func accept_connection(socket_connect *net.TCPConn, doneChannel chan bool) {
	//new message
	var buffer[1024] byte
	var msg = "I am the not\x00"
	for {
		//read from server
		_, err := socket_connect.Read(buffer[:])
		if err != nil {log.Fatal(err)}
		log.Println(string(buffer[:]))
		
		//write to server
		_, err = socket_connect.Write([]byte(msg))
		if err != nil {log.Fatal(err)}

		//delay
		time.Sleep(2000*time.Millisecond)
	}
	doneChannel <- true
}

func main() {
	//TCP-setup
	raddr, err := net.ResolveTCPAddr("tcp", "129.241.187.136:33546")
	if err != nil {log.Fatal(err)}

	socket, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {log.Fatal(err)}
	
	//declare variables
	var msg = "Connect to: 129.241.187.142:20014\x00"
	
	//make listener
	laddr, err := net.ResolveTCPAddr("tcp", "129.241.187.142:20014")
	if err != nil {log.Fatal(err)}
	
	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {log.Fatal(err)}
	
	//server connect back
	_, err = socket.Write([]byte(msg))
	if err != nil {log.Fatal(err)}
	
	socket_connect, err := listener.AcceptTCP()
	if err != nil {log.Fatal(err)}
	
	doneChannel := make(chan bool, 1);
	
	go connect_to_server(socket, doneChannel)
	go accept_connection(socket_connect, doneChannel)
	
	<- doneChannel
	
}
