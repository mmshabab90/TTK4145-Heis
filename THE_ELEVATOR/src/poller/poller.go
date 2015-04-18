package poller

import (
	"../hw"
	"../fsm"
	"../cost"
	"../network"
	"log"
	"time"
	"fmt"
	"../queue"
)

var _ = log.Println
var _ = fmt.Println

type keypress struct {
	button int
	floor  int
}

var connectionMap = make(map[string] network.UdpConnection)
var connectionDeadChan	 = make(chan network.UdpConnection)

func Run() {
	buttonChan := pollButtons()
	floorChan := pollFloors()

	for {
		select {
		case keypress := <-buttonChan:
			fsm.EventButtonPressed(keypress.floor, keypress.button)
		case floor := <-floorChan:
			fsm.EventFloorReached(floor)
		case <-fsm.DoorTimeout:
			fsm.EventDoorTimeout()
		case udpMessage := <-network.ReceiveChan:
			handleMessage(network.ParseMessage(udpMessage))
		case connection := <- connectionDeadChan:
			fmt.Printf("Connection with IP %s is dead\n", connection.Addr)
			delete(connectionMap, connection.Addr) //delete dead connection from map
			queue.ReassignOrders(connection.Addr)
			//for key, _ := range connectionMap {fmt.Println(key)}
		}
	}
}

func pollButtons() <-chan keypress {
	c := make(chan keypress)

	go func() {
		var buttonState [hw.NumFloors][hw.NumButtons]bool

		for {
			for f := 0; f < hw.NumFloors; f++ {
				for b := 0; b < hw.NumButtons; b++ {
					if (f == 0 && b == hw.ButtonCallDown) ||
						(f == hw.NumFloors-1 && b == hw.ButtonCallUp) {
						continue
					}
					if hw.ReadButton(f, b) {
						if !buttonState[f][b] {
							c <- keypress{button: b, floor: f}
						}
						buttonState[f][b] = true
					} else {
						buttonState[f][b] = false
					}
				}
			}
			time.Sleep(time.Millisecond)
		}
	}()

	return c
}

func pollFloors() <-chan int {
	c := make(chan int)

	go func() {
		oldFloor := hw.GetFloor()

		for {
			newFloor := hw.GetFloor()
			if newFloor != oldFloor && newFloor != -1 {
				c <- newFloor
			}
			oldFloor = newFloor
			time.Sleep(time.Millisecond)
		}
	}()

	return c
}

func handleMessage(message network.Message) {
	switch message.Kind {
		case network.Alive:
			if connection, exist := connectionMap[message.Addr]; exist {
				connection.Timer.Reset(1*time.Second)
				fmt.Println("timer reset for IP: ")
				fmt.Println(message.Addr)
			} else {
				newConnection := network.UdpConnection{message.Addr, time.NewTimer(1*time.Second)}
				connectionMap[message.Addr] = newConnection
				fmt.Println("New connection, with IP: ")
				fmt.Println(message.Addr)
				go connectionTimer(&newConnection)
			}
		case network.NewOrder:
			cost, err := cost.CalculateCost(message.Floor, message.Button, fsm.GetFloor(), fsm.GetDirection())
			if err != nil {
				log.Println(err)
			}
			costMessage := network.Message{Kind: network.Cost, Floor: message.Floor, Button: message.Button,	Cost: cost}
			network.Send(costMessage)
		case network.CompleteOrder:
			// remove from queues
			queue.RemoveSharedOrder(message.Floor, message.Button)
			// prob more to do here
		case network.Cost:
			// notify lift assignment routine
	}
}

func connectionTimer(connection *network.UdpConnection) {
	for {
		<- connection.Timer.C
		connectionDeadChan <- *connection
	}
}
