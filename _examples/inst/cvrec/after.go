package cvrec

import (
	oraclert "gfuzz/pkg/oraclert"
	"sync"
)

func Hello() {
	m := sync.Mutex{}

	c := sync.NewCond(&m)
	oraclert.StoreOpInfo("Broadcast", 1)

	c.Broadcast()
	oraclert.StoreOpInfo("Signal", 2)

	c.Signal()
	oraclert.StoreOpInfo("Wait", 3)

	c.Wait()
}
