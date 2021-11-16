package fuzz

import (
	"fmt"
	"gfuzz/pkg/fuzz/exec"
	"gfuzz/pkg/gexec"
	"gfuzz/pkg/oraclert/config"
	"path"
	"strconv"
	"strings"
)

// newExecInput should be the only way to create exec.Input
func NewExecInput(ID uint32, fromID uint32, outputDir string, ge gexec.Executable,
	rtConfig *config.Config, stage exec.Stage) *exec.Input {
	inputID := fmt.Sprintf("%d-%s-%s-%d", ID, stage, ge.String(), fromID)
	dir := path.Join(outputDir, "exec", inputID)
	return &exec.Input{
		ID:             inputID,
		Exec:           ge,
		OracleRtConfig: rtConfig,
		OutputDir:      dir,
		Stage:          stage,
	}
}

func GetExecIDFromInputID(inputID string) (uint32, error) {
	id, err := strconv.Atoi(strings.Split(inputID, "-")[0])
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}

func NewInitExecInput(fc *Context, ge gexec.Executable) *exec.Input {
	ortCfg := config.NewConfig()
	globalID := fc.GetAutoIncGlobalID()
	return NewExecInput(globalID, 0, fc.cfg.OutputDir, ge, ortCfg, exec.InitStage)
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
