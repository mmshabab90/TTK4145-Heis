package config

import (
	"net"
)

// Global system constants
const NumFloors = 4
const NumButtons = 3

const (
	ButtonUp int = iota
	ButtonDown
	ButtonIn
)

const (
	DirDown int = iota - 1
	DirStop
	DirUp
)

type Keypress struct {
	Floor  int
	Button int
}

const MaxInt = int(^uint(0) >> 1)

// Move into network module and pass to other modules
var Laddr *net.UDPAddr //Local address
