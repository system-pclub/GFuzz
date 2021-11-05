package fuzzer

import (
	"context"
	"gfuzz/pkg/fuzz"
	"gfuzz/pkg/fuzz/config"
	"gfuzz/pkg/fuzz/exec"
	"gfuzz/pkg/fuzz/gexecfuzz"
	"gfuzz/pkg/gexec"
	"log"
	"os"
	"time"
)

// Main starts fuzzing with a given list of executables and configuration
func Main(execs []gexec.Executable, config *config.Config) {
	if len(execs) == 0 {
		log.Println("no executables found, exit.")
		os.Exit(0)
	}

	for _, e := range execs {
		log.Printf("found executable: %s", e)
	}

	shuffle(execs)

	fc := fuzz.NewContext(execs, config, &fuzz.HandlerImpl{})
	inputCh := make(chan *exec.Input, config.MaxParallel)

	go func() {
		fc.EachGExecFuzz(func(g *gexecfuzz.GExecFuzz) {
			inputCh <- fuzz.NewInitExecInput(fc, g)
			log.Printf("init %s", g)
		})
	}()

	startWorkers(config.MaxParallel, func(ctx context.Context) {
		queueEntryWorker(ctx, fc, inputCh)
	})

}

// queueEntryWorker handles a queue entry receives from channel
func queueEntryWorker(ctx context.Context, fc *fuzz.Context, ch chan *exec.Input) {
	logger := getWorkerLogger(ctx)
	for {
		select {
		case i := <-ch:
			o, err := Run(ctx, i)
			if err != nil {
				logger.Printf("[input %s] %s", i.ID, err)
			}
			newInputs, err := fc.HandleExec(i, o)
			if err != nil {
				logger.Printf("[input %s] %s", i.ID, err)
			}

			go func() {
				for _, i := range newInputs {
					ch <- i
				}
			}()

		case <-time.After(2 * time.Minute):
			logger.Printf("Timeout. Exited.")
			return
		}
	}
}
