package cvrec

import (
	gooracle "gooracle"
	"sync"
)

func Hello() {
	m := sync.Mutex{}

	c := sync.NewCond(&m)
	gooracle.StoreOpInfo("Broadcast", 0)

	c.Broadcast()
	gooracle.StoreOpInfo("Signal", 1)

	c.Signal()
	gooracle.StoreOpInfo("Wait", 2)

	c.Wait()
}
