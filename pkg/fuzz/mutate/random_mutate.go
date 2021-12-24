package mutate

import (
	"fmt"
	"gfuzz/pkg/fuzz/gexecfuzz"
	"gfuzz/pkg/oraclert/config"
	"gfuzz/pkg/oraclert/output"
	"gfuzz/pkg/selefcm"
	"gfuzz/pkg/utils/rand"
)

type RandomMutateStrategy struct {
	SelEfcmTimeout int
}

func (d *RandomMutateStrategy) Mutate(g *gexecfuzz.GExecFuzz, curr *config.Config, o *output.Output, energy int) ([]*config.Config, error) {
	var cfgs []*config.Config

	// TODO:: If we remove redundant cases, should we count redundant cases into the energy?
	for mutate_idx := 0; mutate_idx < energy; mutate_idx++ {
		cfg := config.NewConfig()
		cfg.SelEfcm.SelTimeout = d.SelEfcmTimeout
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
			numOfSelectCases := records[mutateWhichSelect].Cases
			if numOfSelectCases == 0 {
				return nil, fmt.Errorf("cannot randomly mutate an input with zero number of cases in select %d", mutateWhichSelect)
			}
			mutateToWhatValue := rand.GetRandomWithMax(int(numOfSelectCases))

			cfg.SelEfcm.Efcms = append(cfg.SelEfcm.Efcms, selefcm.SelEfcm{
				ID:   records[mutateWhichSelect].ID,
				Case: mutateToWhatValue,
			})
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
