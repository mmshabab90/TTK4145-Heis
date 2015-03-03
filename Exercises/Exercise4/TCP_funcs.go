package main

import(
	"log"
	"net"
	"time"
)
	
func setUpMaster(addr string){
	//make listener
	masterAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {log.Fatal(err)}
	
	listener, err := net.ListenTCP("tcp", masterAddr)
	if err != nil {log.Fatal(err)}

	for{
		conn, err := listener.AcceptTCP()
		if err != nil {log.Fatal(err)}
		if conn != nil{
			log.Println("Connection established")
		}
	}
} 

func connectToMaster(masterAddr string){
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {log.Fatal(err)}
	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {log.Fatal(err)}

	log.Println("Connection established")
}

