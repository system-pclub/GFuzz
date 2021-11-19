package gexecfuzz

import (
	"gfuzz/pkg/gexec"
	"gfuzz/pkg/oraclert/output"
	"sync"
)

// GExecFuzz all information about fuzzing progress for given executable
type GExecFuzz struct {
	BestScore int
	Exec      gexec.Executable

	// Oracler runtime related
	OrtConfigHashes []string
	// All selects we have seen so far
	OrtSelects map[string]output.SelectRecord
	m          sync.RWMutex
}

func NewGExecFuzz(exec gexec.Executable) *GExecFuzz {
	return &GExecFuzz{
		Exec:      exec,
		BestScore: 0,
	}
}

func (e *GExecFuzz) String() string {
	return e.Exec.String()
}

func (e *GExecFuzz) UpdateSelectRecordsIfNew(records []output.SelectRecord) int {
	e.m.Lock()
	defer e.m.Unlock()
	newSelects := 0
	for _, rec := range records {
		if _, exists := e.OrtSelects[rec.ID]; !exists {
			e.OrtSelects[rec.ID] = rec
			newSelects += 1
		}
	}
	return newSelects
}

func (e *GExecFuzz) GetAllSelectRecords() []output.SelectRecord {
	e.m.Lock()
	defer e.m.Unlock()
	var records []output.SelectRecord
	for _, rec := range e.OrtSelects {
		records = append(records, rec)
	}
	return records
}
