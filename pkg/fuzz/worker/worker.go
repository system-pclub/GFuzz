package worker

import (
	"context"
	"gfuzz/pkg/fuzzer/fuzz"
	"log"
	"strconv"
	"sync"
	"time"
)

// InitWorkers starts maxParallel workers working on inputCh from fuzzer context.
func InitWorkers(maxParallel int, fuzzCtx *fuzz.Context) {
	var wg sync.WaitGroup

	for i := 0; i < maxParallel; i++ {
		wg.Add(1)

		// Start worker
		go func(i int) {
			log.Printf("[Worker %d] Started", i)
			defer wg.Done()
			for {
				select {
				// Receive input
				case task := <-fuzzCtx.runTaskCh:
					log.Printf("[Worker %d] Working on %s\n", i, task.id)
					if ShouldSkipRunTask(fuzzCtx, task) {
						log.Printf("[Worker %d][Task %s] skipped\n", i, task.id)
						continue
					}
					ctx := context.WithValue(context.Background(), "workerID", strconv.Itoa(i))
					result, err := Run(ctx, fuzzCtx, task)
					if err != nil {
						log.Printf("[Worker %d] [Task %s] Error: %s\n", i, task.id, err)
						continue
					}
					err = HandleRunResult(ctx, task, result, fuzzCtx)
					if err != nil {
						log.Printf("[Worker %d] [Task %s] Error: %s\n", i, task.id, err)
						continue
					}
				case <-time.After(3 * time.Minute):
					log.Printf("[Worker %d] Timeout. Exiting...", i)
					return
				}
			}

		}(i)
	}

	wg.Wait()

}
