package main

import (
	"fmt"
	"gfuzz/pkg/fuzz"
	"gfuzz/pkg/fuzz/config"
	"gfuzz/pkg/fuzzer"
	gLog "gfuzz/pkg/fuzzer/log"
	"gfuzz/pkg/gexec"
	ortconfig "gfuzz/pkg/oraclert/config"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var (
	Version string
	Build   string
)

func main() {
	var err error
	parseFlags()

	// flags sanity check
	if opts.Version {
		fmt.Printf("GFuzz Version: %s Build: %s", Version, Build)
		os.Exit(0)
	}

	if opts.OutputDir == "" {
		log.Fatal("--outputDir is required")
	}

	if opts.InstStats == "" {
		log.Fatal("--instStats is required")
	}

	if _, err := os.Stat(opts.OutputDir); os.IsNotExist(err) {
		err := os.Mkdir(opts.OutputDir, os.ModePerm)
		if err != nil {
			log.Fatalf("create output folder failed: %v", err)
		}
	}

	if opts.GoModDir == "" && opts.TestBinGlobs == nil {
		log.Fatal("Either --goModDir or --testBin is required")
	}

	gLog.SetupLogger(filepath.Join(opts.OutputDir, GFUZZ_LOG_FILE))

	log.Printf("GFuzz Version: %s Build: %s", Version, Build)

	// prepare fuzzing configuration
	config := config.NewConfig()
	config.OutputDir, err = filepath.Abs(opts.OutputDir)
	if err != nil {
		log.Fatal("filepath.Abs", err)
	}
	config.MaxParallel = opts.Parallel

	// prepare fuzz targets
	var execs []gexec.Executable
	if opts.TestFunc != "" {
		if opts.TestPkg == "" && opts.TestBin == "" {
			log.Panicln("if --func is given, either --pkg or --testbin is also required")
		}
		if opts.TestPkg != "" {
			execs = append(execs, &gexec.GoPkgTest{
				Func:     opts.TestFunc,
				Package:  opts.TestPkg,
				GoModDir: opts.GoModDir,
			})
		} else {
			execs = append(execs, &gexec.GoBinTest{
				Func: opts.TestFunc,
				Bin:  opts.TestBin,
			})
		}

	} else if opts.TestPkg != "" {
		execs, err = gexec.ListExecutablesInPackage(opts.GoModDir, opts.TestPkg)
		if err != nil {
			log.Printf("ListExecutablesInPackage: %s", err)
		}
	} else if opts.TestBinGlobs != nil {
		execs, err = gexec.ListExecutablesFromTestBinGlobs(opts.TestBinGlobs)
		if err != nil {
			log.Printf("ListExecutablesFromTestBinGlobs: %s", err)
		}
	} else if opts.GoModDir != "" {
		execs, err = gexec.ListExecutablesFromGoModule(opts.GoModDir)
		if err != nil {
			log.Printf("ListExecutablesFromGoModule: %s", err)
		}
	}

	fuzzer.Shuffle(execs)

	fctx := fuzz.NewContext(execs, config, &fuzz.HandlerImpl{})

	if opts.Ortconfig != "" {
		ortcfgbytes, err := ioutil.ReadFile(opts.Ortconfig)
		if err != nil {
			fmt.Errorf("read %s: %s", opts.Ortconfig, err)
		}
		ortconfig, err := ortconfig.Deserilize(ortcfgbytes)
		if err != nil {
			fmt.Errorf("parse %s: %s", opts.Ortconfig, err)
		}
		fuzzer.Replay(fctx, execs[0], config, ortconfig)
		return
	}

	// start fuzzing
	fuzzer.Main(fctx, execs, config)
}
