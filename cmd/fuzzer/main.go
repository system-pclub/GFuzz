package main

import (
	"fmt"
	"gfuzz/pkg/exec"
	"gfuzz/pkg/fuzz"
	"gfuzz/pkg/fuzzer"
	gLog "gfuzz/pkg/fuzzer/log"
	"log"
	"os"
	"path/filepath"
)

var (
	Version string
	Build   string
)

func main() {
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

	var execs []exec.Executable
	var err error
	if opts.TestBinGlobs != nil {
		execs, err = exec.ListExecutablesFromTestBinGlobs(opts.TestBinGlobs)
		if err != nil {
			log.Println("ListExecutablesFromTestBinGlobs", err)
		}
	} else if opts.GoModDir != "" {
		execs, err = exec.ListExecutablesFromGoModule(opts.GoModDir)
		if err != nil {
			log.Println("ListExecutablesFromGoModule", err)
		}
	}

	// prepare fuzzing configuration
	config := fuzz.NewConfig()
	config.OutputDir, err = filepath.Abs(opts.OutputDir)
	if err != nil {
		log.Fatal("filepath.Abs", err)
	}

	// start fuzzing
	fuzzer.Main(execs, config)
}
