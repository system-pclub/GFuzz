package selefcm

import (
	oraclert "gfuzz/pkg/oraclert"
	"sync"
	"time"
)

func SelectWithCh() {
	ch1 := make(chan int)
	ch2 := make(chan struct{})
	switch oraclert.GetSelEfcmSwitchCaseIdx("/home/xsh/code/GFuzz/_examples/inst/selefcm/before.go", "12", 2) {
	case 0:
		select {
		case <-ch1:
			println("ch1!")
		case <-oraclert.SelEfcmTimeout():
			oraclert.StoreLastMySwitchChoice(-1)
			select {
			case <-ch1:
				println("ch1!")
			case <-ch2:
				println("ch2")
			}
		}
	case 1:
		select {
		case <-ch2:
			println("ch2")
		case <-oraclert.SelEfcmTimeout():
			oraclert.StoreLastMySwitchChoice(-1)
			select {
			case <-ch1:
				println("ch1!")
			case <-ch2:
				println("ch2")
			}
		}
	default:
		oraclert.StoreLastMySwitchChoice(-1)
		select {
		case <-ch1:
			println("ch1!")
		case <-ch2:
			println("ch2")
		}
	}
}

func SelectWithDefault() {
	ch1 := make(chan int)
	ch2 := make(chan struct{})
	switch oraclert.GetSelEfcmSwitchCaseIdx("/home/xsh/code/GFuzz/_examples/inst/selefcm/before.go", "24", 3) {
	case 0:
		select {
		case <-ch1:
			println("ch1!")
		case <-oraclert.SelEfcmTimeout():
			oraclert.StoreLastMySwitchChoice(-1)
			select {
			case <-ch1:
				println("ch1!")
			case <-ch2:
				println("ch2")
			default:
				print("default")
			}
		}
	case 1:
		select {
		case <-ch2:
			println("ch2")
		case <-oraclert.SelEfcmTimeout():
			oraclert.StoreLastMySwitchChoice(-1)
			select {
			case <-ch1:
				println("ch1!")
			case <-ch2:
				println("ch2")
			default:
				print("default")
			}
		}
	default:
		oraclert.StoreLastMySwitchChoice(-1)
		select {
		case <-ch1:
			println("ch1!")
		case <-ch2:
			println("ch2")
		default:
			print("default")
		}
	}
}

func SelectWithTimeout() {
	ch1 := make(chan int)
	ch2 := make(chan struct{})
	switch oraclert.GetSelEfcmSwitchCaseIdx("/home/xsh/code/GFuzz/_examples/inst/selefcm/before.go", "38", 3) {
	case 0:
		select {
		case <-ch1:
			println("ch1!")
		case <-oraclert.SelEfcmTimeout():
			oraclert.StoreLastMySwitchChoice(-1)
			select {
			case <-ch1:
				println("ch1!")
			case <-ch2:
				println("ch2")
			case <-time.After(10):
				println("timeout")
			}
		}
	case

		// Example from Paper "Who Goes First? Detecting Go Concurrency Bugs via Message Reordering" ====
		1:
		select {
		case <-ch2:
			println("ch2")
		case <-oraclert.SelEfcmTimeout():
			oraclert.StoreLastMySwitchChoice(-1)
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
		case <-oraclert.SelEfcmTimeout():
			oraclert.StoreLastMySwitchChoice(-1)
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
		oraclert.StoreLastMySwitchChoice(-1)
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
	switch oraclert.GetSelEfcmSwitchCaseIdx("/home/xsh/code/GFuzz/_examples/inst/selefcm/before.go", "65",

		// t is created before
		2) {
	case 0:
		select {
		case <-k.ch:
			println("hey, I got selected")
		case <-oraclert.SelEfcmTimeout():
			oraclert.StoreLastMySwitchChoice(-1)
			select {
			case <-k.ch:
				println("hey, I got selected")
			case <-ticker:
				k.deleteTokenFunc(t)
			}
		}
	case 1:
		select {
		case <-ticker:
			k.deleteTokenFunc(t)
		case <-oraclert.SelEfcmTimeout():
			oraclert.StoreLastMySwitchChoice(-1)
			select {
			case <-k.ch:
				println("hey, I got selected")
			case <-ticker:
				k.deleteTokenFunc(t)
			}
		}
	default:
		oraclert.StoreLastMySwitchChoice(-1)
		select {
		case <-k.ch:
			println("hey, I got selected")
		case <-ticker:
			k.deleteTokenFunc(t)
		}
	}
}
func newDeleteFunc(t *token) {
	t.mu.Lock()
}
