package main

import (
	"fmt"
	"gfuzz/pkg/fuzz/api"
	"gfuzz/pkg/fuzz/config"
	"gfuzz/pkg/fuzz/interest"
	"gfuzz/pkg/fuzz/score"
	"gfuzz/pkg/fuzzer"
	gLog "gfuzz/pkg/fuzzer/log"
	"gfuzz/pkg/gexec"
	ortconfig "gfuzz/pkg/oraclert/config"
	"gfuzz/pkg/utils/arr"
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
		log.Fatal("--out is required")
	}

	if _, err := os.Stat(opts.OutputDir); os.IsNotExist(err) {
		err := os.Mkdir(opts.OutputDir, os.ModePerm)
		if err != nil {
			log.Fatalf("create output folder failed: %v", err)
		}
	}

	if opts.GoModDir == "" && opts.TestBinGlobs == nil {
		log.Fatal("Either --gomod or --testbin is required")
	}

	gLog.SetupLogger(filepath.Join(opts.OutputDir, GFUZZ_LOG_FILE), true)

	log.Printf("GFuzz Version: %s Build: %s", Version, Build)

	// prepare fuzzing configuration
	config := config.NewConfig()
	config.OutputDir, err = filepath.Abs(opts.OutputDir)
	if err != nil {
		log.Fatal("filepath.Abs", err)
	}

	// Set up config from opts
	config.MaxParallel = opts.Parallel
	if config.MaxParallel == 0 {
		config.MaxParallel = 5
	}
	log.Printf("Running with MaxParallel: %v", config.MaxParallel)

	config.IsIgnoreFeedback = opts.IsIgnoreFeedback
	config.IsDisableScore = opts.IsDisableScore
	if config.IsIgnoreFeedback {
		log.Printf("Warning: Ignoring feedback from the oracle. ")
		config.IsDisableScore = true
	}

	config.IsNoMutation = opts.IsNoMutation
	if config.IsNoMutation {
		log.Printf("Warning: Do not apply mutations to all test cases. ")
	}

	config.OracleRtDebug = opts.OracleRtDebug

	config.RandMutateEnergy = opts.RandMutateEnergy
	if config.RandMutateEnergy == 0 {
		// Default 5
		config.RandMutateEnergy = 5
		//config.RandMutateEnergy = 1
	}
	log.Printf("Using Random mutation energy: %v", config.RandMutateEnergy)

	if config.IsDisableScore {
		log.Printf("Warning: Disabling score to prioritize fuzzing entries. ")
	} else {
		log.Printf("Using score to prioritize fuzzing entries. ")
	}

	if config.SelEfcmTimeout == 0 {
		config.SelEfcmTimeout = 500
	}

	// prepare fuzz targets
	var execs []gexec.Executable
	if opts.TestBinGlobs != nil {
		log.Printf("list test bin executables from %v", opts.TestBinGlobs)
		execs, err = gexec.ListExecutablesFromTestBinGlobs(opts.TestBinGlobs)
		if err != nil {
			log.Printf("ListExecutablesFromTestBinGlobs: %s", err)
		}
	} else if opts.GoModDir != "" {
		// output directory for compiled test binary file
		binTestsDir, _ := filepath.Abs(filepath.Join(opts.OutputDir, "tbin"))
		execs, err = gexec.ListExecutablesFromGoModule(opts.GoModDir, opts.TestPkg, true, binTestsDir)
		if err != nil {
			log.Printf("ListExecutablesFromGoModule: %s", err)
		}
	}

	// filter fuzz targets by func or package if provided
	var filteredExecs []gexec.Executable
	for _, e := range execs {

		switch v := e.(type) {
		case *gexec.GoBinTest:
			// The reason we don't filter package here are following:
			// 1. test binary file itself cannot tell which package it comes from
			// 2. we already filtered which packages to compile previous
			if opts.TestFunc != nil && !arr.Contains(opts.TestFunc, v.Func) {
				continue
			}
		case *gexec.GoPkgTest:
			if opts.TestFunc != nil && !arr.Contains(opts.TestFunc, v.Func) {
				continue
			}

			if opts.TestPkg != nil && !arr.Contains(opts.TestPkg, v.Package) {
				continue
			}
		}

		filteredExecs = append(filteredExecs, e)
	}

	//fuzzer.Shuffle(filteredExecs)
	fctx := api.NewContext(filteredExecs, config)

	var scorer api.ScoreStrategy = score.NewScoreStrategyImpl(fctx)
	var interestHdl api.InterestHandler = interest.NewInterestHandlerImpl(fctx)

	if opts.Ortconfig != "" {
		ortcfgbytes, err := ioutil.ReadFile(opts.Ortconfig)
		if err != nil {
			fmt.Errorf("read %s: %s", opts.Ortconfig, err)
		}
		ortconfig, err := ortconfig.Deserilize(ortcfgbytes)
		if err != nil {
			fmt.Errorf("parse %s: %s", opts.Ortconfig, err)
		}
		fuzzer.Replay(fctx, filteredExecs[0], config, ortconfig)
		return
	}

	// start fuzzing
	// gLog.DisableStdoutLog()
	// reportCh := make(chan *terminal.TerminalReport)
	// go terminal.Render(reportCh)
	// go terminal.Feed(reportCh, fctx)
	fuzzer.Main(fctx, filteredExecs, config, interestHdl, scorer)
}
