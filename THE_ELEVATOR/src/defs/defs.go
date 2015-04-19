package defs

import "net"

// Global system constants
const NumButtons = 3
const NumFloors = 4

const (
	ButtonCallUp int = iota
	ButtonCallDown
	ButtonCommand // Rename to ButtonInternal or something
)

const ( // Rename to DirDown etc.
	DirnDown int = iota - 1
	DirnStop
	DirnUp
)

const (
	Alive int = iota +1
	NewOrder
	CompleteOrder
	Cost
)

// Generic network message. No other messages are ever sent on the network.
type Message struct {
	Kind   int
	Floor  int
	Button int
	Cost   int
	Addr   string `json:"-"`
}

var MessageChan = make(chan Message)

var Laddr *net.UDPAddr //Local address
