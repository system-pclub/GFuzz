package mutate

import (
	"gfuzz/pkg/fuzz/gexecfuzz"
	"gfuzz/pkg/oraclert/config"
	"gfuzz/pkg/oraclert/output"
	"gfuzz/pkg/selefcm"
)

type DeterMutateStrategy struct{}

func (d *DeterMutateStrategy) Mutate(g *gexecfuzz.GExecFuzz, curr *config.Config, o *output.Output, energy int) ([]*config.Config, error) {
	if o == nil || o.Selects == nil {
		return nil, nil
	}
	var cfgs []*config.Config

	// loop selects, generate a new config by tweak a different case each time to prioritize
	for selIdx, sel := range o.Selects {
		for i := 0; i < int(sel.Cases); i++ {
			cfg := curr.Copy()
			cfg.SelEfcm.SelTimeout = 500
			cfg.SelEfcm.Efcms = selectsToEfcms(o.Selects)
			cfg.SelEfcm.Efcms[selIdx].Case = i
			cfgs = append(cfgs, cfg)
		}
	}
	return cfgs, nil
}

// selectsToEfcms transfer a list of select records to a list of enforcement
func selectsToEfcms(selects []output.SelectRecord) (efcms []selefcm.SelEfcm) {
	for _, sel := range selects {
		efcms = append(efcms, selefcm.SelEfcm{
			ID:   sel.ID,
			Case: int(sel.Chosen),
		})
	}

	return efcms
}
