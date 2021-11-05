package wgrec

import (
	oraclert "gfuzz/pkg/oraclert"
	"sync"
)

func Hello() {
	wg := sync.WaitGroup{}
	oraclert.StoreOpInfo("Add", 1)

	wg.Add(1)
	oraclert.StoreOpInfo("Wait", 2)
	wg.Wait()
	oraclert.StoreOpInfo("Done", 3)
	wg.Done()
}
