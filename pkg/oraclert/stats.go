package oraclert

import (
	"sync"
)

var (
	getSelEfcmCount uint32
	origSelCount    uint32
	notSelEfcmCount uint32
	opRecords       []uint16
	opMutex         sync.Mutex
)

func recordOp(opID uint16) {
	opMutex.Lock()
	opRecords = append(opRecords, opID)
	opMutex.Unlock()
}
