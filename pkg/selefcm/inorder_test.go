package selefcm

import (
	"testing"
)

func TestNewSelectCaseInOrderHappy(t *testing.T) {
	inputs := []SelEfcm{
		{
			ID:   "abc.go:1",
			Case: 0,
		},
		{
			ID:   "abc.go:1",
			Case: 1,
		},
		{
			ID:   "abc.go:1",
			Case: 2,
		},
	}
	strat := NewSelectCaseInOrder(inputs)
	cr, exist := strat.id2Cr["abc.go:1"]
	if !exist {
		t.FailNow()
	}

	if len(cr.efcms) != 3 {
		t.FailNow()
	}
}

func TestNewSelectCaseInOrderGetCaseHit(t *testing.T) {
	inputs := []SelEfcm{
		{
			ID:   "abc.go:1",
			Case: 0,
		},
		{
			ID:   "abc.go:1",
			Case: 1,
		},
		{
			ID:   "abc.go:1",
			Case: 2,
		},
	}
	strat := NewSelectCaseInOrder(inputs)
	if strat.GetCase("abc.go:1") != 0 {
		t.FailNow()
	}

	if strat.GetCase("abc.go:1") != 1 {
		t.FailNow()
	}

	if strat.GetCase("abc.go:1") != 2 {
		t.FailNow()
	}

	if strat.GetCase("abc.go:1") != 0 {
		t.FailNow()
	}
}

func TestNewSelectCaseInOrderGetCaseNoHit(t *testing.T) {
	inputs := []SelEfcm{}
	strat := NewSelectCaseInOrder(inputs)
	if strat.GetCase("abc.go:1") != -1 {
		t.FailNow()
	}

}
