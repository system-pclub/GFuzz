package fuzz

import (
	"context"
	"fmt"
	"math/rand"
	"path"

	fexec "gfuzz/pkg/fuzz/exec"
	"gfuzz/pkg/oraclert/config"
	"gfuzz/pkg/utils/hash"
)

func newExecInput(fuzzCtx *Context, e *QueueEntry, rtConfig *config.Config) *fexec.Input {
	globalID := fuzzCtx.GetAutoIncGlobalID()
	inputID := fmt.Sprintf("%s:%s", globalID, e.String())
	dir := path.Join(fuzzCtx.cfg.OutputDir, "exec", inputID)
	return &fexec.Input{
		ID:             inputID,
		Exec:           e.Exec,
		OracleRtConfig: rtConfig,
		OutputDir:      dir,
	}
}

// HandleFuzzQueryEntry will handle a single queue entry
func HandleQueueEntry(ctx context.Context, fuzzCtx *Context, e *QueueEntry) ([]*fexec.Input, error) {
	if skipEntry(fuzzCtx, e) {
		return nil, nil
	}
	var fexecs []*fexec.Input
	if e.Stage == InitStage {
		// If stage is InitStage, input's note will be PrintInput and gooracle will record select choices
		fe := newExecInput(fuzzCtx, e, e.OracleRtConfig)
		fexecs = append(fexecs, fe)
	} else if e.Stage == DeterStage {
		// If stage is InitStage, input's note will be not PrintInput and expect to have some select choice enforcement
		fe := newExecInput(fuzzCtx, e, e.OracleRtConfig)
		fexecs = append(fexecs, fe)
	} else if e.Stage == CalibStage {
		fe := newExecInput(fuzzCtx, e, e.OracleRtConfig)
		fexecs = append(fexecs, fe)
	} else if e.Stage == RandStage {
		randNum := rand.Int31n(101)
		if e.BestScore < int(randNum) {
			// if skip, simply add entry to the tail
			return nil, nil
		}
		// energy is too large
		currentFuzzingEnergy := (e.BestScore / 10) + 1
		generatedSelectsHash := make(map[string]bool)
		for randFuzzIdx := 0; randFuzzIdx < currentFuzzingEnergy; randFuzzIdx++ {
			newRtCgf, err := fuzzCtx.cfg.MutateStrategy.Mutate(e.OracleRtConfig)
			if err != nil {
				return nil, err
			}
			selectsHash := hash.AsSha256(newRtCgf)
			if _, exist := generatedSelectsHash[selectsHash]; exist {
				continue
			}
			generatedSelectsHash[selectsHash] = true
			fe := newExecInput(fuzzCtx, e, newRtCgf)
			fexecs = append(fexecs, fe)
		}
	} else {
		return nil, fmt.Errorf("incorrect stage %s", e.Stage)
	}

	return fexecs, nil

}

// HandleExecOutput handles an fuzz execution output
func HandleExecOutput(fuzzCtx *Context, input *fexec.Input, output *fexec.Output) error {

	return nil
}

// skipEntry return true if given fuzz entry need to be skipped
func skipEntry(fuzzCtx *Context, e *QueueEntry) bool {
	// only check when it is in rand stage
	if e.Stage != RandStage {
		return false
	}

	// TODO
	return false
}
