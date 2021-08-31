package selefcm

import (
	gooracle "gooracle"
	"time"
)

func Hello() {
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
