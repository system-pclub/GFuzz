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
	cfgs, err := mts.Mutate(nil, curr, output)

	if err != nil {
		t.Fail()
	}

	if len(cfgs) != 7 {
		t.Fail()
	}

}
