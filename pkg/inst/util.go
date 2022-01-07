package inst

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io/fs"
	"os"

	"golang.org/x/tools/go/ast/astutil"
)

func AddImport(fs *token.FileSet, ast *ast.File, name string, path string) error {
	for _, vecImportSpec := range astutil.Imports(fs, ast) {
		for _, importSpec := range vecImportSpec {
			if importSpec != nil && importSpec.Name != nil {

				if importSpec.Name.Name == name && importSpec.Path.Value == path { // instrumented before
					// ideomatic
					return nil
				}
			}

		}
	}

	ok := astutil.AddNamedImport(fs, ast, name, path)
	if !ok {
		return fmt.Errorf("failed to add import %s %s", name, path)
	}
	return nil
}

// DumpAstFile serialized AST to given file
func DumpAstFile(fset *token.FileSet, astFile *ast.File, dstFile string) error {
	if astFile == nil {
		return fmt.Errorf("found nil ast file for %s", dstFile)
	}
	fi, err := os.Stat(dstFile)
	var mode fs.FileMode
	if err != nil {
		// return any error except not exist
		if !os.IsNotExist(err) {
			return err
		} else {
			// if file not exist, use default mode
			mode = 0666
		}
	} else {
		// if file exist, use same mode
		mode = fi.Mode()
	}
	w, err := os.OpenFile(dstFile, os.O_CREATE|os.O_WRONLY, mode)
	defer w.Close()
	if err != nil {
		return err
	}

	err = format.Node(w, fset, astFile)
	if err != nil {
		return err
	}
	return nil
}
