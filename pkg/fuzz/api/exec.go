package api

import (
	"encoding/json"
	"fmt"
	"gfuzz/pkg/gexec"
	"gfuzz/pkg/oraclert/config"
	"gfuzz/pkg/oraclert/output"
	"path"
	"path/filepath"
)

// Stage indicates how we treat/response to an input and corresponding output
type Stage string

const (
	// InitStage simply run the empty without any mutation
	InitStage Stage = "init"

	// DeterStage is to create input by tweak select choice one by one
	DeterStage Stage = "deter"

	// CalibStage choose an input from queue to run (prepare for rand)
	CalibStage Stage = "calib"

	// RandStage randomly mutate select choice
	RandStage Stage = "rand"

	// Run with custom/pre-prepared oracle runtime configuration
	ReplayStage Stage = "replay"
)

// Input contains all information about a single execution
// (usually by fuzzer)
type Input struct {
	// ID is the unique identifer for this execution.
	ID string
	// OracleRtConfig is the configuration for the oracle runtime.
	OracleRtConfig *config.Config
	// Exec is the command to trigger a program with oracle runtime.
	Exec gexec.Executable
	// OutputDir is the output directory for this execution
	OutputDir string

	Stage Stage
}

// Output contains all useful information after a single execution
type Output struct {
	OracleRtOutput *output.Output
	BugIDs         []string
	Timeout        bool
}

func (i *Input) GetOrtConfigFilePath() (string, error) {
	return filepath.Abs(path.Join(i.OutputDir, "ort_config"))
}

func (i *Input) GetOutputFilePath() (string, error) {
	return filepath.Abs(path.Join(i.OutputDir, "stdout"))
}

func (i *Input) GetOrtOutputFilePath() (string, error) {
	return filepath.Abs(path.Join(i.OutputDir, "ort_output"))
}

func Serialize(l *Input) ([]byte, error) {
	if l == nil {
		return []byte{}, nil
	}

	return json.Marshal(l)
}

func Deserilize(data []byte) (*Input, error) {
	l := Input{}
	err := json.Unmarshal(data, &l)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

// newExecInput should be the only way to create exec.Input
func NewExecInput(ID uint32, fromID uint32, outputDir string, ge gexec.Executable,
	rtConfig *config.Config, stage Stage) *Input {
	inputID := fmt.Sprintf("%d-%s-%s-%d", ID, stage, ge.String(), fromID)
	dir := path.Join(outputDir, "exec", inputID)
	return &Input{
		ID:             inputID,
		Exec:           ge,
		OracleRtConfig: rtConfig,
		OutputDir:      dir,
		Stage:          stage,
	}
}

func NewInitExecInput(fctx *Context, ge gexec.Executable) *Input {
	ortCfg := config.NewConfig()
	ortCfg.SelEfcm.SelTimeout = fctx.Cfg.SelEfcmTimeout
	globalID := fctx.GetAutoIncGlobalID()
	return NewExecInput(globalID, 0, fctx.Cfg.OutputDir, ge, ortCfg, InitStage)
}
