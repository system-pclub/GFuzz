package interest

import (
	"fmt"
	"gfuzz/pkg/fuzz/api"
	"gfuzz/pkg/fuzz/mutate"
	"gfuzz/pkg/utils/hash"
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
	ortCfgHash := hash.AsSha256(o.OracleRtOutput)
	if h.fctx.UpdateOrtOutputHash(ortCfgHash) {
		return true, nil
	}
	return false, nil
}
func (h *InterestHandlerImpl) HandleInterest(i *api.InterestInput) error {
	i.HandledCnt += 1
	// if interested input has not been executed, execute first
	if !i.Executed {
		i.Executed = true
		h.fctx.ExecInputCh <- i.Input
		return nil
	}

	if i.Output == nil {
		// if executed is true but output is nil
		// it could be still in channel pending to run
		return nil
	}

	// if interested input has been executed, then try to mutate and send to execution according to its stage
	switch i.Input.Stage {
	case api.InitStage:
		if i.HandledCnt == 1 {
			// if handle init first time, treat it as normal init stage
			return handleInitStageInput(h.fctx, i.Input, i.Output)

		} else {
			// if not the first time, treat it as random
			return handleRandStageInput(h.fctx, i.Input, i.Output)
		}
		// we are handling the output from the input with init stage
	case api.DeterStage:
		// we are handling the output from the input with deter stage
		return handleDeterStageInput(h.fctx, i.Input, i.Output)
	case api.CalibStage:
		// we are handling the output from the input with calib stage
		return handleCalibStageInput(h.fctx, i.Input, i.Output)
	case api.RandStage:
		// we are handling the output from the input with rand stage
		return handleRandStageInput(h.fctx, i.Input, i.Output)
	case api.ReplayStage:
		// no need to handle replay
		return nil
	default:
		return fmt.Errorf("unexpected stage: %s", i.Input.Stage)
	}

}

func handleInitStageInput(fctx *api.Context, i *api.Input, o *api.Output) error {

	g := fctx.GetQueueEntryByGExecID(i.Exec.String())
	execID, err := getExecIDFromInputID(i.ID)
	if err != nil {
		return err
	}
	var deterInputs []*api.Input
	var mts mutate.OrtConfigMutateStrategy = &mutate.DeterMutateStrategy{}

	cfgs, err := mts.Mutate(g, i.OracleRtConfig, o.OracleRtOutput)
	if err != nil {
		return err
	}

	for _, cfg := range cfgs {
		deterInputs = append(deterInputs, api.NewExecInput(fctx.GetAutoIncGlobalID(), execID, fctx.Cfg.OutputDir, g.Exec, cfg, api.DeterStage))
	}

	for _, input := range deterInputs {
		fctx.ExecInputCh <- input
	}

	return nil
}

func handleDeterStageInput(fctx *api.Context, i *api.Input, o *api.Output) error {
	g := fctx.GetQueueEntryByGExecID(i.Exec.String())
	execID, err := getExecIDFromInputID(i.ID)
	if err != nil {
		return err
	}

	input := api.NewExecInput(fctx.GetAutoIncGlobalID(), execID, fctx.Cfg.OutputDir, g.Exec, i.OracleRtConfig, api.CalibStage)
	fctx.ExecInputCh <- input
	return nil
}

func handleCalibStageInput(fctx *api.Context, i *api.Input, o *api.Output) error {
	g := fctx.GetQueueEntryByGExecID(i.Exec.String())
	execID, err := getExecIDFromInputID(i.ID)
	if err != nil {
		return err
	}

	input := api.NewExecInput(fctx.GetAutoIncGlobalID(), execID, fctx.Cfg.OutputDir, g.Exec, i.OracleRtConfig, api.RandStage)
	fctx.ExecInputCh <- input
	return nil
}

func handleRandStageInput(fctx *api.Context, i *api.Input, o *api.Output) error {
	g := fctx.GetQueueEntryByGExecID(i.Exec.String())
	execID, err := getExecIDFromInputID(i.ID)
	if err != nil {
		return err
	}
	var randInputs []*api.Input
	var mts mutate.OrtConfigMutateStrategy = &mutate.RandomMutateStrategy{}

	cfgs, err := mts.Mutate(g, i.OracleRtConfig, o.OracleRtOutput)
	if err != nil {
		return err
	}

	for _, cfg := range cfgs {
		randInputs = append(randInputs, api.NewExecInput(fctx.GetAutoIncGlobalID(), execID, fctx.Cfg.OutputDir, g.Exec, cfg, api.DeterStage))
	}

	for _, input := range randInputs {
		fctx.ExecInputCh <- input
	}

	return nil
}

func getExecIDFromInputID(inputID string) (uint32, error) {
	id, err := strconv.Atoi(strings.Split(inputID, "-")[0])
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}
