package mutate

import (
	"gfuzz/pkg/fuzz/gexecfuzz"
	"gfuzz/pkg/oraclert/config"
	"gfuzz/pkg/oraclert/output"
	"gfuzz/pkg/selefcm"
	"gfuzz/pkg/utils/rand"
	"log"
)

// NfbRandomMutateStrategy is for rand stage with non-feedback mode
type NfbRandomMutateStrategy struct {
	SelEfcmTimeout      int
	FixedSelEfcmTimeout bool
	RandomTimeoutIncr   bool
}

func (d *NfbRandomMutateStrategy) Mutate(g *gexecfuzz.GExecFuzz, curr *config.Config, o *output.Output, energy int) ([]*config.Config, error) {
	var cfgs []*config.Config
	// baseCfg is based on output's select records
	baseCfg := config.NewConfig()
	// copy output's selects into new cfg
	for _, sel := range o.Selects {
		baseCfg.SelEfcm.Efcms = append(baseCfg.SelEfcm.Efcms, selefcm.SelEfcm{
			ID:   sel.ID,
			Case: int(sel.Chosen),
		})
	}

	for mutate_idx := 0; mutate_idx < energy; mutate_idx++ {
		cfg := baseCfg.Copy()
		if !d.FixedSelEfcmTimeout {
			if d.RandomTimeoutIncr {
				cfg.SelEfcm.SelTimeout = rand.GetRandomWithMax(10500) + 1
			} else {
				cfg.SelEfcm.SelTimeout = curr.SelEfcm.SelTimeout + 1000
				if cfg.SelEfcm.SelTimeout > 10500 {
					cfg.SelEfcm.SelTimeout = 500
				}
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
		// mutateChance := rand.GetRandomWithMax(numOfSelects)
		for mutateIdx := 0; mutateIdx < numOfSelects; mutateIdx++ {
			selectedSel := records[mutateIdx]
			randCase := rand.GetRandomWithMax(int(selectedSel.Cases))
			if selectedSel.Cases == 0 {
				log.Printf("cannot randomly mutate an input with zero number of cases in select %s", selectedSel.ID)
				continue
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
