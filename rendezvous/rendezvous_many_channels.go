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

func helper(N int, ready, wait chan Signal) {
	for i := N; i > 0; i-- {
		<-ready
	}
	for i := N; i > 0; i-- {
		wait <- Signal{}
	}
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

	N := 10
	ready := make(chan Signal)
	wait := make(chan Signal)

	for i := 0; i < N; i++ {
		name := string('a' + i)
		goWg(wg, func() { worker(name, ready, wait) })
	}

	goWg(wg, func() { helper(N, ready, wait) })
}
