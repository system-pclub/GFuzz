package mutate

import (
	"gfuzz/pkg/selefcm"
	"testing"
)

// func TestMemRandMutateHappy(t *testing.T) {
// 	triggeredCases := make(map[string][]int)

// 	var strat OrtConfigMutateStrategy
// 	strat = &MemRandMutateStrategy{}

// 	gef := gexecfuzz.NewGExecFuzz(nil)
// 	gef.CaseRecords = triggeredCases

// 	currCfg := config.NewConfig()
// 	output := &output.Output{}
// 	genCfgs, err := strat.Mutate(gef, currCfg, output, 5)

// 	if err != nil {
// 		t.Fail()
// 	}

// 	if len(genCfgs) != 5 {
// 		t.Fail()
// 	}

// }

func TestIsEfcmAvailable(t *testing.T) {
	triggeredCases := make(map[string][]int)
	triggeredCases["abc.go:1"] = []int{1}

	if !isEfcmAvailable(triggeredCases, &selefcm.SelEfcm{
		ID:   "abc.go:1",
		Case: 2,
	}) {
		t.Fail()
	}

	if isEfcmAvailable(triggeredCases, &selefcm.SelEfcm{
		ID:   "abc.go:1",
		Case: 1,
	}) {
		t.Fail()
	}

	if !isEfcmAvailable(triggeredCases, &selefcm.SelEfcm{
		ID:   "abcdefg.go:1",
		Case: 1,
	}) {
		t.Fail()
	}
}
