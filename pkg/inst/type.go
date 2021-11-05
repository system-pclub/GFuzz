package inst

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
)

// InstContext contains all information needed to instrument one single Golang source code.
type InstContext struct {
	File            string
	OriginalContent []byte
	FS              *token.FileSet
	AstFile         *ast.File
}

// InstPass shapes the pass used for instrumenting a single Golang source code
type InstPass interface {

	// Name returns the name of the pass
	Name() string

	// Deps returns a list of dependent passes
	Deps() []string

	GetPreApply(iCtx *InstContext) func(*astutil.Cursor) bool

	GetPostApply(iCtx *InstContext) func(*astutil.Cursor) bool
}

// PassRegistry records all registered passes
type PassRegistry struct {
	// pass name => pass
	n2p map[string]InstPass
}
