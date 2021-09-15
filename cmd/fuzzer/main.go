package main

import (
	"fmt"
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

	if opts.Version {
		fmt.Printf("GFuzz Version: %s Build: %s", Version, Build)
		os.Exit(0)
	}

	if opts.OutputDir == "" {
		log.Fatal("--outputDir is required")
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

	setupLogger(filepath.Join(opts.OutputDir, GFUZZ_LOG_FILE))

	log.Printf("GFuzz Version: %s Build: %s", Version, Build)

}
