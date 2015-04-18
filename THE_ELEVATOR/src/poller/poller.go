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
	"../defs"
	"errors"
)

var _ = log.Println
var _ = fmt.Println

type keypress struct {
	button int
	floor  int
}

var connectionMap = make(map[string] network.UdpConnection)
var connectionDeadChan	 = make(chan network.UdpConnection)
const resetTime = 1*time.Second

var costChan = make(chan network.Message)
type reply struct {
	cost int
	lift string
}
type order struct {
	floor int
	button int
}

func liftAssigner() {
	// collect cost values from all lifts
	// decide which lift gets the order when all lifts
	// in alive-list have answered or after a timeout
	// either send the decision on network or pray that all
	// lifts make the same choice every time

	// spawn a goroutine for each order to be assigned?
	go func() {
		assignmentQueue := make(map[order][]reply)
		// assignmentQueue := make([]order)
		// assignmentQueue := [defs.NumFloors][defs.NumButtons][]reply
		for {
			select {
			case message := <- costChan:
				order, reply = parse(message)
				if value, exist := assignmentQueue[order]; exist {
					// sjekk om val fins i listen
					// append om ikke
					found := false
					for _, e := range value {
						if e == reply {
							found = true
						}
					}
					if !found {
						assignmentQueue[order] = append(assignmentQueue[order], reply)
					}
				} else {
					assignmentQueue[order] = []reply{reply}
				}






				if !isOrderInAss(assignmentQueue, message) {
					// add message to assQue:
					newOrder := order{floor:message.floor, button:message.button}
					assignmentQueue = append(assignmentQueue, newOrder)
					index := findIndexInAss(assignmentQueue, newOrder)
					go waitForCosts(*assignmentQueue[index])
				}
				index := findIndexInAss(assignmentQueue, newOrder)
			default:
				// do nothing
			}
		}
	}()
}

func parse(m network.Message) order, reply {
	return order{floor:m.Floor, button:m.Button}, reply{cost:m.Cost, lift:m.Addr}
}

func isOrderInAss(que []order, message network.Message) bool {
	for _, e := range que {
        if e.floor == message.floor && e.button == message.button {
        	return true
        }
    }
    return false
}

func findIndexInAss(que []order, ord order) int {
	for i, e := range que {
		if e == ord {
			return i
		}
	}
	return -1
}

func waitForCosts(ord *order) {

}






func Init() {
	if err := hw.Init(); err != nil {
		log.Fatal(err)
	}
	fsm.Init()
	network.Init()

	run()
}

func run() {
	buttonChan := pollButtons()
	floorChan := pollFloors()

	for {
		select {
		case keypress := <-buttonChan:
			fsm.EventButtonPressed(keypress.floor, keypress.button)
		case floor := <-floorChan:
			fsm.EventFloorReached(floor)
		case <-fsm.DoorTimeoutChan:
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
		var buttonState [defs.NumFloors][defs.NumButtons]bool

		for {
			for f := 0; f < defs.NumFloors; f++ {
				for b := 0; b < defs.NumButtons; b++ {
					if (f == 0 && b == defs.ButtonCallDown) ||
						(f == defs.NumFloors-1 && b == defs.ButtonCallUp) {
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
				connection.Timer.Reset(resetTime)
				fmt.Println("timer reset for IP: ")
				fmt.Println(message.Addr)
			} else {
				newConnection := network.UdpConnection{message.Addr, time.NewTimer(resetTime)}
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
			costMessage := &network.Message{
				Kind: network.Cost,
				Floor: message.Floor,
				Button: message.Button,
				Cost: cost}
			network.Send(costMessage)
		case network.CompleteOrder:
			// remove from queues
			queue.RemoveSharedOrder(message.Floor, message.Button)
			// prob more to do here
		case network.Cost:
			costChan <- message
	}
}

func connectionTimer(connection *network.UdpConnection) {
	for {
		<- connection.Timer.C
		connectionDeadChan <- *connection
	}
}
