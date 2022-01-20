package main

import (
	"fmt"
	"gfuzz/pkg/inst"
	"gfuzz/pkg/inst/pass"
	"gfuzz/pkg/inst/stats"
	"gfuzz/pkg/utils/fs"
	"log"
	"os"
	"runtime/pprof"
)

var (
	Version      string
	Build        string
	numOfHandled uint32
)

func main() {

	parseFlags()

	if opts.CPUProfile != "" {
		f, err := os.Create(opts.CPUProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()

	}

	if opts.Version {
		fmt.Printf("GFuzz Version: %s Build: %s", Version, Build)
		os.Exit(0)
	}

	reg := inst.NewPassRegistry()

	// register passes
	reg.Register("selefcm", func() inst.InstPass { return &pass.SelEfcmPass{} })
	reg.Register("chrec", func() inst.InstPass { return &pass.ChRecPass{} })
	reg.Register("cvrec", func() inst.InstPass { return &pass.CvRecPass{} })
	reg.Register("mtxrec", func() inst.InstPass { return &pass.MtxRecPass{} })
	reg.Register("wgrec", func() inst.InstPass { return &pass.WgRecPass{} })
	reg.Register("oracle", func() inst.InstPass { return &pass.OraclePass{} })
	reg.Register("chlc", func() inst.InstPass { return pass.NewChLifeCyclePass() })

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
	// var wg sync.WaitGroup
	// toInstSrcCh := make(chan string)
	// for i := 1; i <= int(opts.Parallel); i++ {
	// 	wg.Add(1)
	// 	go func() {
	for _, src := range goSrcFiles {
		err := HandleSrcFile(src, reg, passes)
		if err != nil {
			log.Printf("HandleSrcFile %s: %s", src, err)
		}
	}

	// 		wg.Done()

	// 	}()
	// }

	// for _, src := range goSrcFiles {
	// 	toInstSrcCh <- src
	// }

	// close(toInstSrcCh)

	// wg.Wait()

	// handle output
	if opts.StatsOut != "" {
		err := stats.ToFile(opts.StatsOut)
		if err != nil {
			log.Fatalf("writing stats to %s: %s", opts.StatsOut, err)
		}
	}

	log.Printf("successfully handled %d/%d file(s)", numOfHandled, len(goSrcFiles))
}
