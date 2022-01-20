package inst

import (
	"testing"

	"golang.org/x/tools/go/ast/astutil"
)

type MockPass struct{}

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
	reg.Register("mock", func() InstPass { return &MockPass{} })

	exist := reg.HasPass("mock")
	if !exist {
		t.Fail()
	}

}
