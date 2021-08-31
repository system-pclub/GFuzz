package wgrec

import (
	gooracle "gooracle"
	"sync"
)

func Hello() {
	wg := sync.WaitGroup{}
	gooracle.StoreOpInfo("Add", 0)

	wg.Add(1)
	gooracle.StoreOpInfo("Wait", 1)
	wg.Wait()
	gooracle.StoreOpInfo("Done", 2)
	wg.Done()
}
