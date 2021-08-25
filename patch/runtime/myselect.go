package runtime

import "sync/atomic"

var MapSelectInfo map[string]SelectInfo // useful only when RecordSelectChoice is true
//var MapInput map[string]SelectInfo // useful only when RecordSelectChoice is false
var RecordSelectChoice bool = false
var MuFirstInput mutex
var Uint32SelectCount uint32
var BoolSelectCount bool



type SelectInfo struct {
	StrFileName string
	StrLineNum string
	IntNumCase int
	IntPrioCase int
}

func StoreSelectInput(intNumCase, intChosenCase int) {
	newSelectInput := NewSelectInputFromRuntime(intNumCase, intChosenCase, 3)
	if newSelectInput.IntNumCase != 0 { // IntNumCase would be 0 if this select is not instrumented (e.g., in SDK) and we can't mutate it
		lock(&MuFirstInput)
		MapSelectInfo[newSelectInput.StrFileName + ":" + newSelectInput.StrLineNum] = newSelectInput
		unlock(&MuFirstInput)
	}
}

func NewSelectInputFromRuntime(intNumCase, intPrioCase int, intLayerCallee int) SelectInfo {
	// if A contains select, select calls StoreSelectInput, StoreSelectInput calls this function, then intLayerCallee is 3
	const size = 64 << 10
	buf := make([]byte, size)
	buf = buf[:Stack(buf, false)]
	strStack := string(buf)
	stackSingleGo := ParseStackStr(strStack)
	if len(stackSingleGo.VecFuncLine) < intLayerCallee {
		return SelectInfo{}
	}
	selectInput := SelectInfo{
		StrFileName: stackSingleGo.VecFuncFile[intLayerCallee - 1], // where is the select
		StrLineNum:  LastMySwitchLineNum(),
		IntNumCase:  LastMySwitchOriSelectNumCase(),
		IntPrioCase: -1,
	}

	if LastMySwitchChoice() == -1 {
		// Executing the original select
		selectInput.IntPrioCase = intPrioCase
	} else {
		// Executing our select, so the chosen case is same to switch's choice
		selectInput.IntPrioCase = LastMySwitchChoice()
	}

	return selectInput
}

func LastMySwitchLineNum() string {
	if getg().lastMySwitchLineNum != "" {
		return getg().lastMySwitchLineNum
	} else {
		return "0"
	}
}

func LastMySwitchOriSelectNumCase() int {
	return getg().lastMySwitchOriSelectNumCase
}

func LastMySwitchChoice() int {
	return getg().lastMySwitchChoice
}

func StoreLastMySwitchSelectNumCase(numCase int) {
	getg().lastMySwitchOriSelectNumCase = numCase
}

func StoreLastMySwitchChoice(choice int) {
	getg().lastMySwitchChoice = choice
}

func StoreLastMySwitchLineNum(strLine string) {
	getg().lastMySwitchLineNum = strLine // no need for synchronization.
}

func SelectCount() {
	atomic.AddUint32(&Uint32SelectCount, 1)
}

func ReadSelectCount() uint32 {
	return atomic.LoadUint32(&Uint32SelectCount)
}