package mutate

import (
	"fmt"
	"gfuzz/pkg/fuzz/gexecfuzz"
	"gfuzz/pkg/oraclert/config"
	"gfuzz/pkg/oraclert/output"
	"gfuzz/pkg/selefcm"
	"gfuzz/pkg/utils/rand"
)

// NfbRandomMutateStrategy is for rand stage with non-feedback mode
type NfbRandomMutateStrategy struct {
	SelEfcmTimeout      int
	FixedSelEfcmTimeout bool
}

func (d *NfbRandomMutateStrategy) Mutate(g *gexecfuzz.GExecFuzz, curr *config.Config, o *output.Output, energy int) ([]*config.Config, error) {
	var cfgs []*config.Config
	for mutate_idx := 0; mutate_idx < energy; mutate_idx++ {
		cfg := config.NewConfig()
		if !d.FixedSelEfcmTimeout {
			cfg.SelEfcm.SelTimeout = curr.SelEfcm.SelTimeout + 1000
			if cfg.SelEfcm.SelTimeout > 10000 {
				cfg.SelEfcm.SelTimeout = 500
			}
		} else {
			cfg.SelEfcm.SelTimeout = d.SelEfcmTimeout
		}
		// get all select records we have seen so far for this execution
		records := o.Selects
		numOfSelects := len(records)
		if numOfSelects == 0 {
			return nil, nil
		}

		// Mutate random number of select
		mutateChance := rand.GetRandomWithMax(numOfSelects)
		for mutateIdx := 0; mutateIdx < mutateChance; mutateIdx++ {
			mutateWhichSelect := rand.GetRandomWithMax(numOfSelects)
			selectedSel := records[mutateWhichSelect]
			randCase := rand.GetRandomWithMax(int(selectedSel.Cases))
			if selectedSel.Cases == 0 {
				return nil, fmt.Errorf("cannot randomly mutate an input with zero number of cases in select %d", mutateWhichSelect)
			}
			cfg.SelEfcm.Efcms = append(cfg.SelEfcm.Efcms, selefcm.SelEfcm{
				ID:   selectedSel.ID,
				Case: randCase,
			})
		}

		cfgs = append(cfgs, cfg)

	}

	return cfgs, nil
}
