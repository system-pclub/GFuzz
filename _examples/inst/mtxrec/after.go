package mtxrec

import (
	oraclert "gfuzz/pkg/oraclert"
	"sync"
)

func Hello() {
	m := sync.Mutex{}
	oraclert.StoreOpInfo("Lock", 1)

	m.Lock()
	oraclert.StoreOpInfo("Unlock", 2)

	m.Unlock()

	rwm := sync.RWMutex{}
	oraclert.StoreOpInfo("Lock", 3)

	rwm.Lock()
	oraclert.StoreOpInfo("RLock", 4)
	rwm.RLock()
	oraclert.StoreOpInfo("Unlock", 5)
	rwm.Unlock()
	oraclert.StoreOpInfo("RUnlock", 6)
	rwm.RUnlock()
	oraclert.StoreOpInfo("Lock", 7)
	rwm.RLocker().Lock()
	oraclert.StoreOpInfo("Unlock", 8)
	rwm.RLocker().Unlock()

}
