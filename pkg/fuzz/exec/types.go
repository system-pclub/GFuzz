package exec

import (
	"encoding/json"
	gExec "gfuzz/pkg/exec"
	"gfuzz/pkg/oraclert/config"
	"gfuzz/pkg/oraclert/output"
	"path"
	"path/filepath"
)

// Input contains all information about a single execution
// (usually by fuzzer)
type Input struct {
	// ID is the unique identifer for this execution.
	ID string
	// OracleRtConfig is the configuration for the oracle runtime.
	OracleRtConfig *config.Config
	// Exec is the command to trigger a program with oracle runtime.
	Exec gExec.Executable
	// OutputDir is the output directory
	// for this execution
	OutputDir string
}

// Output contains all useful information after a single execution
type Output struct {
	OracleRtOutput *output.Output
	BugIDs         []string
	IsTimeout      bool
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
