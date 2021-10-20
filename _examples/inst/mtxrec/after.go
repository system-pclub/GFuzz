package mtxrec

import (
	"gfuzz/pkg/oraclert"
	"sync"
)

func Hello() {
	m := sync.Mutex{}
	oraclert.StoreOpInfo("Lock", 0)

	m.Lock()
	oraclert.StoreOpInfo("Unlock", 1)

	m.Unlock()

	rwm := sync.RWMutex{}
	oraclert.StoreOpInfo("Lock", 2)

	rwm.Lock()
	oraclert.StoreOpInfo("RLock", 3)
	rwm.RLock()
	oraclert.StoreOpInfo("Unlock", 4)
	rwm.Unlock()
	oraclert.StoreOpInfo("RUnlock", 5)
	rwm.RUnlock()
	oraclert.StoreOpInfo("Lock", 6)
	rwm.RLocker().Lock()
	oraclert.StoreOpInfo("Unlock", 7)
	rwm.RLocker().Unlock()

}
