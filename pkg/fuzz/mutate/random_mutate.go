package mutate

import (
	"fmt"
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
	IgnoreFeedback      bool
}

func (d *RandomMutateStrategy) Mutate(g *gexecfuzz.GExecFuzz, curr *config.Config, o *output.Output, energy int) ([]*config.Config, error) {
	var cfgs []*config.Config
	idxcache := make(map[string]int)
	// TODO:: If we remove redundant cases, should we count redundant cases into the energy?
	for mutate_idx := 0; mutate_idx < energy; mutate_idx++ {
		cfg := config.NewConfig()
		if !d.FixedSelEfcmTimeout {
			cfg.SelEfcm.SelTimeout = curr.SelEfcm.SelTimeout + 1000
			if cfg.SelEfcm.SelTimeout > 10000 {
				cfg.SelEfcm.SelTimeout = 1000
			}
		} else {
			cfg.SelEfcm.SelTimeout = d.SelEfcmTimeout
		}
		mutateMethod := rand.GetRandomWithMax(10)
		// get all select records we have seen so far for this executable
		records := g.GetAllSelectRecords()
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
				return nil, fmt.Errorf("cannot randomly mutate an input with zero number of cases in select %d", mutateWhichSelect)
			}
			mutateToWhatValue := rand.GetRandomWithMax(int(numOfSelectCases))

			if d.IgnoreFeedback {

				cfg.SelEfcm.Efcms = append(cfg.SelEfcm.Efcms, selefcm.SelEfcm{
					ID:   selectedSel.ID,
					Case: mutateToWhatValue,
				})
			} else {
				// use feedback to avoid random to duplicated case
				selectedInCurr := false
				for _, ef := range curr.SelEfcm.Efcms {
					if ef.ID == selectedSel.ID {
						selectedInCurr = true
						prevIdxOffset := idxcache[ef.ID]
						if prevIdxOffset == -1 {
							// -1 means this ID has been generated all cases
							break
						}
						newCase := (ef.Case + prevIdxOffset + 1) % int(selectedSel.Cases)

						if newCase == ef.Case {
							idxcache[ef.ID] = -1
							break
						}
						log.Println("used feedback to avoid redundent random generation")
						idxcache[ef.ID] += 1
						cfg.SelEfcm.Efcms = append(cfg.SelEfcm.Efcms, selefcm.SelEfcm{
							ID:   selectedSel.ID,
							Case: newCase,
						})
						break
					}
				}

				if !selectedInCurr {
					cfg.SelEfcm.Efcms = append(cfg.SelEfcm.Efcms, selefcm.SelEfcm{
						ID:   selectedSel.ID,
						Case: mutateToWhatValue,
					})
				}
			}

		} else {
			// Mutate random number of select
			mutateChance := rand.GetRandomWithMax(numOfSelects)
			for mutateIdx := 0; mutateIdx < mutateChance; mutateIdx++ {
				mutateWhichSelect := rand.GetRandomWithMax(numOfSelects)
				numOfSelectCases := records[mutateWhichSelect].Cases
				if numOfSelectCases == 0 {
					return nil, fmt.Errorf("cannot randomly mutate an input with zero number of cases in select %d", mutateWhichSelect)
				}
				mutateToWhatValue := rand.GetRandomWithMax(int(numOfSelectCases))
				cfg.SelEfcm.Efcms = append(cfg.SelEfcm.Efcms, selefcm.SelEfcm{
					ID:   records[mutateWhichSelect].ID,
					Case: mutateToWhatValue,
				})
			}
		}

		cfgs = append(cfgs, cfg)

	}

	return cfgs, nil
}
