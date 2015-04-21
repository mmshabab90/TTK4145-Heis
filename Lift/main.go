package main

import (
	"./src/defs"
	"./src/fsm"
	"./src/hw"
	"./src/network"
	"./src/queue"
	"errors"
	"fmt"
	"log"
	"time"
)

const debugPrint = false

var _ = log.Println
var _ = fmt.Println
var _ = errors.New

type keypress struct {
	button int
	floor  int
}

var onlineLifts = make(map[string]network.UdpConnection)

var deadChan = make(chan network.UdpConnection)
var costChan = make(chan defs.Message)
var orderTimeoutChan = make(chan order)

type reply struct {
	cost int
	lift string
}
type order struct {
	floor  int
	button int
	// timeout bool
	// timer   *time.Timer
}

func main() {
	if err := hw.Init(); err != nil {
		log.Fatal(err)
	}
	fsm.Init()
	network.Init()

	liftAssigner()
	run()
}

func run() {
	buttonChan := pollButtons()
	floorChan := pollFloors()

	for {
		select {
		case keypress := <-buttonChan:
			switch keypress.button {
			case defs.ButtonCommand:
				fsm.EventInternalButtonPressed(keypress.floor, keypress.button)
			case defs.ButtonCallUp, defs.ButtonCallDown:
				fsm.EventExternalButtonPressed(keypress.floor, keypress.button)
			default:
				// maybe care about bad button here
			}

		case floor := <-floorChan:
			fsm.EventFloorReached(floor)
		case udpMessage := <-network.ReceiveChan:
			handleMessage(network.ParseMessage(udpMessage))
		case connection := <-deadChan:
			handleDeadLift(connection.Addr)
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
		oldFloor := hw.Floor()

		for {
			newFloor := hw.Floor()
			if newFloor != oldFloor && newFloor != -1 {
				c <- newFloor
			}
			oldFloor = newFloor
			time.Sleep(time.Millisecond)
		}
	}()

	return c
}

// consider moving each case into a function
func handleMessage(message defs.Message) {
	/*fmt.Println("Received")
	network.PrintMessage(message)
	fmt.Println()*/
	switch message.Kind {
	case defs.Alive:
		if connection, exist := onlineLifts[message.Addr]; exist {
			connection.Timer.Reset(defs.ResetTime)
			if debugPrint {
				fmt.Printf("Timer reset for IP %s\n", message.Addr)
			}
		} else {
			newConnection := network.UdpConnection{message.Addr, time.NewTimer(defs.ResetTime)}
			onlineLifts[message.Addr] = newConnection
			if debugPrint {
				fmt.Printf("New connection with IP %s\n", message.Addr)
			}
			go connectionTimer(&newConnection)
		}
	case defs.NewOrder:
		fmt.Printf("handleMessage(): NewOrder message: f=%d b=%d from lift %s\n",
			message.Floor+1, message.Button, message.Addr[12:15])

		cost := queue.CalculateCost(message.Floor, message.Button, fsm.Floor(), hw.Floor(), fsm.Direction())

		costMessage := defs.Message{
			Kind:   defs.Cost,
			Floor:  message.Floor,
			Button: message.Button,
			Cost:   cost}
		//fmt.Printf("handleMessage(): NewOrder sends cost message: f=%d b=%d (with cost %d) from me\n", costMessage.Floor+1, costMessage.Button, costMessage.Cost)
		defs.MessageChan <- costMessage
	case defs.CompleteOrder:
		fmt.Println("handleMessage(): CompleteOrder message")
		// remove from queues
		queue.RemoveRemoteOrdersAt(message.Floor)

		// prob more to do here
	case defs.Cost:
		fmt.Printf("handleMessage(): Cost message: f=%d b=%d with cost %d from lift %s\n", message.Floor+1, message.Button, message.Cost, message.Addr[12:15])
		costChan <- message
	}
}

func handleDeadLift(deadAddr string) {
	fmt.Printf("Connection to IP %s is dead!\n", deadAddr)
	delete(onlineLifts, deadAddr)
	queue.ReassignOrders(deadAddr)
}

func connectionTimer(connection *network.UdpConnection) {
	for { //don't think this needs to be a for-loop
		<-connection.Timer.C
		deadChan <- *connection
	}
}

// func orderTimer(newOrder *order) {
// 	<-newOrder.timer.C
// 	orderTimeoutChan <- *newOrder

// }

func liftAssigner() {
	// collect cost values from all lifts
	// decide which lift gets the order when all lifts
	// in alive-list have answered or after a timeout
	// either send the decision on network or pray that all
	// lifts make the same choice every time
	go func() {
		assignmentQueue := make(map[order][]reply)
		for {
			select {
			case message := <-costChan:
				newOrder, newReply := split(message)
				// Check if order in queue
				if value, exist := assignmentQueue[newOrder]; exist {
					// Check if lift in list of that order
					found := false
					for _, e := range value {
						if e == newReply {
							found = true
						}
					}
					// Add it if not found
					if !found {
						assignmentQueue[newOrder] = append(assignmentQueue[newOrder], newReply)
						//newOrder.timer.Reset(10 * time.Second)
					}
				} else {
					// If order not in queue at all, init order list with it
					assignmentQueue[newOrder] = []reply{newReply}
					// newOrder.timer = time.NewTimer(10 * time.Second)
					// go orderTimer(&newOrder)
				}
				evaluateLists(assignmentQueue)
			/*case newOrder := <-orderTimeoutChan:
				fmt.Printf("\n\n ORDER TIMED OUT!\n")
				// newOrder.timeout = true
				evaluateLists(assignmentQueue)*/
			}
		}
	}()
}

func split(m defs.Message) (order, reply) {
	return order{floor: m.Floor, button: m.Button}, reply{cost: m.Cost, lift: m.Addr}
}

// evaluateLists goes through the map of orders with associated costs, checks
// if any orders have received answers from all live lifts, and finds the
// the best candidate for all such orders. The best candidate is added to the
// shared queue.
// This is very cryptic and ungood.
func evaluateLists(que map[order][]reply) {
	// Loop thru all lists
	fmt.Printf("Lists: ")
	fmt.Println(que)
	for key, replyList := range que {
		// Check if the list is complete
		if len(replyList) == len(onlineLifts) /*|| key.timeout*/ {
			fmt.Printf("Laddr = %v\n", defs.Laddr)
			var (
				lowCost = defs.MaxInt
				lowAddr string
			)
			// Loop thru costs in each complete list
			for _, reply := range replyList {
				// ny kost: reply.cost
				// gammel kost: lowCost
				// ny ip: reply.lift
				// gammel ip: lowAddr

				// hvis ny bedre enn gammel: best = ny
				// hvis ny og gammel like bra og best: ny = lavest ip
				if reply.cost < lowCost {
					lowCost = reply.cost
					lowAddr = reply.lift
				} else if reply.cost == lowCost {
					if reply.lift < lowAddr {
						lowCost = reply.cost
						lowAddr = reply.lift
					}
				}
			}
			// Print winner:
			fmt.Printf("Lift %s won order f=%d b=%d\n", lowAddr[12:15], key.floor+1, key.button)
			// Assign order key to lift
			queue.AddRemoteOrder(key.floor, key.button, lowAddr)
			//queue.PrintQueues()
			if lowAddr == defs.Laddr.String() {
				fsm.EventExternalOrderGivenToMe()
			}
			// Empty list
			delete(que, key)
			// SUPERIMPORTANT: NOTIFY ABOUT EVENT NEW ORDER
		}
	}
}
