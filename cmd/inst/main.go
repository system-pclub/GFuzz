package main

import (
	"gfuzz/pkg/inst"
	"gfuzz/pkg/inst/pass"
	"log"
)

func main() {
	parseFlags()

	reg := inst.NewPassRegistry()

	// register passes
	reg.AddPass(&pass.SelEfcmPass{})

	if len(opts.Args.Globs) == 0 {
		log.Panic("at least one glob need to present")
	}

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
			files, err := listGoSrcByGlob(g)
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

	// handle go source files
	for _, src := range goSrcFiles {
		iCtx, err := inst.NewInstContext(src)
		if err != nil {
			log.Panic(err)
		}

		err = inst.Run(iCtx, reg, passes)
		if err != nil {
			log.Panic(err)
		}
	}
}
