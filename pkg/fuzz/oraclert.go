package fuzz

import (
	"log"
)

func ShouldSkipInput(fuzzCtx *FuzzContext, exec string, i *OracleRtInput) bool {
	if i == nil {
		return true
	}

	// drop it if list of selects has been run more than 3 times
	if GetCountsOfSelects(exec, i.VecSelect) > 5 {
		log.Printf("drop %s since the select inputs run more than 3 times", i)
		return false
	}

	// drop it if it's source already timeout more than 3 times
	fuzzCtx.timeoutTargetsLock.RLock()
	timeoutCnt, exist := fuzzCtx.timeoutTargets[exec]
	fuzzCtx.timeoutTargetsLock.RUnlock()
	if exist {
		if timeoutCnt > 3 {
			log.Printf("drop %s since it has timeout more than 3 times", i)
			return true
		}
	}

	return false
}

func (i *OracleRtInput) DeepCopy() *OracleRtInput {
	newInput := &OracleRtInput{
		Note:          i.Note,
		SelectDelayMS: i.SelectDelayMS,
		VecSelect:     []SelectInput{},
	}
	for _, selectInput := range i.VecSelect {
		newInput.VecSelect = append(newInput.VecSelect, copySelectInput(selectInput))
	}
	return newInput
}

func FindRecordHashInSlice(recordHash string, recordHashSlice []string) bool {
	for _, searchRecordHash := range recordHashSlice {
		if recordHash == searchRecordHash {
			return true
		}
	}
	return false
}

func copySelectInput(sI SelectInput) SelectInput {
	return SelectInput{
		StrFileName: sI.StrFileName,
		IntLineNum:  sI.IntLineNum,
		IntNumCase:  sI.IntNumCase,
		IntPrioCase: sI.IntPrioCase,
	}
}
