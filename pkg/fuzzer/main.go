package fuzzer

import (
	"context"
	"gfuzz/pkg/fuzz/api"
	"gfuzz/pkg/fuzz/config"
	"gfuzz/pkg/fuzz/gexecfuzz"
	"gfuzz/pkg/gexec"
	ortconfig "gfuzz/pkg/oraclert/config"
	"log"
	"os"
	"time"
)

// Reply run the fuzzing with given oracle runtime configuration and given executable
func Replay(fctx *api.Context, ge gexec.Executable, config *config.Config, rtConfig *ortconfig.Config) {
	ctx := context.Background()
	i := api.NewExecInput(fctx.GetAutoIncGlobalID(), 0, config.OutputDir, ge, rtConfig, api.ReplayStage)
	Run(ctx, i)
}

// Main starts fuzzing with a given list of executables and configuration
func Main(fctx *api.Context, execs []gexec.Executable, config *config.Config,
	interestHdl api.InterestHandler, scorer api.ScoreStrategy) {
	if len(execs) == 0 {
		log.Println("no executables found, exit.")
		os.Exit(0)
	}

	for _, e := range execs {
		log.Printf("found executable: %s", e)
	}

	// initialize interested inputs by generating init stage input for each executables
	fctx.EachGExecFuzz(func(g *gexecfuzz.GExecFuzz) {
		i := api.NewInitExecInput(fctx, g.Exec)
		fctx.Interests.Add(api.NewUnexecutedInterestInput(i))
	})

	// endless loop to handle interested inputs
	go func() {
		for {
			// handle interested inputs one by one
			fctx.Interests.Each(interestHdl)
		}
	}()

	// start a group of workers to handle fuzz execution in parallel
	startWorkers(config.MaxParallel, func(ctx context.Context) {
		execWorker(ctx, fctx, scorer)
	})

}

// execWorker handles a execution inputs from channel
func execWorker(ctx context.Context, fc *api.Context, scorer api.ScoreStrategy) {
	logger := getWorkerLogger(ctx)

	for {
		select {
		case i := <-fc.ExecInputCh:
			logger.Printf("start %s", i.ID)
			o, err := Run(ctx, i)
			if err != nil {
				logger.Printf("%s: %s", i.ID, err)
			}
			err = HandleExec(i, o, fc, scorer)
			if err != nil {
				logger.Printf("%s: %s", i.ID, err)
			}
		case <-time.After(2 * time.Minute):
			logger.Printf("Timeout. Exited.")
			return
		}
	}
}
