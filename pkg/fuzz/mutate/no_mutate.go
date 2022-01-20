package mutate

import (
	"gfuzz/pkg/fuzz/gexecfuzz"
	"gfuzz/pkg/oraclert/config"
	"gfuzz/pkg/oraclert/output"
)

type NoMutateStrategy struct {
}

func (d *NoMutateStrategy) Mutate(g *gexecfuzz.GExecFuzz, curr *config.Config, o *output.Output, energy int) ([]*config.Config, error) {
	var cfgs []*config.Config
	for i := 0; i < energy; i++ {
		cfg := config.NewConfig()
		cfgs = append(cfgs, cfg)
	}

	return cfgs, nil
}
