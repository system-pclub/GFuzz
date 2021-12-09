package defaultp

import (
	oraclert "gfuzz/pkg/oraclert"
	aaa "sync"
)

func Hello() {
	m := aaa.Mutex{}

	c := aaa.NewCond(&m)
	oraclert.StoreOpInfo("Broadcast",

		//asdfasdf
		1)

	c.Broadcast()
	oraclert.StoreOpInfo("Signal", 2)

	c.Signal()
	oraclert.StoreOpInfo("Wait",

		//asdfadfas
		3)
	oraclert.StoreOpInfo("Wait", 5)

	c.Wait()

	w := aaa.WaitGroup{}
	oraclert.StoreOpInfo("Wait", 4)
	oraclert.StoreOpInfo("Wait", 6)

	w.Wait()
}
