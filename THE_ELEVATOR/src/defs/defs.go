package defs

const NumButtons = 3
const NumFloors = 4
const (
	ButtonCallUp int = iota
	ButtonCallDown
	ButtonCommand
)
const (
	DirnDown int = iota - 1
	DirnStop
	DirnUp
)