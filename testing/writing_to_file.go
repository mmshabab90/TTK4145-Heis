// Writing files in Go follows similar patterns to the
// ones we saw earlier for reading.

package main

import (
    "fmt"
    "io/ioutil"
    "encoding/gob"
    "os"
    "math/rand"
    "time"
)

var _ = fmt.Println
var _ = ioutil.ReadFile

type Thing struct {
    A int
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    rand.Seed(time.Now().UTC().UnixNano())

    data := new(Thing)

    for i := 0; i < 3; i++ {
        fmt.Println("-----")

        if err := data.Load("temp"); err == os.PathError {
            fmt.Println(err)
        }
        if err != nil

        fmt.Printf("loaded value is %v\n", data.A)

        value := rand.Intn(100)
        fmt.Printf("new value is    %v\n", value)
        data.A = value

        if err := data.Save("temp"); err != nil {
            fmt.Println("save", err)
        }
        time.Sleep(100*time.Millisecond)
    }
}

func (t *Thing) Load(filename string) error {

    fi, err := os.Open(filename)
    if err !=nil {
        return err
    }
    defer fi.Close()

    decoder := gob.NewDecoder(fi)
    err = decoder.Decode(&t)
    if err !=nil {
        return err
    }

    return nil
}

func (t *Thing) Save(filename string) error {

    fi, err := os.Create(filename)
    if err !=nil {
        return err
    }
    defer fi.Close()

    encoder := gob.NewEncoder(fi)
    err = encoder.Encode(t)
    if err !=nil {
        return err
    }

    return nil
}
