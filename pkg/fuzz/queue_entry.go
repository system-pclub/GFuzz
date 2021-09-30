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
	OracleRtConfig       *config.Config
	Exec                 gexec.Executable
	OracleRtConfigHashes []string
}

func (e *QueueEntry) String() string {
	return fmt.Sprintf("%s:%s", e.Exec, e.Stage)
}
