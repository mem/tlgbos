package main

import (
	"fmt"
	"sync"
)

type Signal struct{}

type Barrier struct {
	ready chan Signal
	wait  chan Signal
	done  chan struct{}
}

func (b Barrier) Join() {
	b.ready <- Signal{}
	<-b.wait
}

func NewBarrier(N int) Barrier {
	ready := make(chan Signal)
	wait := make(chan Signal)
	done := make(chan struct{})

	helper := func() {
		for {
			for i := N; i > 0; i-- {
				select {
				case <-ready:
				case <-done:
					close(ready)
					close(wait)
					return
				}
			}
			for i := N; i > 0; i-- {
				select {
				case wait <- Signal{}:
				case <-done:
					close(ready)
					close(wait)
					return
				}
			}
		}
	}

	go helper()

	return Barrier{
		ready: ready,
		wait:  wait,
		done:  done,
	}
}

func (b Barrier) Done() {
	close(b.done)
}

func worker(name string, N int, b1, b2 Barrier) {
	for i := N; i > 0; i-- {
		fmt.Printf("Worker %s entering: %d\n", name, N-i)
		b1.Join()
		fmt.Printf("Worker %s leaving: %d\n", name, N-i)
		b2.Join()
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

	N, M := 10, 5

	b1 := NewBarrier(N)
	b2 := NewBarrier(N)

	for i := 0; i < N; i++ {
		name := string('a' + i)
		goWg(wg, func() { worker(name, M, b1, b2) })
	}

	wg.Wait()

	b1.Done()
	b2.Done()
}
