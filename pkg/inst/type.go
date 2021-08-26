package inst

import (
	"go/ast"
)

// InstContext contains all information needed to instrument one single Golang source code.
type InstContext struct {
	FileName string
	Ast      ast.File
}

// InstPass shapes the pass used for instrumenting a single Golang source code
type InstPass interface {
	Name() string
	Run(iCtx *InstContext) error
}

// PassRegistry records all registered passes
type PassRegistry struct {
	// pass name => pass
	n2p map[string]InstPass
}
