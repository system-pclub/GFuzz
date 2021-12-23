package gexecfuzz

import (
	"gfuzz/pkg/oraclert/config"
	"gfuzz/pkg/selefcm"
	"gfuzz/pkg/utils/hash"
	"testing"
)

func TestHasOrtCfgHash(t *testing.T) {

	gexec := NewGExecFuzz(nil)
	cfg1 := &config.Config{
		DumpSelects: true,
		SelEfcm: selefcm.SelEfcmConfig{
			SelTimeout: 500,
			Efcms: []selefcm.SelEfcm{
				{
					ID:   "a.go",
					Case: 2,
				},
			},
		},
	}
	cfg2 := &config.Config{
		DumpSelects: true,
		SelEfcm: selefcm.SelEfcmConfig{
			SelTimeout: 500,
			Efcms: []selefcm.SelEfcm{
				{
					ID:   "a.go",
					Case: 1,
				},
			},
		},
	}
	dupCfg2 := &config.Config{
		DumpSelects: true,
		SelEfcm: selefcm.SelEfcmConfig{
			SelTimeout: 500,
			Efcms: []selefcm.SelEfcm{
				{
					ID:   "a.go",
					Case: 1,
				},
			},
		},
	}
	gexec.RecordOrtCfgHash(hash.AsSha256(cfg1))
	gexec.RecordOrtCfgHash(hash.AsSha256(cfg2))

	if !gexec.HasOrtCfgHash(hash.AsSha256(dupCfg2)) {
		t.Fail()
	}
}
