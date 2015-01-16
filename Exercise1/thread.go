
package main

import(
	. "fmt"
	"runtime"
	"time"
)

var(
	i int = 0
) 

func thread1foo(){
	for j := 0; j < 10000; j++{
		i++
	}
}

func thread2foo(){
	for j := 0; j < 10000; j++{
		i--
	}
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	go thread1foo()

	go thread2foo()

	time.Sleep(100*time.Millisecond)
	
	Println(i)
}
