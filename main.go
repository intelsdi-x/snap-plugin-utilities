package main

import (
	. "github.com/intelsdi-x/pulse-plugin-utilities/pipeline"
	"fmt"
)

type Add struct {}

func (a Add) Process(in chan interface{}) chan interface {} {
	out := make(chan interface{})
	go func() {
		for i := range in {
			val := i.(int)
			out <- val + 10
		}
		close(out)
	}()
	return out
}


func main() {
	p := NewPipeline(Add{})

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Printf("Sending %d\n", i)
			p.Enqueue(i)
		}
		fmt.Printf("Closing pipeline\n")
		p.Close()
	}()

	p.Dequeue(func(i interface{}) {
		fmt.Printf("Recieving %v\n", i)
	})
}
