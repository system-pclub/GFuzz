package mutate

import (
	"gfuzz/pkg/oraclert/config"
	"gfuzz/pkg/oraclert/output"
	"testing"
)

func TestDeterMutateStrategy(t *testing.T) {
	var mts OrtConfigMutateStrategy = &DeterMutateStrategy{}
	curr := config.NewConfig()
	curr.DumpSelects = true

	output := &output.Output{
		Selects: []output.SelectRecord{
			{
				ID:     "abc.go:1",
				Cases:  5,
				Chosen: 1,
			},
			{
				ID:     "abc.go:2",
				Cases:  2,
				Chosen: 1,
			},
		},
	}
	cfgs, err := mts.Mutate(nil, curr, output, 0)

	if err != nil {
		t.Fail()
	}

	if len(cfgs) != 7 {
		t.Fail()
	}

}

func TestSelectsToEfcms(t *testing.T) {
	selects := []output.SelectRecord{
		{
			ID:     "abc.go:1",
			Cases:  5,
			Chosen: 1,
		},
		{
			ID:     "abc.go:2",
			Cases:  2,
			Chosen: 0,
		},
	}

	efcms := selectsToEfcms(selects)

	if efcms[0].ID != "abc.go:1" {
		t.Fail()
	}

	if efcms[0].Case != 1 {
		t.Fail()
	}

	if efcms[1].ID != "abc.go:2" {
		t.Fail()
	}

	if efcms[1].Case != 0 {
		t.Fail()
	}
}
