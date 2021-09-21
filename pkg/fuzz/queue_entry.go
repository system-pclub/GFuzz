package fuzz

import (
	"fmt"
	"log"
	"math/rand"

	gExec "gfuzz/pkg/exec"
	"gfuzz/pkg/oraclert"
)

// QueueEntry records all information about fuzzing progress for given executable
type QueueEntry struct {
	Stage                FuzzStage
	BestScore            int
	ExecutionCount       int
	OracleRtConfig       *oraclert.Config
	Exec                 gExec.Executable
	OracleRtConfigHashes []string
}

func (e *QueueEntry) String() string {
	return fmt.Sprintf("%s:%s:%d", e.Exec, e.Stage, e.Idx)
}

// HandleFuzzQueryEntry will handle a single entry from fuzzCtx's fuzzingQueue
// Notes:
//   1. e is expected to be dequeue from fuzzCtx's fuzzingQueue
func HandleFuzzQueryEntry(e *QueueEntry, fuzzCtx *FuzzContext) error {
	// TODO: better way to print FuzzQueryEntry, maybe ID or string of input?
	log.Printf("handle entry: %s\n", e)

	if shouldDropFuzzQueryEntry(fuzzCtx, e) {
		return nil
	}

	var runTasks []*RunTask

	if e.Stage == InitStage {
		// If stage is InitStage, input's note will be PrintInput and gooracle will record select choices
		t, err := NewRunTask(e.CurrInput, e.Stage, e.Idx, e.ExecutionCount, e)
		if err != nil {
			return err
		}
		runTasks = append(runTasks, t)
	} else if e.Stage == DeterStage {
		// If stage is InitStage, input's note will be not PrintInput and expect to have some select choice enforcement
		t, err := NewRunTask(e.CurrInput, e.Stage, e.Idx, e.ExecutionCount, e)
		if err != nil {
			return err
		}
		runTasks = append(runTasks, t)
	} else if e.Stage == CalibStage {
		t, err := NewRunTask(e.CurrInput, e.Stage, e.Idx, e.ExecutionCount, e)
		if err != nil {
			return err
		}
		runTasks = append(runTasks, t)
	} else if e.Stage == RandStage {

		randNum := rand.Int31n(101)
		if e.BestScore < int(randNum) {
			log.Printf("[%s] randomly skipped", e)
			// if skip, simply add entry to the tail
			fuzzCtx.EnqueueQueryEntry(e)
			return nil
		}
		// energy is too large
		currentFuzzingEnergy := (e.BestScore / 10) + 1
		generatedSelectsHash := make(map[string]bool)
		execCount := e.ExecutionCount
		log.Printf("[%+v] randomly mutate with energy %d", *e, currentFuzzingEnergy)
		for randFuzzIdx := 0; randFuzzIdx < currentFuzzingEnergy; randFuzzIdx++ {
			randomInput, err := RandomMutateInput(e.CurrInput)
			if err != nil {
				log.Printf("[%s] randomly mutate input fail: %s, continue", e, err)
				continue
			}
			selectsHash := GetHashOfSelects(randomInput.VecSelect)
			if _, exist := generatedSelectsHash[selectsHash]; exist {
				log.Printf("[%s][%d] skip generated input because of duplication", e, randFuzzIdx)
				continue
			}
			generatedSelectsHash[selectsHash] = true
			log.Printf("[%s][%d] successfully generate input", e, randFuzzIdx)
			t, err := NewRunTask(randomInput, e.Stage, e.Idx, execCount, e)
			if err != nil {
				return err
			}
			runTasks = append(runTasks, t)
			execCount += 1
		}
		e.ExecutionCount = execCount
		fuzzCtx.EnqueueQueryEntry(e)
	} else {
		return fmt.Errorf("incorrect stage found: %s", e.Stage)
	}

	for _, t := range runTasks {
		fuzzCtx.runTaskCh <- t
	}

	return nil

}

// shouldDropFuzzQueryEntry return true if given fuzz entry need to be dropped
func shouldDropFuzzQueryEntry(fuzzCtx *FuzzContext, e *QueueEntry) bool {
	// only check when it is in rand stage
	if e.Stage != RandStage {
		return false
	}
	return ShouldSkipInput(fuzzCtx, e.Exec.String(), e.CurrInput)
}
