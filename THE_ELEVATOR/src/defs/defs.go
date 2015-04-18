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

type Message struct {
	Kind   int
	Floor  int
	Button int
	Cost   int
	Addr   string `json:"-"`
}

chan MessageChan Message
