package pass

import "gfuzz/pkg/inst"

// SelEfcm, select enforcement pass, instrument the 'select' keyword,
// turn it into a select with multiple cases, each case represent one
// original select's case and a timeout case.
type SelEfcmPass struct{}

func (p *SelEfcmPass) Name() string {
	return "mock"
}

func (p *SelEfcmPass) Run(iCtx *inst.InstContext) error {
	return nil
}
