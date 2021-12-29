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

	// mapping from oracle config hash to timeout count
	ortCfg2TimeoutCnt map[string]uint32
	ortCfgHash        map[string]struct{}
	// All selects we have seen so far
	OrtSelects  map[string]output.SelectRecord
	OrtChannels map[string]output.ChanRecord
	OrtTuples   map[uint32]uint32
	CaseRecords map[string][]int
	m           sync.RWMutex
}

func NewGExecFuzz(exec gexec.Executable) *GExecFuzz {
	return &GExecFuzz{
		Exec:              exec,
		BestScore:         0,
		OrtSelects:        make(map[string]output.SelectRecord),
		OrtChannels:       make(map[string]output.ChanRecord),
		OrtTuples:         make(map[uint32]uint32),
		CaseRecords:       make(map[string][]int),
		ortCfg2TimeoutCnt: make(map[string]uint32),
		ortCfgHash:        make(map[string]struct{}),
	}
}

func (e *GExecFuzz) String() string {
	return e.Exec.String()
}

func (e *GExecFuzz) Clean() {
	e.m.Lock()
	defer e.m.Unlock()
	e.OrtSelects = make(map[string]output.SelectRecord)
	e.OrtChannels = make(map[string]output.ChanRecord)
	e.OrtTuples = make(map[uint32]uint32)
	e.ortCfgHash = make(map[string]struct{})
	e.CaseRecords = make(map[string][]int)
}

func (e *GExecFuzz) RecordCase(rec output.SelectRecord) {
	e.m.Lock()
	defer e.m.Unlock()

	if cases, ok := e.CaseRecords[rec.ID]; ok {
		// since length cases usually small, we don't need map
		for _, c := range cases {
			// check is this case has been recorded
			if c == int(rec.Chosen) {
				continue
			}
		}

		cases = append(cases, int(rec.Chosen))

		e.CaseRecords[rec.ID] = cases
	} else {
		cases := []int{int(rec.Chosen)}
		e.CaseRecords[rec.ID] = cases
	}

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

func (e *GExecFuzz) UpdateChannelRecordsIfNew(records map[string]output.ChanRecord) int {
	e.m.Lock()
	defer e.m.Unlock()
	newChannels := 0

	for k, v := range records {
		if _, exists := e.OrtChannels[k]; !exists {
			// If the channel hasn't been seen before
			e.OrtChannels[k] = v
			newChannels += 1
		} else {
			// If the channel has a new state.
			curSavedRecord := e.OrtChannels[k]
			isNewChannel := false
			// If new notclosed.
			if curSavedRecord.NotClosed == false && v.NotClosed == true {
				curSavedRecord.NotClosed = true
				isNewChannel = true
			}
			// If new closed
			if curSavedRecord.Closed == false && v.Closed == true {
				curSavedRecord.Closed = true
				isNewChannel = true
			}
			// If new PeakBuf == CapBuf
			if curSavedRecord.PeakBuf != v.PeakBuf && v.PeakBuf == v.CapBuf {
				curSavedRecord.PeakBuf = v.PeakBuf
				isNewChannel = true
			}
			if isNewChannel {
				newChannels += 1
			}
		}
	}

	return newChannels
}

func (e *GExecFuzz) UpdateTupleRecordsIfNew(records map[uint32]uint32) int {
	e.m.Lock()
	defer e.m.Unlock()
	newTuples := 0

	for k, v := range records {
		if _, exists := e.OrtTuples[k]; !exists {
			e.OrtTuples[k] = v
			newTuples += 1
		}
	}

	return newTuples
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

func (e *GExecFuzz) RecordTimeoutOrtCfgHash(h string) {
	e.m.Lock()
	defer e.m.Unlock()
	e.ortCfg2TimeoutCnt[h] += 1
}

func (e *GExecFuzz) HasTimeoutOrtCfgHash(h string) bool {
	e.m.RLock()
	defer e.m.RUnlock()
	_, exist := e.ortCfg2TimeoutCnt[h]
	return exist
}

func (e *GExecFuzz) RecordOrtCfgHash(h string) {
	e.m.Lock()
	defer e.m.Unlock()
	e.ortCfgHash[h] = struct{}{}
}

func (e *GExecFuzz) HasOrtCfgHash(h string) bool {
	e.m.RLock()
	defer e.m.RUnlock()
	_, exist := e.ortCfgHash[h]
	return exist
}
