package mutate

import (
	"gfuzz/pkg/fuzz/gexecfuzz"
	"gfuzz/pkg/oraclert/config"
	"gfuzz/pkg/oraclert/output"
	"gfuzz/pkg/selefcm"
	"gfuzz/pkg/utils/rand"
	"log"
)

type RandomMutateStrategy struct {
	SelEfcmTimeout      int
	FixedSelEfcmTimeout bool
}

func (d *RandomMutateStrategy) Mutate(g *gexecfuzz.GExecFuzz, curr *config.Config, o *output.Output, energy int) ([]*config.Config, error) {
	var cfgs []*config.Config
	idxcache := make(map[string]int)

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
			cfg.SelEfcm.SelTimeout = curr.SelEfcm.SelTimeout + 1000
			if cfg.SelEfcm.SelTimeout > 10500 {
				cfg.SelEfcm.SelTimeout = 500
			}
		} else {
			cfg.SelEfcm.SelTimeout = d.SelEfcmTimeout
		}
		mutateMethod := rand.GetRandomWithMax(10)
		// get all select records we have seen so far for this execution
		records := o.Selects
		numOfSelects := len(records)
		if numOfSelects == 0 {
			return nil, nil
		}

		if mutateMethod < 8 {
			// Mutate one select per time
			mutateWhichSelect := rand.GetRandomWithMax(numOfSelects)
			selectedSel := records[mutateWhichSelect]
			numOfSelectCases := selectedSel.Cases
			if numOfSelectCases == 0 {
				log.Printf("cannot randomly mutate an input with zero number of cases in select %d", mutateWhichSelect)
				continue
			}

			// use feedback to avoid random to duplicated case
			for _, rec := range records {
				if rec.ID == selectedSel.ID {
					log.Println("used feedback to avoid redundent random generation")
					prevIdxOffset := idxcache[rec.ID]
					if prevIdxOffset == -1 {
						// -1 means this ID has been generated all cases
						break
					}
					newCase := (int(rec.Chosen) + prevIdxOffset + 1) % int(selectedSel.Cases)

					if newCase == int(rec.Chosen) {
						idxcache[rec.ID] = -1
						break
					}
					idxcache[rec.ID] += 1
					cfg.SelEfcm.Efcms[mutateWhichSelect].Case = newCase
					break
				}
			}

		} else {
			// Mutate random number of select
			mutateChance := rand.GetRandomWithMax(numOfSelects)
			for mutateIdx := 0; mutateIdx < mutateChance; mutateIdx++ {
				mutateWhichSelect := rand.GetRandomWithMax(numOfSelects)
				selectedSel := records[mutateWhichSelect]
				numOfSelectCases := selectedSel.Cases
				if numOfSelectCases == 0 {
					log.Printf("cannot randomly mutate an input with zero number of cases in select %d", mutateWhichSelect)
					continue
				}
				cfg.SelEfcm.Efcms[mutateWhichSelect].Case = (int(selectedSel.Chosen) + 1) % int(selectedSel.Cases)
			}
		}

		cfgs = append(cfgs, cfg)

	}

	return cfgs, nil
}
