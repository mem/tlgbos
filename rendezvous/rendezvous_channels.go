package main

import (
	"fmt"
	"sync"
)

type Signal struct{}

func worker(name string, ready, wait chan Signal) {
	fmt.Printf("%s%d\n", name, 1)
	ready <- Signal{}
	<-wait
	fmt.Printf("%s%d\n", name, 2)
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

	aReady := make(chan Signal, 1)
	bReady := make(chan Signal, 1)

	goWg(wg, func() { worker("a", aReady, bReady) })
	goWg(wg, func() { worker("b", bReady, aReady) })
}
