package interest

import (
	"gfuzz/pkg/oraclert/config"
	"gfuzz/pkg/selefcm"
	"gfuzz/pkg/utils/hash"
	"testing"
)

func TestCfgHashEq(t *testing.T) {
	cfg1 := &config.Config{
		SelEfcm: selefcm.SelEfcmConfig{
			SelTimeout: 500,
			Efcms: []selefcm.SelEfcm{
				{
					ID:   "abc.go:123",
					Case: 1,
				},
			},
		},
	}
	cfg2 := &config.Config{
		SelEfcm: selefcm.SelEfcmConfig{
			SelTimeout: 500,
			Efcms: []selefcm.SelEfcm{
				{
					ID:   "abc.go:123",
					Case: 1,
				},
			},
		},
	}
	if hash.AsSha256(cfg1) != hash.AsSha256(cfg2) {
		t.Fail()
	}
}

func TestCfgHashNotEq(t *testing.T) {
	cfg1 := &config.Config{
		SelEfcm: selefcm.SelEfcmConfig{
			SelTimeout: 1000,
			Efcms: []selefcm.SelEfcm{
				{
					ID:   "abc.go:123",
					Case: 1,
				},
			},
		},
	}
	cfg2 := &config.Config{
		SelEfcm: selefcm.SelEfcmConfig{
			SelTimeout: 500,
			Efcms: []selefcm.SelEfcm{
				{
					ID:   "abc.go:123",
					Case: 1,
				},
			},
		},
	}
	if hash.AsSha256(cfg1) == hash.AsSha256(cfg2) {
		t.Fail()
	}
}
