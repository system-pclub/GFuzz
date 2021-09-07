package selefcm

import (
	"sync"
	"time"
)

func SelectWithCh() {
	ch1 := make(chan int)
	ch2 := make(chan struct{})

	select {
	case <-ch1:
		println("ch1!")
	case <-ch2:
		println("ch2")
	}
}

func SelectWithDefault() {
	ch1 := make(chan int)
	ch2 := make(chan struct{})

	select {
	case <-ch1:
		println("ch1!")
	case <-ch2:
		println("ch2")
	default:
		print("default")
	}
}

func SelectWithTimeout() {
	ch1 := make(chan int)
	ch2 := make(chan struct{})

	select {
	case <-ch1:
		println("ch1!")
	case <-ch2:
		println("ch2")
	case <-time.After(10):
		println("timeout")
	}
}

// Example from Paper "Who Goes First? Detecting Go Concurrency Bugs via Message Reordering" ====
type keeper struct {
	ch              chan int
	deleteTokenFunc func(*token)
}
type token struct {
	k  *keeper
	mu sync.RWMutex
}

func (t *token) assignTokenToUser() { //Goroutine1
	t.mu.Lock()
	t.k.ch <- 0
	t.mu.Unlock()
}
func (k *keeper) run() { //Goroutine2
	ticker := time.NewTicker()
	select {
	case <-k.ch:
		println("hey, I got selected")
	case <-ticker:
		k.deleteTokenFunc(t) // t is created before
	}
}
func newDeleteFunc(t *token) {
	t.mu.Lock()
}
