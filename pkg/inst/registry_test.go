package inst

import "testing"

type MockPass struct{}

func (p *MockPass) Name() string {
	return "mock"
}

func (p *MockPass) Run(iCtx *InstContext) error {
	return nil
}
func TestAddGetRegistryPass(t *testing.T) {
	reg := NewPassRegistry()
	p := &MockPass{}
	reg.AddPass(p)

	gp, err := reg.GetPass("mock")
	if err != nil {
		t.Fail()
	}

	if gp != p {
		t.Fail()
	}
}
func TestHasRegistryPass(t *testing.T) {
	reg := NewPassRegistry()
	p := &MockPass{}
	reg.AddPass(p)

	if reg.HasPass("abcdef") {
		t.Fail()
	}
}
