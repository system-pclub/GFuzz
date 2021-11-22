package fs

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

func ListFilesByGlob(glob string) ([]string, error) {
	var fsys fs.FS
	startWithRoot := false
	if strings.HasPrefix(glob, "/") {
		fsys = os.DirFS("/")
		startWithRoot = true
		glob = glob[1:]
	} else {
		fsys = os.DirFS(".")
	}
	fmt.Printf("%s, %s\n", fsys, glob)
	matches, err := doublestar.Glob(fsys, glob)
	if err != nil {
		return nil, err
	}
	if startWithRoot {
		var adjustedMaches []string
		for _, m := range matches {
			adjustedMaches = append(adjustedMaches, "/"+m)
		}
		return adjustedMaches, nil
	}
	return matches, nil
}
