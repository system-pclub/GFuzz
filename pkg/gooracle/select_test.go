package gooracle

import (
	"runtime"
	"testing"
)

func TestNewSelectCaseInOrderHappy(t *testing.T) {
	inputs := []runtime.SelectInfo{
		runtime.SelectInfo{
			StrFileName: "abc.go",
			StrLineNum:  "1",
			IntNumCase:  3,
			IntPrioCase: 0,
		},
		runtime.SelectInfo{
			StrFileName: "abc.go",
			StrLineNum:  "1",
			IntNumCase:  3,
			IntPrioCase: 1,
		},
		runtime.SelectInfo{
			StrFileName: "abc.go",
			StrLineNum:  "1",
			IntNumCase:  3,
			IntPrioCase: 2,
		},
	}
	strat := NewSelectCaseInOrder(inputs)
	cr, exist := strat.id2Cr["abc.go:1"]
	if !exist {
		t.FailNow()
	}

	if len(cr.inputs) != 3 {
		t.FailNow()
	}
}

func TestNewSelectCaseInOrderGetCaseHit(t *testing.T) {
	inputs := []runtime.SelectInfo{
		runtime.SelectInfo{
			StrFileName: "abc.go",
			StrLineNum:  "1",
			IntNumCase:  3,
			IntPrioCase: 0,
		},
		runtime.SelectInfo{
			StrFileName: "abc.go",
			StrLineNum:  "1",
			IntNumCase:  3,
			IntPrioCase: 1,
		},
		runtime.SelectInfo{
			StrFileName: "abc.go",
			StrLineNum:  "1",
			IntNumCase:  3,
			IntPrioCase: 2,
		},
	}
	strat := NewSelectCaseInOrder(inputs)
	if strat.GetCase("abc.go", 1, 3) != 0 {
		t.FailNow()
	}

	if strat.GetCase("abc.go", 1, 3) != 1 {
		t.FailNow()
	}

	if strat.GetCase("abc.go", 1, 3) != 2 {
		t.FailNow()
	}

	if strat.GetCase("abc.go", 1, 3) != 0 {
		t.FailNow()
	}
}

func TestNewSelectCaseInOrderGetCaseNoHit(t *testing.T) {
	inputs := []runtime.SelectInfo{}
	strat := NewSelectCaseInOrder(inputs)
	if strat.GetCase("abc.go", 1, 3) != -1 {
		t.FailNow()
	}

}
