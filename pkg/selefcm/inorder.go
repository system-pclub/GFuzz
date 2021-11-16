package selefcm

import (
	"sync/atomic"
)

// SelectCaseInOrder will return case number according to the inputs' order.
type SelectCaseInOrder struct {
	// id2Cr is map from: select identifer (filname + line number) => case recorder
	// We don't need to add mutex is here since it aims to be read-only after initialization:
	// means no new SelectCaseRecorder will be added.
	id2Cr map[string]*SelectCaseRecorder
}

type SelectCaseRecorder struct {
	// A list of enforcements for same select ID
	efcms []SelEfcm
	// Last returned element whose index in efcms array
	lastEfcmIdx int32
}

// NewSelectCaseInOrder creates a SelectCaseInOrder by given list of inputs.
func NewSelectCaseInOrder(efcms []SelEfcm) *SelectCaseInOrder {
	var id2Cr map[string]*SelectCaseRecorder = make(map[string]*SelectCaseRecorder)
	for _, efcm := range efcms {
		if _, exist := id2Cr[efcm.ID]; !exist {
			id2Cr[efcm.ID] = &SelectCaseRecorder{
				efcms:       []SelEfcm{efcm},
				lastEfcmIdx: -1,
			}
		} else {
			id2Cr[efcm.ID].efcms = append(id2Cr[efcm.ID].efcms, efcm)
		}
	}
	return &SelectCaseInOrder{
		id2Cr: id2Cr,
	}
}

// GetCase return the case index application should choose in that select
func (s *SelectCaseInOrder) GetCase(selectID string) int {
	cr, exist := s.id2Cr[selectID]
	if !exist {
		return -1
	}
	var newIdx int32

	// lock-free update lastInputIdx
	for {
		oldIdx := atomic.LoadInt32(&cr.lastEfcmIdx)
		newIdx = (oldIdx + 1) % int32(len(cr.efcms))
		res := atomic.CompareAndSwapInt32(&cr.lastEfcmIdx, oldIdx, newIdx)
		if res {
			break
		}
	}

	return cr.efcms[newIdx].Case
}
