package main

import (
	"fmt"
)

const nF = 3
const nB = 3

type queue struct {
	local [nF][nB] int
}

func (q *queue) alter() {
	for f := 0; f < nF; f++ {
		for b := 0; b < nB; b++ {
			q.local[f][b] = f*nB + b
		}
	}
}

func main() {
	myQ := new(queue)
	fmt.Println(myQ)
	
	myQ.alter()
	fmt.Println(myQ)
}