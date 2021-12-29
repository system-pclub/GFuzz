package mutate

import (
	"gfuzz/pkg/fuzz/gexecfuzz"
	"gfuzz/pkg/oraclert/config"
	"gfuzz/pkg/oraclert/output"
	"gfuzz/pkg/selefcm"
)

// generate missing cases by history records if possible
type MemRandMutateStrategy struct {
	BackUpStrat OrtConfigMutateStrategy
}

func (d *MemRandMutateStrategy) Mutate(g *gexecfuzz.GExecFuzz, curr *config.Config, o *output.Output, energy int) (cfgs []*config.Config, err error) {
	// baseCfg is based on output's select records
	baseCfg := config.NewConfig()
	baseCfg.SelEfcm.SelTimeout = curr.SelEfcm.SelTimeout
	// copy output's selects into new cfg
	for _, sel := range o.Selects {
		baseCfg.SelEfcm.Efcms = append(baseCfg.SelEfcm.Efcms, selefcm.SelEfcm{
			ID:   sel.ID,
			Case: int(sel.Chosen),
		})
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
			g.RecordSelEfcm(efcmToTest)

			if len(cfgs) == energy {
				return
			}
		}
	}

	// infers not enough number of cfg generated
	diff := energy - len(cfgs)
	otherCfgs, err := d.BackUpStrat.Mutate(g, curr, o, diff)
	if err != nil {
		return nil, err
	}
	cfgs = append(cfgs, otherCfgs...)
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
