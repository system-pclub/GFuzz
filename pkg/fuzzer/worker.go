package fuzzer

import (
	"context"
	"fmt"
	"gfuzz/pkg/fuzz"
	gLog "gfuzz/pkg/fuzzer/log"
	"log"
	"strconv"
	"sync"
	"time"
)

// startWorkers starts parallel workers working on inputCh from fuzzer context.
func startWorkers(parallel int, worker func(context.Context)) {
	var wg sync.WaitGroup
	for i := 1; i <= parallel; i++ {
		wg.Add(1)
		ctx := newWorkerContext(strconv.Itoa(i))

		// Start worker
		go func() {
			logger := getWorkerLogger(ctx)
			logger.Printf("[Worker %d] started", i)
			defer wg.Done()
			worker(ctx)
		}()
	}

	wg.Wait()
}

func queueEntryWorker(ctx context.Context, queueEntryCh chan fuzz.QueueEntry) {
	logger := getWorkerLogger(ctx)

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
				log.Printf("[Worker %d][Task %s] Error: %s\n", i, task.id, err)
				continue
			}
			err = HandleRunResult(ctx, task, result, fuzzCtx)
			if err != nil {
				log.Printf("[Worker %d][Task %s] Error: %s\n", i, task.id, err)
				continue
			}
		case <-time.After(3 * time.Minute):
			log.Printf("[Worker %d] Timeout. Exited", i)
			return
		}
	}
}

const CTX_KEY_WORKER_ID = "WORKER_ID"

func newWorkerContext(workerID string) context.Context {
	return context.WithValue(context.Background(), CTX_KEY_WORKER_ID, workerID)
}

func getWorkerID(context context.Context) string {
	return context.Value(CTX_KEY_WORKER_ID).(string)
}

func getWorkerLogger(context context.Context) *log.Logger {
	workerID := getWorkerID(context)
	return gLog.NewLogger(fmt.Sprintf("[worker %s]", workerID))
}
