package hello

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestChannelBug(t *testing.T) {

	ch := make(chan int)
	go func() {
		ch <- 1
	}()

	select {
	case <-ch:
		fmt.Println("Normal")
	case <-time.After(300 * time.Millisecond):
		fmt.Println("Should be buggy")
	}

}

func TestWgBug(t *testing.T) {

	var wg sync.WaitGroup
	wg.Add(1)

	wg.Wait()
}
