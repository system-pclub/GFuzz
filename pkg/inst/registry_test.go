package inst

import (
	"testing"

	"golang.org/x/tools/go/ast/astutil"
)

type MockPass struct{}

func (p *MockPass) Name() string {
	return "mock"
}

func (p *MockPass) Deps() []string {
	return nil
}

func (p *MockPass) Before(*InstContext) {
}

func (p *MockPass) After(*InstContext) {
}

func (p *MockPass) GetPreApply(*InstContext) func(*astutil.Cursor) bool {
	return nil
}

func (p *MockPass) GetPostApply(*InstContext) func(*astutil.Cursor) bool {
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
