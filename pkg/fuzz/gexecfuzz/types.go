package gexecfuzz

import (
	"gfuzz/pkg/gexec"
	"gfuzz/pkg/oraclert/output"
)

// GExecFuzz all information about fuzzing progress for given executable
type GExecFuzz struct {
	BestScore int
	Exec      gexec.Executable

	// Oracler runtime related
	OrtConfigHashes []string
	// All selects we have seen so far
	OrtSelects map[string]output.SelectRecord
}

func NewGExecFuzz(exec gexec.Executable) *GExecFuzz {
	return &GExecFuzz{
		Exec:      exec,
		BestScore: 0,
	}
}

func (e *GExecFuzz) String() string {
	return e.Exec.String()
}
