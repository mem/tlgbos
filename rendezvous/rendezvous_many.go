package main

import (
	"fmt"
	"sync"
)

func worker(name string, rendezvous *sync.WaitGroup) {
	fmt.Printf("%s%d\n", name, 1)
	rendezvous.Done()
	rendezvous.Wait()
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

	rendezvous := &sync.WaitGroup{}

	N := 10
	rendezvous.Add(N)

	for i := 0; i < N; i++ {
		name := string('a' + i)
		goWg(wg, func() { f(name, rendezvous) })
	}
}
