package fs

import (
	"os"

	"github.com/bmatcuk/doublestar/v4"
)

func ListFilesByGlob(glob string) ([]string, error) {
	fsys := os.DirFS(".")
	matches, err := doublestar.Glob(fsys, glob)
	if err != nil {
		return nil, err
	}
	return matches, nil
}
