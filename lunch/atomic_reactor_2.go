package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func eatLunch(name string) {
	fmt.Println(name, "started having lunch")
	time.Sleep(time.Duration(30+rand.Intn(30)) * 10 * time.Millisecond)
	fmt.Println(name, "finished having lunch")
}

func alice(l sync.Locker) {
	l.Lock()
	eatLunch("Alice")
	l.Unlock()
}

func bob(l sync.Locker) {
	l.Lock()
	eatLunch("Bob")
	l.Unlock()
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

	l := &sync.Mutex{}

	goWg(wg, func() { alice(l) })

	goWg(wg, func() { bob(l) })
}
