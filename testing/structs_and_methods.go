package main

import (
	"fmt"
)

const nF = 3
const nB = 3
const maxLifts = 10

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

type reply struct {
	cost int
	lift string
}

func main() {
	// myQ := new(queue)
	// fmt.Println(myQ)
	
	// myQ.alter()
	// fmt.Println(myQ)

	// fmt.Println(queue{})

	ass := make(map[string][]int)
	ass["key"] = []int{24}
	fmt.Println(ass)

	ass["key"] = append(ass["key"], 98)
	fmt.Println(ass)
}