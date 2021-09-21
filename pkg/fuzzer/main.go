package fuzzer

import (
	"gfuzz/pkg/exec"
	"gfuzz/pkg/fuzz"
	"log"
	"os"
)

// Main starts fuzzing with a given list of executables and configuration
func Main(execs []exec.Executable, config *fuzz.Config) {
	if len(execs) == 0 {
		log.Println("no executables found, exit.")
		os.Exit(0)
	}

	for e := range execs {
		log.Println("found executable: %s", e)
	}

	shuffle(execs)

	startWorkers(config.MaxParallel)
}
