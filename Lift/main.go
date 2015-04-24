package main

import (
	def "./src/config"
	"./src/fsm"
	"./src/hw"
	"./src/network"
	"./src/queue"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"os/exec"
	"time"
)

const debugPrint = false

var _ = log.Println
var _ = fmt.Println
var _ = errors.New

//Start a new terminal when restart.Run()
var restart = exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go")

var onlineLifts = make(map[string]network.UdpConnection)

var deadChan = make(chan network.UdpConnection)
var costChan = make(chan def.Message)
var costTimeoutChan = make(chan order)

type reply struct {
	cost int
	lift string
}
type order struct { //bad name?
	floor   int
	button  int
	timeout bool
	timer   *time.Timer
}

func main() {
	if err := hw.Init(); err != nil {
		log.Fatal(err)
	}

	e := fsm.EventChannels{
		NewOrder:     make(chan bool),
		FloorReached: make(chan int)}
	fsm.Init(e)

	network.Init()

	//handle ctrl+c
	safeKill() //bad name?

	liftAssigner(e.NewOrder)
	go poll(e)
	queue.Init(e.NewOrder)
	
	

	for { //nicer solution?
		time.Sleep(100 * time.Second)
	}
}

func poll(e fsm.EventChannels) {
	buttonChan := pollButtons()
	floorChan := pollFloors()

	for {
		select {
		case keypress := <-buttonChan:
			switch keypress.Button {
			case def.ButtonCommand:
				queue.AddLocalOrder(keypress.Floor, keypress.Button)
			case def.ButtonUp, def.ButtonDown:
				def.MessageChan <- def.Message{
					Kind:   def.NewOrder,
					Floor:  keypress.Floor,
					Button: keypress.Button}
			}
		case floor := <-floorChan:
			e.FloorReached <- floor
		case udpMessage := <-network.ReceiveChan:
			handleMessage(network.ParseMessage(udpMessage))
		case connection := <-deadChan:
			handleDeadLift(connection.Addr)
		case order:= <-queue.OrderTimeoutChan:
			fmt.Println("order in queue timed out, takes it myself")
			queue.RemoveRemoteOrdersAt(order.Floor)
			queue.AddRemoteOrder(order.Floor, order.Button , def.Laddr.String())			
		}
	}
}

func pollButtons() <-chan def.Keypress {
	c := make(chan def.Keypress)

	go func() {
		var buttonState [def.NumFloors][def.NumButtons]bool

		for {
			for f := 0; f < def.NumFloors; f++ {
				for b := 0; b < def.NumButtons; b++ {
					if (f == 0 && b == def.ButtonDown) ||
						(f == def.NumFloors-1 && b == def.ButtonUp) {
						continue
					}
					if hw.ReadButton(f, b) {
						if !buttonState[f][b] {
							c <- def.Keypress{Button: b, Floor: f}
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
func handleMessage(message def.Message) {

	network.PrintMessage(message)

	switch message.Kind {
	case def.Alive:
		if connection, exist := onlineLifts[message.Addr]; exist {
			connection.Timer.Reset(def.ResetTime)
			if debugPrint {
				fmt.Printf("Timer reset for IP %s\n", message.Addr)
			}
		} else {
			newConnection := network.UdpConnection{message.Addr, time.NewTimer(def.ResetTime)}
			onlineLifts[message.Addr] = newConnection
			if debugPrint {
				fmt.Printf("New connection with IP %s\n", message.Addr)
			}
			go connectionTimer(&newConnection)
		}
	case def.NewOrder:
		fmt.Printf("handleMessage(): NewOrder message: f=%d b=%d from lift %s\n",
			message.Floor+1, message.Button, message.Addr[12:15])

		cost := queue.CalculateCost(message.Floor, message.Button, fsm.Floor(), hw.Floor(), fsm.Direction())

		costMessage := def.Message{
			Kind:   def.Cost,
			Floor:  message.Floor,
			Button: message.Button,
			Cost:   cost}
		//fmt.Printf("handleMessage(): NewOrder sends cost message: f=%d b=%d (with cost %d) from me\n", costMessage.Floor+1, costMessage.Button, costMessage.Cost)
		def.MessageChan <- costMessage
	case def.CompleteOrder:
		fmt.Println("handleMessage(): CompleteOrder message")
		// remove from queues
		queue.RemoveRemoteOrdersAt(message.Floor)

		// prob more to do here
	case def.Cost:
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
	<-connection.Timer.C
	deadChan <- *connection
}

func costTimer(newOrder *order) {
	<-newOrder.timer.C
	costTimeoutChan <- *newOrder
}

//bad variable names am bad!
func liftAssigner(newOrderChan chan bool) {
	// collect cost values from all lifts
	// decide which lift gets the order when all lifts
	// in alive-list have answered or after a timeout
	// either send the decision on network or pray that all
	// lifts make the same choice every time
	go func() {
		assignmentQueue := make(map[order][]reply)
		var newOrder order
		for {
			select {
			case message := <-costChan:
				//newKey, newReply := split(message) //newKey is actually the worst name ever
				newOrder.makeNewOrder(message)
				newReply := getReply(message)

				for oldOrder := range assignmentQueue {
					if newOrder.isSameOrder(oldOrder) {
						newOrder = oldOrder
					}
				}
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
						newOrder.timer.Reset(10 * time.Second)
					}
				} else {
					// If order not in queue at all, init order list with it
					newOrder.timer = time.NewTimer(10 * time.Second)
					assignmentQueue[newOrder] = []reply{newReply}
					go costTimer(&newOrder)
				}
				evaluateLists(&assignmentQueue, newOrderChan)
			case newOrder := <-costTimeoutChan:
				fmt.Printf("\n ORDER TIMED OUT!\n\n")
				newOrder.setTimeout(true)
				evaluateLists(&assignmentQueue, newOrderChan)
			}
		}
	}()
}

func split(m def.Message) (order, reply) {
	return order{floor: m.Floor, button: m.Button}, reply{cost: m.Cost, lift: m.Addr}
}

func getReply(m def.Message) reply {
	return reply{cost: m.Cost, lift: m.Addr}
}

// evaluateLists goes through the map of orders with associated costs, checks
// if any orders have received answers from all live lifts, and finds the
// the best candidate for all such orders. The best candidate is added to the
// shared queue.
// This is very cryptic and ungood.
func evaluateLists(que *(map[order][]reply), newOrderChan chan bool) {
	// Loop thru all lists
	fmt.Printf("Lists: ")
	fmt.Println(*que)
	for key, replyList := range *que {
		// Check if the list is complete
		if len(replyList) == len(onlineLifts) || key.timeout {
			fmt.Printf("Laddr = %v\n", def.Laddr)
			var (
				lowCost = def.MaxInt
				lowAddr string
			)
			// Loop thru costs in each complete list
			for _, reply := range replyList {
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
			
			// Empty list and stop timer
			key.timer.Stop()
			delete(*que, key)
		}
	}
}

//safeKill gets the motor to stop when the program is killed with ctrl+c
func safeKill() {
	var c = make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		hw.SetMotorDirection(def.DirStop)
		//restart.Run() //vil vi restarte med ctrl+c?
		log.Fatal("[FATAL]\tUser terminated program")
	}()
}

// --------------- METHODS FOR ORDER TYPE: ---------------

func (o *order) makeNewOrder(msg def.Message) {
	o.floor = msg.Floor
	o.button = msg.Button
}

func (o *order) isSameOrder(other order) bool {
	if other.floor == o.floor && other.button == o.button {
		return true
	} else {
		return false
	}
}

func (o *order) setTimeout(b bool) {
	o.timeout = b
}
