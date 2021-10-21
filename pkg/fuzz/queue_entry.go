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

func newQueueEntry(exec gexec.Executable) *QueueEntry {
	return &QueueEntry{
		Exec:      exec,
		Stage:     InitStage,
		BestScore: 0,
		OracleRtConfig: &config.Config{
			RecordSelect: true,
		},
	}
}

func (e *QueueEntry) String() string {
	return fmt.Sprintf("%s-%s", e.Exec, e.Stage)
}
