package fuzzer

import (
	"context"
	"gfuzz/pkg/exec"
	"gfuzz/pkg/fuzz"
	"log"
	"os"
	"time"
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

	fuzzCtx := fuzz.NewContext(execs)
	eCh := make(chan *fuzz.QueueEntry, config.MaxParallel)

	startWorkers(config.MaxParallel, func(ctx context.Context) {
		queueEntryWorker(ctx, fuzzCtx, eCh)
	})
}

// queueEntryWorker handles a queue entry receives from channel
func queueEntryWorker(ctx context.Context, fuzzCtx *fuzz.Context, eCh chan *fuzz.QueueEntry) {
	logger := getWorkerLogger(ctx)

	for {
		select {
		case e := <-eCh:
			inputs, err := fuzz.HandleQueueEntry(ctx, fuzzCtx, e)
			if err != nil {
				logger.Printf("[entry %s] %s", e, err)
			}

			for _, i := range inputs {
				out, err := Run(ctx, i)
				if err != nil {
					logger.Printf("[entry %s] %s", e, err)
				}
			}

		case <-time.After(2 * time.Minute):
			logger.Printf("Timeout. Exited.")
			return
		}
	}
}
