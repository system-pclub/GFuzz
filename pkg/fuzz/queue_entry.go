package fuzz

import (
	"fmt"

	gexec "gfuzz/pkg/exec"
	"gfuzz/pkg/oraclert/config"
)

// QueueEntry records all information about fuzzing progress for given executable
type QueueEntry struct {
	Stage                FuzzStage
	BestScore            int
	ExecutionCount       int
	OracleRtConfig       *config.Config
	Exec                 gexec.Executable
	OracleRtConfigHashes []string
}

func (e *QueueEntry) String() string {
	return fmt.Sprintf("%s:%s:%d", e.Exec, e.Stage, e.ExecutionCount)
}
