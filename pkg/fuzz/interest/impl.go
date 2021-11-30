package interest

import (
	"fmt"
	"gfuzz/pkg/fuzz/api"
	"gfuzz/pkg/fuzz/mutate"
	"gfuzz/pkg/utils/hash"
	"gfuzz/pkg/fuzz/score"
	"math"
	"strconv"
	"strings"
)

type InterestHandlerImpl struct {
	fctx *api.Context
}

func NewInterestHandlerImpl(fctx *api.Context) api.InterestHandler {
	return &InterestHandlerImpl{
		fctx: fctx,
	}
}
func (h *InterestHandlerImpl) IsInterested(i *api.Input, o *api.Output) (bool, error) {

	if h.fctx.Cfg.IsIgnoreFeedback == true {
		// if isIgnoreFeedback is true, we treat every feedback as interesting and directly return.
		return true, nil
	}

	oTupleCopy := make(map[uint32]uint32)
	for k, v := range o.OracleRtOutput.Tuples {
		log2element := uint32(math.Log2(float64(v)))
		if log2element > 0 {
			oTupleCopy[k] = log2element
		} else {
			oTupleCopy[k] = 0
		}
	}
	ortCfgHash := hash.AsSha256(oTupleCopy)
	if h.fctx.UpdateOrtOutputHash(ortCfgHash) {
		return true, nil
	}
	return false, nil
}
func (h *InterestHandlerImpl) HandleInterest(i *api.InterestInput) (ret bool, err error) {

	// if init input has not been executed, execute first
	if i.Input.Stage == api.InitStage && !i.Executed {
		i.HandledCnt += 1
		i.Executed = true
		h.fctx.ExecInputCh <- i.Input
		return true, nil
	}

	if i.Output == nil {
		// if executed is true but output is nil
		// it could be still in channel pending to run
		return false, nil
	}

	// if interested input has been executed, then try to mutate and send to execution according to its stage
	switch i.Input.Stage {
	case api.InitStage:
		// we are handling the output from the input with init stage

		//if i.HandledCnt == 1 {
		//	// if handle init first time, treat it as normal init stage
		//	err = handleInitStageInput(h.fctx, i.Input, i.Output)
		//} else {
		//	// if not the first time, treat it as random
		//}

		// Yu:: Directly jump to RandStage if we see the init_stage for the second time.
		// Yu :: The first time of execution should be covered by !i.Executed.
		ret, err = handleRandStageInput(h.fctx, i.Input, i.Output)
	case api.DeterStage:
		// we are handling the output from the input with deter stage
		//err = handleDeterStageInput(h.fctx, i.Input, i.Output)

		// Yu:: Should not be seen in the exec. But if seen, treat it as rand.
		ret, err = handleRandStageInput(h.fctx, i.Input, i.Output)
	case api.CalibStage:
		// we are handling the output from the input with calib stage
		//err = handleCalibStageInput(h.fctx, i.Input, i.Output)

		// Yu:: Should not be seen in the exec. But if seen, treat it as rand.
		ret, err = handleRandStageInput(h.fctx, i.Input, i.Output)
	case api.RandStage:
		// we are handling the output from the input with rand stage
		ret, err = handleRandStageInput(h.fctx, i.Input, i.Output)
	case api.ReplayStage:
		// no need to handle replay

	default:
		err = fmt.Errorf("unexpected stage: %s", i.Input.Stage)
	}

	if err != nil {
		return false, err
	}

	i.HandledCnt += 1

	return ret, nil

}

func handleInitStageInput(fctx *api.Context, i *api.Input, o *api.Output) (bool, error) {

	//g := fctx.GetQueueEntryByGExecID(i.Exec.String())
	//execID, err := getExecIDFromInputID(i.ID)
	//if err != nil {
	//	return false, err
	//}
	//var deterInputs []*api.Input
	//var mts mutate.OrtConfigMutateStrategy = &mutate.DeterMutateStrategy{}
	//
	//if o.OracleRtOutput == nil {
	//	return false, nil
	//}
	//cfgs, err := mts.Mutate(g, i.OracleRtConfig, o.OracleRtOutput, 0)
	//if err != nil {
	//	return false, err
	//}
	//
	//for _, cfg := range cfgs {
	//	deterInputs = append(deterInputs, api.NewExecInput(fctx.GetAutoIncGlobalID(), execID, fctx.Cfg.OutputDir, g.Exec, cfg, api.DeterStage))
	//}
	//
	//if len(deterInputs) == 0 {
	//	return false, nil
	//}
	//for _, input := range deterInputs {
	//	fctx.ExecInputCh <- input
	//}

	return true, nil
}

func handleDeterStageInput(fctx *api.Context, i *api.Input, o *api.Output) (bool, error) {
	g := fctx.GetQueueEntryByGExecID(i.Exec.String())
	execID, err := getExecIDFromInputID(i.ID)
	if err != nil {
		return false, err
	}

	input := api.NewExecInput(fctx.GetAutoIncGlobalID(), execID, fctx.Cfg.OutputDir, g.Exec, i.OracleRtConfig, api.CalibStage)
	fctx.ExecInputCh <- input
	return true, nil
}

func handleCalibStageInput(fctx *api.Context, i *api.Input, o *api.Output) (bool, error) {
	g := fctx.GetQueueEntryByGExecID(i.Exec.String())
	execID, err := getExecIDFromInputID(i.ID)
	if err != nil {
		return false, err
	}

	input := api.NewExecInput(fctx.GetAutoIncGlobalID(), execID, fctx.Cfg.OutputDir, g.Exec, i.OracleRtConfig, api.RandStage)
	fctx.ExecInputCh <- input
	return true, nil
}

func handleRandStageInput(fctx *api.Context, i *api.Input, o *api.Output) (bool, error) {
	g := fctx.GetQueueEntryByGExecID(i.Exec.String())
	execID, err := getExecIDFromInputID(i.ID)
	if err != nil {
		return false, err
	}
	var randInputs []*api.Input
	var mts mutate.OrtConfigMutateStrategy = &mutate.RandomMutateStrategy{}


	randMutateEnergy := fctx.Cfg.RandMutateEnergy

	if fctx.IsUsingScore {
		var scoreFunc api.ScoreStrategy = score.NewScoreStrategyImpl(fctx)
		curScore, _ := scoreFunc.Score(i, o)
		if curScore > fctx.GlobalBestScore {
			fctx.GlobalBestScore = curScore
		}
		randMutateEnergy = int(float64(randMutateEnergy) * (float64(curScore) / float64(fctx.GlobalBestScore)))
	}

	cfgs, err := mts.Mutate(g, i.OracleRtConfig, o.OracleRtOutput, randMutateEnergy)

	if err != nil {
		return false, err
	}

	for _, cfg := range cfgs {
		randInputs = append(randInputs, api.NewExecInput(fctx.GetAutoIncGlobalID(), execID, fctx.Cfg.OutputDir, g.Exec, cfg, api.RandStage))
	}

	for _, input := range randInputs {
		fctx.ExecInputCh <- input
	}

	return true, nil
}

func getExecIDFromInputID(inputID string) (uint32, error) {
	id, err := strconv.Atoi(strings.Split(inputID, "-")[0])
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}
