package cvrec

import "sync"

func Hello() {
	m := sync.Mutex{}

	c := sync.NewCond(&m)

	c.Broadcast()

	c.Signal()

	c.Wait()
}
