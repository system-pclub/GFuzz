package interest

import (
	"gfuzz/pkg/oraclert/output"
	"gfuzz/pkg/selefcm"
	"testing"
)

func TestIsEfcmCoveredHappy(t *testing.T) {
	efcms := []selefcm.SelEfcm{
		{
			ID:   "abc.go:1",
			Case: 1,
		},
		{
			ID:   "abc.go:3",
			Case: 4,
		},
	}

	records := []output.SelectRecord{
		{
			ID:     "abc.go:1",
			Chosen: 1,
		},
		{
			ID:     "abc.go:2",
			Chosen: 2,
		},
		{
			ID:     "abc.go:3",
			Chosen: 4,
		},
	}

	if !IsEfcmCovered(efcms, records) {
		t.Fail()
	}
}

func TestIsEfcmCoveredUnhappy(t *testing.T) {
	efcms := []selefcm.SelEfcm{
		{
			ID:   "abc.go:1",
			Case: 1,
		},
		{
			ID:   "abc.go:3",
			Case: 4,
		},
	}

	records := []output.SelectRecord{
		{
			ID:     "abc.go:1",
			Chosen: 1,
		},
		{
			ID:     "abc.go:2",
			Chosen: 2,
		},
		{
			ID:     "abc.go:3",
			Chosen: 3,
		},
	}

	if IsEfcmCovered(efcms, records) {
		t.Fail()
	}
}
