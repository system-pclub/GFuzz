package main

import (
	"fmt"
	"gfuzz/pkg/inst"
	"gfuzz/pkg/inst/pass"
	"gfuzz/pkg/inst/stats"
	"gfuzz/pkg/utils/fs"
	"gfuzz/pkg/utils/gofmt"
	"io/ioutil"
	"log"
	"os"
)

var (
	Version string
	Build   string
)

func main() {

	parseFlags()

	if opts.Version {
		fmt.Printf("GFuzz Version: %s Build: %s", Version, Build)
		os.Exit(0)
	}

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
		// default to run all passes if no pass(s) is/are given
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

		var dst string
		if opts.Out != "" {
			dst = opts.Out
		} else {
			// dump AST in-place
			dst = iCtx.File

		}
		err = inst.DumpAstFile(iCtx.FS, iCtx.AstFile, dst)
		if err != nil {
			log.Print(err)
			continue
		}

		// check if output is valid, revert if error happened
		if gofmt.HasSyntaxError(dst) {
			if opts.IgnoreSyntaxErr {
				ioutil.WriteFile(dst, iCtx.OriginalContent, 0666)
			} else {
				log.Panicf("syntax error found at file '%s'", dst)
			}
		}
	}

	// handle output
	if opts.StatsOut != "" {
		err := stats.ToFile(opts.StatsOut)
		if err != nil {
			log.Fatalf("writing stats to %s: %s", opts.StatsOut, err)
		}
	}
}
