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

	for _, e := range execs {
		log.Printf("found executable: %s", e)
	}

	shuffle(execs)

	fuzzCtx := fuzz.NewContext(execs, config)
	eCh := make(chan *fuzz.QueueEntry, config.MaxParallel)

	go func() {
		for {
			// TODO: use interface to handle strategy for next entry
			next := fuzzCtx.NextQueueEntry()
			eCh <- next
		}
	}()

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
				o, err := Run(ctx, i)
				if err != nil {
					logger.Printf("[entry %s] %s", e, err)
				}
				err = fuzz.HandleExecOutput(fuzzCtx, i, o)
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
