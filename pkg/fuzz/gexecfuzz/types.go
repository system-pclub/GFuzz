package gexecfuzz

import (
	"gfuzz/pkg/gexec"
	"gfuzz/pkg/oraclert/output"
	"gfuzz/pkg/selefcm"
	"gfuzz/pkg/utils/hash"
	"sync"
)

// GExecFuzz all information about fuzzing progress for given executable
type GExecFuzz struct {
	BestScore int
	Exec      gexec.Executable

	// mapping from oracle config hash to timeout count
	EfcmHash2TimeoutCnt map[string]uint32
	// All selects we have seen so far
	OrtSelects       map[string]output.SelectRecord
	OrtChannels      map[string]output.ChanRecord
	OrtTuples        map[uint32]uint32
	InputSelectsHash map[string]struct{}
	m                sync.RWMutex
}

func NewGExecFuzz(exec gexec.Executable) *GExecFuzz {
	return &GExecFuzz{
		Exec:                exec,
		BestScore:           0,
		OrtSelects:          make(map[string]output.SelectRecord),
		OrtChannels:         make(map[string]output.ChanRecord),
		OrtTuples:           make(map[uint32]uint32),
		InputSelectsHash:    make(map[string]struct{}),
		EfcmHash2TimeoutCnt: make(map[string]uint32),
	}
}

func (e *GExecFuzz) String() string {
	return e.Exec.String()
}

func (e *GExecFuzz) UpdateInputSelectEfcmsIfNew(Efcms []selefcm.SelEfcm) int {
	iSelectMap := make(map[string][]int)
	newSelectNum := 0

	for _, v := range Efcms {
		vSelectId := v.ID
		vSelectChosen := v.Case

		if val, ok := iSelectMap[vSelectId]; ok {
			val = append(val, vSelectChosen)
			iSelectMap[vSelectId] = val
		} else {
			val := []int{vSelectChosen}
			iSelectMap[vSelectId] = val
		}
	}

	ortCfgHash := hash.AsSha256(iSelectMap)
	if _, exist := e.InputSelectsHash[ortCfgHash]; !exist {
		e.InputSelectsHash[ortCfgHash] = struct{}{}
		newSelectNum += 1
	}

	return newSelectNum
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

func (e *GExecFuzz) RecordTimeoutEfcm(efcm string) {
	e.m.Lock()
	defer e.m.Unlock()
	e.EfcmHash2TimeoutCnt[efcm] += 1
}

func (e *GExecFuzz) HasTimeoutEfcm(efcm string) bool {
	e.m.RLock()
	defer e.m.RUnlock()
	_, exist := e.EfcmHash2TimeoutCnt[efcm]
	return exist
}
