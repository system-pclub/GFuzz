package selectstrtg

import (
	"runtime"
	"strconv"
)

// SelectCaseInOrder will return case number according to the inputs' order.
type SelectCaseInOrder struct {
	// select identifer (filname + line number) => case recorder
	id2Cr map[string]*SelectCaseRecorder
}

type SelectCaseRecorder struct {
	inputs       []runtime.SelectInfo
	lastInputIdx int
}

// NewSelectCaseInOrder creates a SelectCaseInOrder by given list of inputs.
func NewSelectCaseInOrder(inputs []runtime.SelectInfo) *SelectCaseInOrder {
	var id2Cr map[string]*SelectCaseRecorder = make(map[string]*SelectCaseRecorder)
	for _, input := range inputs {
		selectID := input.StrFileName + ":" + input.StrLineNum
		if _, exist := id2Cr[selectID]; !exist {
			id2Cr[selectID] = &SelectCaseRecorder{
				inputs:       []runtime.SelectInfo{input},
				lastInputIdx: -1,
			}
		} else {
			id2Cr[selectID].inputs = append(id2Cr[selectID].inputs, input)
		}
	}
	return &SelectCaseInOrder{
		id2Cr: id2Cr,
	}
}

// GetCase return the case index application should choose in that select
func (s *SelectCaseInOrder) GetCase(filename string, line, numOfCases int) int {
	lineStr := strconv.Itoa(line)
	selectID := filename + ":" + lineStr
	cr, exist := s.id2Cr[selectID]
	if !exist {
		return -1
	}

	idx := (cr.lastInputIdx + 1) % len(cr.inputs)
	cr.lastInputIdx = idx
	return cr.inputs[idx].IntPrioCase
}
