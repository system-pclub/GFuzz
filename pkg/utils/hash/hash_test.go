package hash

import (
	"gfuzz/pkg/oraclert/config"
	"gfuzz/pkg/selefcm"
	"testing"
)

func TestHashEqInOracleRtConfig1(t *testing.T) {
	cfg1 := config.NewConfig()
	cfg2 := config.NewConfig()

	cfg1.SelEfcm.Copy().Efcms = append(cfg1.SelEfcm.Copy().Efcms, selefcm.SelEfcm{
		ID:   "1",
		Case: 2,
	})
	cfg1.SelEfcm.Copy().Efcms = append(cfg1.SelEfcm.Copy().Efcms, selefcm.SelEfcm{
		ID:   "2",
		Case: 4,
	})

	cfg2.SelEfcm.Copy().Efcms = append(cfg1.SelEfcm.Copy().Efcms, selefcm.SelEfcm{
		ID:   "2",
		Case: 4,
	})
	cfg2.SelEfcm.Copy().Efcms = append(cfg1.SelEfcm.Copy().Efcms, selefcm.SelEfcm{
		ID:   "1",
		Case: 2,
	})

	if AsSha256(cfg1) != AsSha256(cfg2) {
		t.Fail()
	}
}

func TestHashEqInOracleRtConfig2(t *testing.T) {
	cfg1 := config.NewConfig()
	cfg2 := config.NewConfig()

	cfg1.SelEfcm.Copy().Efcms = append(cfg1.SelEfcm.Copy().Efcms, selefcm.SelEfcm{
		ID:   "1",
		Case: 2,
	})
	cfg1.SelEfcm.Copy().Efcms = append(cfg1.SelEfcm.Copy().Efcms, selefcm.SelEfcm{
		ID:   "2",
		Case: 4,
	})

	cfg2.SelEfcm.Copy().Efcms = append(cfg1.SelEfcm.Copy().Efcms, selefcm.SelEfcm{
		ID:   "1",
		Case: 2,
	})
	cfg2.SelEfcm.Copy().Efcms = append(cfg1.SelEfcm.Copy().Efcms, selefcm.SelEfcm{
		ID:   "2",
		Case: 4,
	})

	if AsSha256(cfg1) != AsSha256(cfg2) {
		t.Fail()
	}
}
