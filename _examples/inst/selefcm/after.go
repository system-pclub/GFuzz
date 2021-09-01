package selefcm

import (
	gooracle "gooracle"
	"time"
)

func Hello() {
	ch1 := make(chan int)
	ch2 := make(chan struct{})
	switch gooracle.ReadSelect("/gfuzz/_examples/inst/selefcm/before.go", 9, 3) {
	case 0:
		select {
		case <-ch1:
			println("ch1!")
		case <-gooracle.SelectTimeout():
			gooracle.StoreLastMySwitchChoice(-1)
			select {
			case <-ch1:
				println("ch1!")
			case <-ch2:
				println("ch2")
			case <-time.After(10):
				println("timeout")
			}
		}
	case 1:
		select {
		case <-ch2:
			println("ch2")
		case <-gooracle.SelectTimeout():
			gooracle.StoreLastMySwitchChoice(-1)
			select {
			case <-ch1:
				println("ch1!")
			case <-ch2:
				println("ch2")
			case <-time.After(10):
				println("timeout")
			}
		}
	case 2:
		select {
		case <-time.After(10):
			println("timeout")
		case <-gooracle.SelectTimeout():
			gooracle.StoreLastMySwitchChoice(-1)
			select {
			case <-ch1:
				println("ch1!")
			case <-ch2:
				println("ch2")
			case <-time.After(10):
				println("timeout")
			}
		}
	default:
		gooracle.StoreLastMySwitchChoice(-1)
		select {
		case <-ch1:
			println("ch1!")
		case <-ch2:
			println("ch2")
		case <-time.After(10):
			println("timeout")
		}
	}
}
