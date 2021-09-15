package main

import (
	"strings"

	"gfuzz/pkg/utils/fs"
)

func listGoSrcByDir(dir string) ([]string, error) {
	ptn := dir
	if !strings.HasSuffix(dir, "/") {
		ptn = dir + "/"
	}
	ptn = ptn + "**/*.go"
	return fs.ListFilesByGlob(ptn)
}
