package cvrec

import (
	"gfuzz/pkg/oraclert"
	"sync"
)

func Hello() {
	m := sync.Mutex{}

	c := sync.NewCond(&m)
	oraclert.StoreOpInfo("Broadcast", 0)

	c.Broadcast()
	oraclert.StoreOpInfo("Signal", 1)

	c.Signal()
	oraclert.StoreOpInfo("Wait", 2)

	c.Wait()
}
