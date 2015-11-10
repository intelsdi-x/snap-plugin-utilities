package main

import (
	"fmt"
	"time"
	. "github.com/intelsdi-x/pulse-plugin-utilities/source"

)

func main() {
	time.Sleep(time.Millisecond)
	ech := make(chan error)
	out := make(chan interface{})
	s := Source{"ls", []string{"-ltr"}}
	go s.Generate(out, ech)

	LOOP:
	for {
		time.Sleep(100 * time.Millisecond)
		select {
		case data, ok := <- out:
			if !ok {
				fmt.Printf("Waiting for out ...\n")
				time.Sleep(1 * time.Second)
				fmt.Printf("Out empty!\n")
				break LOOP;
			}
			fmt.Printf(">>> Recieving {%v}\n", data)
		case e := <- ech:
			fmt.Printf("ERRROR {%v}\n", e)
			break LOOP;
		case <-time.After(time.Second * 2):
			fmt.Printf("No activity\n")
			break LOOP
		}
	}
}
