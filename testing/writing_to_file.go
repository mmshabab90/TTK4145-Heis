// Writing files in Go follows similar patterns to the
// ones we saw earlier for reading.

package main

import (
    "fmt"
    "io/ioutil"
    "encoding/gob"
    "bytes"
    "log"
)

var _ = fmt.Println
var _ = ioutil.ReadFile

type thing struct {
    a int
    b int
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    // data1 := "hello\ngo\n"

    // err := ioutil.WriteFile("temp", []byte(data1), 0644)
    // check(err)

    // data2, err := ioutil.ReadFile("temp")
    // check(err)

    // fmt.Println(string(data2) == data1)



    // data3 := thing{1, 2}
    // err := ioutil.WriteFile("structfile", []byte(data3), 0644)
    // check(err)

    // data4, err := ioutil.ReadFile("structfile")
    // check(err)

    // fmt.Println(thing(data4))



    var network bytes.Buffer // Stand-in for the network.

    // Create an encoder and send a value.
    enc := gob.NewEncoder(&network)
    err := enc.Encode(thing{3, 4})
    if err != nil {log.Fatal("encode:", err)}

    // Create a decoder and receive a value.
    dec := gob.NewDecoder(&network)
    var v thing
    err = dec.Decode(&v)
    if err != nil {log.Fatal("decode:", err)}
    fmt.Println(v)
}