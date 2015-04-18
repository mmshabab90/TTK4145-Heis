package main

type cost struct {
	cost int
	lift string
}
type order struct {
	floor int
	button int
}

func main() {
	go func() {
		// This is a map that indexes on a struct and returns a slice:
		assignmentQueue := make(map[order][]cost)
	}()
}
