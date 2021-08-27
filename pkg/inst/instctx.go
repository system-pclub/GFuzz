package inst

import (
	"bytes"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
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
		File:    goSrcFile,
		FS:      fs,
		AstFile: ast,
	}, nil
}

func (i *InstContext) Dump() error {
	buf := &bytes.Buffer{}
	err := format.Node(buf, i.FS, i.AstFile)
	if err != nil {
		return fmt.Errorf("error formatting new code: %s in file:%s", err.Error(), i.File)
	}

	newSource := buf.Bytes()
	fi, err := os.Stat(i.File)
	if err != nil {
		return fmt.Errorf("Error in os.Stat file: %s\tError:%s", i.File, err.Error())
	}
	err = ioutil.WriteFile(i.File, newSource, fi.Mode())
	return err
}
