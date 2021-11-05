package pass

import (
	"sync"
)

var opID uint16
var opMu sync.Mutex
var records []string
var recMu sync.Mutex

func getNewOpID() uint16 {
	opMu.Lock()
	defer opMu.Unlock()
	opID++
	return opID
}

func addRecord(record string) {
	recMu.Lock()
	defer recMu.Unlock()
	records = append(records, record)
}
