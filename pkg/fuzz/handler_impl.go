package fuzz

import (
	"fmt"
	"gfuzz/pkg/fuzz/exec"
	"gfuzz/pkg/fuzz/mutate"
)

type HandlerImpl struct{}

func (h *HandlerImpl) Handle(fc *Context, i *exec.Input, o *exec.Output) ([]*exec.Input, error) {
	switch i.Stage {
	case exec.InitStage:
		// we are handling the output from the input with init stage
		return handleInitStageInput(fc, i, o)
	case exec.DeterStage:
		// we are handling the output from the input with deter stage
		return handleDeterStageInput(fc, i, o)
	case exec.CalibStage:
		// we are handling the output from the input with calib stage
		return handleCalibStageInput(fc, i, o)
	case exec.RandStage:
		// we are handling the output from the input with rand stage
		return handleRandStageInput(fc, i, o)
	case exec.ReplayStage:
		// no need to handle replay
		return nil, nil
	default:
		return nil, fmt.Errorf("unexpected stage: %s", i.Stage)
	}
}

func handleInitStageInput(fc *Context, i *exec.Input, o *exec.Output) ([]*exec.Input, error) {
	fc.lock.Lock()
	defer fc.lock.Unlock()
	g := fc.GetQueueEntryByGExecID(i.Exec.String())
	execID, err := GetExecIDFromInputID(i.ID)
	if err != nil {
		return nil, err
	}
	var deterInputs []*exec.Input
	var mts mutate.OrtConfigMutateStrategy = &mutate.DeterMutateStrategy{}

	cfgs, err := mts.Mutate(g, i.OracleRtConfig, o.OracleRtOutput)
	if err != nil {
		return nil, err
	}

	for _, cfg := range cfgs {
		deterInputs = append(deterInputs, NewExecInput(fc.GetAutoIncGlobalID(), execID, fc.cfg.OutputDir, g.Exec, cfg, exec.DeterStage))
	}

	return deterInputs, nil
}

func handleDeterStageInput(fc *Context, i *exec.Input, o *exec.Output) ([]*exec.Input, error) {
	fc.lock.Lock()
	defer fc.lock.Unlock()
	g := fc.GetQueueEntryByGExecID(i.Exec.String())
	execID, err := GetExecIDFromInputID(i.ID)
	if err != nil {
		return nil, err
	}
	return []*exec.Input{
		NewExecInput(fc.GetAutoIncGlobalID(), execID, fc.cfg.OutputDir, g.Exec, i.OracleRtConfig, exec.DeterStage),
	}, nil
}

func handleCalibStageInput(fc *Context, i *exec.Input, o *exec.Output) ([]*exec.Input, error) {
	fc.lock.Lock()
	defer fc.lock.Unlock()
	// g := fc.GetQueueEntryByGExecID(i.Exec.String())
	return nil, nil
}

func handleRandStageInput(fc *Context, i *exec.Input, o *exec.Output) ([]*exec.Input, error) {
	fc.lock.Lock()
	defer fc.lock.Unlock()
	// g := fc.GetQueueEntryByGExecID(i.Exec.String())
	return nil, nil
}
