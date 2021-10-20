package wgrec

import (
	"gfuzz/pkg/oraclert"
	"sync"
)

func Hello() {
	wg := sync.WaitGroup{}
	oraclert.StoreOpInfo("Add", 0)

	wg.Add(1)
	oraclert.StoreOpInfo("Wait", 1)
	wg.Wait()
	oraclert.StoreOpInfo("Done", 2)
	wg.Done()
}
