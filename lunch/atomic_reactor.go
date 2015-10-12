package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Signal struct{}

func eatLunch(name string) {
	fmt.Println(name, "started having lunch")
	time.Sleep(time.Duration(30+rand.Intn(30)) * 10 * time.Millisecond)
	fmt.Println(name, "finished having lunch")
}

func alice(ch chan Signal) {
	<-ch
	eatLunch("Alice")
	ch <- Signal{}
}

func bob(ch chan Signal) {
	<-ch
	eatLunch("Bob")
	ch <- Signal{}
}

func goWg(wg *sync.WaitGroup, f func()) {
	wg.Add(1)
	go func() {
		f()
		wg.Done()
	}()
}

func main() {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	ch := make(chan Signal, 1)
	ch <- Signal{}

	goWg(wg, func() { alice(ch) })

	goWg(wg, func() { bob(ch) })
}
