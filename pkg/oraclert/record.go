package oraclert

import "runtime"

func StoreOpInfo(strOpType string, uint16OpID uint16) {
	runtime.StoreChOpInfo(strOpType, uint16OpID)
	recordOp(uint16OpID)
}

func StoreChMakeInfo(ch interface{}, uint16OpID uint16) interface{} {
	runtime.StoreChOpInfo("ChMake", uint16OpID)
	runtime.LinkChToLastChanInfo(ch)
	recordOp(uint16OpID)
	return ch
}

func CurrentGoAddPrime(ch interface{}) {
	runtime.CurrentGoAddCh(ch)
}

func CurrentGoAddCh(ch interface{}) {
	runtime.CurrentGoAddCh(ch)
}

func CurrentGoAddWaitgroup(wg interface{}) {
	runtime.CurrentGoAddWaitgroup(wg)
}

func CurrentGoAddMutex(mu interface{}) {
	runtime.CurrentGoAddMutex(mu)
}

func CurrentGoAddCond(cond interface{}) {
	runtime.CurrentGoAddCond(cond)
}
