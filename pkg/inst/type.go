package inst

import (
	"go/ast"
	"go/token"
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

	// Doing analysis or instrumentation for a single Go source file
	Run(iCtx *InstContext) error

	// Deps returns a list of dependent passes
	Deps() []string
}

// PassRegistry records all registered passes
type PassRegistry struct {
	// pass name => pass
	n2p map[string]InstPass
}
