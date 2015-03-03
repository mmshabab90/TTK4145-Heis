package main

import(
	. "fmt"
	"runtime"
	"time"
)

var(
	i int = 0
) 

<<<<<<< HEAD
func thread1foo(){
	for j := 0; j < 1000000; j++{
=======
func increase() {
	for j := 0; j < 1000000; j++ {
>>>>>>> bc43ba34e43f9a46131a21e82a1536a7d15b2b61
		i++
	}
}

<<<<<<< HEAD
func thread2foo(){
	for j := 0; j < 1000000; j++{
=======
func decrease() {
	for j := 0; j < 1000000; j++ {
>>>>>>> bc43ba34e43f9a46131a21e82a1536a7d15b2b61
		i--
	}
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	go increase()
	go decrease()

	time.Sleep(100*time.Millisecond)
	
	Println(i)
}
