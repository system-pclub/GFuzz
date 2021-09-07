package main

import (
	"os"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

func listGoSrcByDir(dir string) ([]string, error) {
	ptn := dir
	if !strings.HasSuffix(dir, "/") {
		ptn = dir + "/"
	}
	ptn = ptn + "**/*.go"
	return listGoSrcByGlob(ptn)
}

func listGoSrcByGlob(glob string) ([]string, error) {
	fsys := os.DirFS(".")
	matches, err := doublestar.Glob(fsys, glob)
	if err != nil {
		return nil, err
	}
	return matches, nil
}
