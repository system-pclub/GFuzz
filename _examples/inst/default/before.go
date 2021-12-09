package defaultp

import aaa "sync"

func Hello() {
	m := aaa.Mutex{}

	c := aaa.NewCond(&m)

	c.Broadcast()

	c.Signal()

	//asdfasdf
	c.Wait()
	//asdfadfas

	w := aaa.WaitGroup{}

	w.Wait()
}
