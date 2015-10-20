package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Semaphore interface {
	decrement()
	increment()
}

type SyncSemaphore struct {
	n int32
	c *sync.Cond
}

func NewSyncSemaphore(N int32) Semaphore {
	return &SyncSemaphore{
		n: N,
		c: sync.NewCond(new(sync.Mutex)),
	}
}

func (s *SyncSemaphore) decrement() {
	s.c.L.Lock()
	for s.n <= 0 {
		s.c.Wait()
	}
	s.n--
	s.c.L.Unlock()
}

func (s *SyncSemaphore) increment() {
	s.c.L.Lock()
	s.n++
	s.c.L.Unlock()
	s.c.Signal()
}

type ChSemaphore chan struct{}

func NewChSemaphore(N int32) Semaphore {
	s := ChSemaphore(make(chan struct{}, N))
	for i := N; i > 0; i-- {
		s.increment()
	}
	return s
}

func (s ChSemaphore) decrement() {
	<-s
}

func (s ChSemaphore) increment() {
	s <- struct{}{}
}

func worker(s Semaphore, n *int32, stop int32, name string) {
	for done := false; !done; {
		s.decrement()
		if m := atomic.AddInt32(n, 1); m <= stop {
			fmt.Println(name, m)
		} else {
			done = true
		}
		s.increment()
	}
}

func main() {
	wg := &sync.WaitGroup{}

	N := 10
	s := NewChSemaphore(int32(N))
	n := new(int32)

	for i := 0; i < 2*N; i++ {
		name := string('a' + i)
		goWg(wg, func() {
			worker(s, n, 1000000, name)
		})
	}

	wg.Wait()
}

func goWg(wg *sync.WaitGroup, f func()) {
	wg.Add(1)
	go func() {
		f()
		wg.Done()
	}()
}
