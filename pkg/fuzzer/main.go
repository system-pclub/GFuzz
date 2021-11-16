package fuzzer

import (
	"context"
	"gfuzz/pkg/fuzz"
	"gfuzz/pkg/fuzz/config"
	"gfuzz/pkg/fuzz/exec"
	"gfuzz/pkg/fuzz/gexecfuzz"
	"gfuzz/pkg/gexec"
	ortconfig "gfuzz/pkg/oraclert/config"
	"log"
	"os"
	"time"
)

func Replay(fctx *fuzz.Context, ge gexec.Executable, config *config.Config, rtConfig *ortconfig.Config) {
	ctx := context.Background()
	i := fuzz.NewExecInput(fctx.GetAutoIncGlobalID(), 0, config.OutputDir, ge, rtConfig, exec.ReplayStage)
	Run(ctx, i)
}

// Main starts fuzzing with a given list of executables and configuration
func Main(fctx *fuzz.Context, execs []gexec.Executable, config *config.Config) {
	if len(execs) == 0 {
		log.Println("no executables found, exit.")
		os.Exit(0)
	}

	for _, e := range execs {
		log.Printf("found executable: %s", e)
	}

	inputCh := make(chan *exec.Input)

	go func() {
		fctx.EachGExecFuzz(func(g *gexecfuzz.GExecFuzz) {
			inputCh <- fuzz.NewInitExecInput(fctx, g.Exec)
			log.Printf("init %s", g)
		})
	}()

	startWorkers(config.MaxParallel, func(ctx context.Context) {
		queueEntryWorker(ctx, fctx, inputCh)
	})

}

// queueEntryWorker handles a queue entry receives from channel
func queueEntryWorker(ctx context.Context, fc *fuzz.Context, ch chan *exec.Input) {
	logger := getWorkerLogger(ctx)

	for {
		select {
		case i := <-ch:
			logger.Printf("start %s", i.ID)
			o, err := Run(ctx, i)
			if err != nil {
				logger.Printf("%s: %s", i.ID, err)
			}
			newInputs, err := fc.HandleExec(i, o)
			if err != nil {
				logger.Printf("%s: %s", i.ID, err)
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
