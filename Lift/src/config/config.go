package config

import (
	"os/exec"
)

// Global system constants
const NumButtons = 3
const NumFloors = 4

const (
	BtnUp int = iota
	BtnDown
	BtnInside
)

const (
	DirDown int = iota - 1
	DirStop
	DirUp
)

// Local IP address
var Laddr string

//message kind constants
const (
	Alive int = iota + 1
	NewOrder
	CompleteOrder
	Cost
)

type Keypress struct {
	Button int
	Floor  int
}

// Generic network message. No other messages are ever sent on the network.
type Message struct {
	Category int
	Floor    int
	Button   int
	Cost     int
	Addr     string `json:"-"`
}

var OutgoingMsg = make(chan Message, 10) // vurder buff //todo: rename outgoing/outbox (and move?)
var SyncLightsChan = make(chan bool)     // todo move!
var CloseConnectionChan = make(chan bool)

// Start a new terminal when restart.Run()
var Restart = exec.Command("gnome-terminal", "-x", "sh", "-c", "lift")

// Colours for printing to console
const Clr0 = "\x1b[30;1m" // Dark grey
const ClrR = "\x1b[31;1m" // Red
const ClrG = "\x1b[32;1m" // Green
const ClrY = "\x1b[33;1m" // Yellow
const ClrB = "\x1b[34;1m" // Blue
const ClrM = "\x1b[35;1m" // Magenta
const ClrC = "\x1b[36;1m" // Cyan
const ClrW = "\x1b[37;1m" // White
const ClrN = "\x1b[0m"    // Grey (neutral)
