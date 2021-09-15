package main

import (
	"gfuzz/pkg/inst"
	"gfuzz/pkg/inst/pass"
	"gfuzz/pkg/utils/fs"

	"log"
	"os"
)

func main() {

	parseFlags()

	reg := inst.NewPassRegistry()

	// register passes
	reg.AddPass(&pass.SelEfcmPass{})
	reg.AddPass(&pass.ChRecPass{})
	reg.AddPass(&pass.CvRecPass{})
	reg.AddPass(&pass.MtxRecPass{})
	reg.AddPass(&pass.WgRecPass{})
	reg.AddPass(&pass.OraclePass{})

	// prepare passes
	var passes []string
	if len(opts.Passes) > 0 {
		passes = opts.Passes
	} else {
		passes = reg.ListOfPassNames()
	}

	// prepare go source files
	var goSrcFiles []string

	if len(opts.Args.Globs) > 0 {
		for _, g := range opts.Args.Globs {
			files, err := fs.ListFilesByGlob(g)
			if err != nil {
				log.Panic(err)
			}
			goSrcFiles = append(goSrcFiles, files...)
		}
	}

	if opts.Dir != "" {
		files, err := listGoSrcByDir(opts.Dir)
		if err != nil {
			log.Panic(err)
		}
		goSrcFiles = append(goSrcFiles, files...)
	}

	if opts.File != "" {
		// TODO: validate file:
		// if exist, if .go
		goSrcFiles = append(goSrcFiles, opts.File)
	}

	if opts.Out != "" && len(goSrcFiles) != 1 {
		log.Panic("--out is only allow with instrumenting single golang source file")
	}

	if len(goSrcFiles) == 0 {
		log.Println("No go source file(s) found")
		os.Exit(0)
	}

	// handle go source files
	// TODO: use goroutine to accelerate
	for _, src := range goSrcFiles {
		iCtx, err := inst.NewInstContext(src)
		if err != nil {
			log.Print(err)
			continue
		}

		err = inst.Run(iCtx, reg, passes)
		if err != nil {
			log.Print(err)
			continue
		}

		if opts.Out != "" {
			err = inst.DumpAstFile(iCtx.FS, iCtx.AstFile, opts.Out)
		} else {
			// dump AST in-place
			err = inst.DumpAstFile(iCtx.FS, iCtx.AstFile, iCtx.File)
		}
		if err != nil {
			log.Print(err)
			continue
		}
	}
}
