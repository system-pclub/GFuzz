package oraclert

import (
	"gfuzz/pkg/oraclert/env"
	"gfuzz/pkg/selefcm"
	"io/ioutil"
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

var (
	efcmStrat  selefcm.SelectCaseStrategy
	selTimeout int
)

func init() {
	efcmFile := os.Getenv(env.GFUZZ_ENV_SELEFCM_FILE)
	data, err := ioutil.ReadFile(efcmFile)
	if err == nil {
		l, err := selefcm.Deserilize(data)
		if err == nil {
			// We can create different strategies according to our needs
			efcmStrat = selefcm.NewSelectCaseInOrder(l.Efcms)
			selTimeout = l.SelTimeout
		}
	} else {
		println(err)
	}
}

// GetSelEfcmCaseIdx will be instrumented to each select in target program.
func GetSelEfcmSwitchCaseIdx(selectID string) int {
	atomic.AddUint32(&getSelEfcmCount, 1)
	idx := efcmStrat.GetCase(selectID)
	if idx != -1 {
		runtime.StoreLastMySwitchChoice(idx)
		return idx
	} else {
		atomic.AddUint32(&notSelEfcmCount, 1)
		runtime.StoreLastMySwitchChoice(-1)
		return -1 // let switch choose the default case
	}
}

func StoreLastMySwitchChoice(choice int) {
	if choice == -1 {
		atomic.AddUint32(&origSelCount, 1)
	}
	runtime.StoreLastMySwitchChoice(choice)
}

func SelEfcmTimeout() <-chan time.Time {
	// if this channel wins, remember to call "runtime.StoreLastMySwitchChoice(-1)", which means we will use the original select
	return time.After(time.Duration(selTimeout) * time.Millisecond)
}
