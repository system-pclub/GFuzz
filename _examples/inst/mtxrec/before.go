package mtxrec

import "sync"

func Hello() {
	m := sync.Mutex{}

	m.Lock()

	m.Unlock()

	rwm := sync.RWMutex{}

	rwm.Lock()
	rwm.RLock()
	rwm.Unlock()
	rwm.RUnlock()
	rwm.RLocker().Lock()
	rwm.RLocker().Unlock()

}
