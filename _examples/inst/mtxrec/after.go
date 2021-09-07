package mtxrec

import (
	gooracle "gooracle"
	"sync"
)

func Hello() {
	m := sync.Mutex{}
	gooracle.StoreOpInfo("Lock", 0)

	m.Lock()
	gooracle.StoreOpInfo("Unlock", 1)

	m.Unlock()

	rwm := sync.RWMutex{}
	gooracle.StoreOpInfo("Lock", 2)

	rwm.Lock()
	gooracle.StoreOpInfo("RLock", 3)
	rwm.RLock()
	gooracle.StoreOpInfo("Unlock", 4)
	rwm.Unlock()
	gooracle.StoreOpInfo("RUnlock", 5)
	rwm.RUnlock()
	gooracle.StoreOpInfo("Lock", 6)
	rwm.RLocker().Lock()
	gooracle.StoreOpInfo("Unlock", 7)
	rwm.RLocker().Unlock()

}
