package main

import (
	"fmt"
	"sync"
	"time"
)

type DanceStep struct {
	data string
}

func leader(name string, queue chan<- chan DanceStep) {
	ch := make(chan DanceStep)
	queue <- ch
	for i := 0; i < 5; i++ {
		time.Sleep(50 * time.Millisecond)
		ch <- DanceStep{data: fmt.Sprintf("%s in step %d", name, i)}
	}
	close(ch)
}

func follower(name string, queue <-chan chan DanceStep) {
	ch := <-queue
	for s := range ch {
		fmt.Printf("%s following %s\n", name, s.data)
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

	N := 4

	queue := make(chan chan DanceStep)

	for i := 0; i < N; i++ {
		name := string('a' + i)
		goWg(wg, func() { leader("l"+name, queue) })
		goWg(wg, func() { follower("f"+name, queue) })
	}

	wg.Wait()
}
