package mutate

import (
	"gfuzz/pkg/fuzz/gexecfuzz"
	"gfuzz/pkg/oraclert/config"
	"gfuzz/pkg/oraclert/output"
	"gfuzz/pkg/selefcm"
	"gfuzz/pkg/utils/rand"
	"log"
)

// generate missing cases by history records if possible
type MemRandMutateStrategy struct {
	SelEfcmTimeout      int
	FixedSelEfcmTimeout bool
}

func (d *MemRandMutateStrategy) Mutate(g *gexecfuzz.GExecFuzz, curr *config.Config, o *output.Output, energy int) (cfgs []*config.Config, err error) {
	if len(o.Selects) == 0 {
		return nil, nil
	}

	// baseCfg is based on output's select records
	baseCfg := config.NewConfig()
	if !d.FixedSelEfcmTimeout {
		baseCfg.SelEfcm.SelTimeout = curr.SelEfcm.SelTimeout + 1000
		if baseCfg.SelEfcm.SelTimeout > 10500 {
			baseCfg.SelEfcm.SelTimeout = 500
		}
	} else {
		baseCfg.SelEfcm.SelTimeout = d.SelEfcmTimeout
	}
	// copy output's selects into new cfg
	for _, sel := range o.Selects {
		baseCfg.SelEfcm.Efcms = append(baseCfg.SelEfcm.Efcms, selefcm.SelEfcm{
			ID:   sel.ID,
			Case: int(sel.Chosen),
		})
	}

	// make sure we have already record all currrent's efcm
	// since some efcm might be generated from other strategy
	for _, rec := range o.Selects {
		g.RecordCase(rec)
	}

	for idx, sel := range o.Selects {
		efcmToTest := &selefcm.SelEfcm{
			ID: sel.ID,
		}
		for i := 0; i < int(sel.Cases); i++ {
			efcmToTest.Case = i
			if isEfcmAvailable(g.CaseRecords, efcmToTest) {
				cfg := baseCfg.Copy()
				cfg.SelEfcm.Efcms[idx].Case = efcmToTest.Case
				cfgs = append(cfgs, cfg)
			}

			if len(cfgs) == energy {
				return
			}
		}
	}

	// infers not enough number of cfg generated
	diff := energy - len(cfgs)
	numOfSelects := len(o.Selects)
	for i := 0; i < diff; i++ {
		cfg := baseCfg.Copy()

		// Mutate random number of select
		mutateChance := rand.GetRandomWithMax(numOfSelects)
		for mutateIdx := 0; mutateIdx < mutateChance; mutateIdx++ {
			mutateWhichSelect := rand.GetRandomWithMax(numOfSelects)
			selectedSel := o.Selects[mutateWhichSelect]
			numOfSelectCases := selectedSel.Cases
			if numOfSelectCases == 0 {
				log.Printf("cannot randomly mutate an input with zero number of cases in select %d", mutateWhichSelect)
				continue
			}
			cfg.SelEfcm.Efcms[mutateWhichSelect].Case = (int(selectedSel.Chosen) + 1) % int(selectedSel.Cases)
		}
		cfgs = append(cfgs, cfg)
	}
	return
}

func isEfcmAvailable(records map[string][]int, efcm *selefcm.SelEfcm) bool {
	cases, exist := records[efcm.ID]
	if exist {
		for _, c := range cases {
			if c == efcm.Case {
				// return false if we found this case has been generated before
				return false
			}
		}
	}

	return true
}
