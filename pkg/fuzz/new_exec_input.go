package fuzz

import (
	"fmt"
	"gfuzz/pkg/fuzz/exec"
	"gfuzz/pkg/fuzz/gexecfuzz"
	"gfuzz/pkg/oraclert/config"
	"path"
)

// newExecInput should be the only way to create exec.Input
func newExecInput(fc *Context, g *gexecfuzz.GExecFuzz,
	rtConfig *config.Config, stage exec.Stage) *exec.Input {
	globalID := fc.GetAutoIncGlobalID()
	inputID := fmt.Sprintf("%d-%s-%s", globalID, stage, g.String())
	dir := path.Join(fc.cfg.OutputDir, "exec", inputID)
	return &exec.Input{
		ID:             inputID,
		Exec:           g.Exec,
		OracleRtConfig: rtConfig,
		OutputDir:      dir,
		Stage:          stage,
	}
}

func NewInitExecInput(fc *Context, g *gexecfuzz.GExecFuzz) *exec.Input {
	ortCfg := config.NewConfig()
	return newExecInput(fc, g, ortCfg, exec.InitStage)
}

// HandleFuzzQueryEntry will handle a single queue entry
// func HandleQueueEntry(ctx context.Context, fuzzCtx *Context, e *QueueEntry) ([]*exec.Input, error) {

// 	var fexecs []*exec.Input
// 	if e.Stage == InitStage {
// 		// If stage is InitStage, input's note will be PrintInput and gooracle will record select choices
// 		fe := newExecInput(fuzzCtx, e, e.OracleRtConfig)
// 		fexecs = append(fexecs, fe)
// 	} else if e.Stage == DeterStage {
// 		// If stage is InitStage, input's note will be not PrintInput and expect to have some select choice enforcement
// 		fe := newExecInput(fuzzCtx, e, e.OracleRtConfig)
// 		fexecs = append(fexecs, fe)
// 	} else if e.Stage == CalibStage {
// 		fe := newExecInput(fuzzCtx, e, e.OracleRtConfig)
// 		fexecs = append(fexecs, fe)
// 	} else if e.Stage == RandStage {
// 		randNum := rand.Int31n(101)
// 		if e.BestScore < int(randNum) {
// 			// if skip, simply add entry to the tail
// 			return nil, nil
// 		}
// 		// energy is too large
// 		currentFuzzingEnergy := (e.BestScore / 10) + 1
// 		generatedSelectsHash := make(map[string]bool)
// 		for randFuzzIdx := 0; randFuzzIdx < currentFuzzingEnergy; randFuzzIdx++ {
// 			newRtCgf, err := fuzzCtx.mt.Mutate(e.OracleRtConfig)
// 			if err != nil {
// 				return nil, err
// 			}
// 			selectsHash := hash.AsSha256(newRtCgf)
// 			if _, exist := generatedSelectsHash[selectsHash]; exist {
// 				continue
// 			}
// 			generatedSelectsHash[selectsHash] = true
// 			fe := newExecInput(fuzzCtx, e, newRtCgf)
// 			fexecs = append(fexecs, fe)
// 		}
// 	} else {
// 		return nil, fmt.Errorf("incorrect stage %s", e.Stage)
// 	}

// 	return fexecs, nil

// }
