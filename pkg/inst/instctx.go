package inst

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"log"
)

// NewInstContext creates a InstContext by given Golang source file
func NewInstContext(goSrcFile string) (*InstContext, error) {
	oldSource, err := ioutil.ReadFile(goSrcFile)
	if err != nil {
		return nil, err
	}

	fs := token.NewFileSet()
	astF, err := parser.ParseFile(fs, goSrcFile, oldSource, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	conf := types.Config{
		Importer: importer.ForCompiler(fs, "source", nil),
		Error: func(err error) {
			log.Printf("'%s' type checker: %s", goSrcFile, err)
		},
	}
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}

	conf.Check(astF.Name.Name, fs, []*ast.File{astF}, info)

	return &InstContext{
		File:            goSrcFile,
		OriginalContent: oldSource,
		FS:              fs,
		Type:            info,
		AstFile:         astF,
		Metadata:        make(map[string]interface{}),
	}, nil
}
