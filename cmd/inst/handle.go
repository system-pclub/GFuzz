package main

import (
	"gfuzz/pkg/inst"
	"gfuzz/pkg/utils/gofmt"
	"io/ioutil"
	"log"
)

func HandleSrcFile(src string, reg *inst.PassRegistry, passes []string) error {
	iCtx, err := inst.NewInstContext(src)
	if err != nil {
		return err
	}

	err = inst.Run(iCtx, reg, passes)
	if err != nil {
		return err
	}

	var dst string
	if opts.Out != "" {
		dst = opts.Out
	} else {
		// dump AST in-place
		dst = iCtx.File

	}
	err = inst.DumpAstFile(iCtx.FS, iCtx.AstFile, dst)
	if err != nil {
		return err
	}

	// check if output is valid, revert if error happened
	if gofmt.HasSyntaxError(dst) {
		if opts.IgnoreSyntaxErr {
			ioutil.WriteFile(dst, iCtx.OriginalContent, 0666)
			log.Printf("recover '%s' from syntax error\n", dst)
		} else {
			log.Panicf("syntax error found at file '%s'\n", dst)

		}
	}

	return nil
}
