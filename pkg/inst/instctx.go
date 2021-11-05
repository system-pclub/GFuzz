package inst

import (
	"go/parser"
	"go/token"
	"io/ioutil"
)

// NewInstContext creates a InstContext by given Golang source file
func NewInstContext(goSrcFile string) (*InstContext, error) {
	oldSource, err := ioutil.ReadFile(goSrcFile)
	if err != nil {
		return nil, err
	}

	fs := token.NewFileSet()
	ast, err := parser.ParseFile(fs, goSrcFile, oldSource, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	return &InstContext{
		File:            goSrcFile,
		OriginalContent: oldSource,
		FS:              fs,
		AstFile:         ast,
	}, nil
}
