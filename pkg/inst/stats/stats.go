package stats

import "gfuzz/pkg/stats"

type GoOp string

const (
	ChSend  GoOp = "chsend"
	ChMake  GoOp = "chmake"
	ChRecv  GoOp = "chrecv"
	ChClose GoOp = "chclose"
)

var (
	codeStats = stats.NewStats()
)

func IncGoOp(op GoOp) {
	codeStats.Inc(string(op))
}

func SetSelectNumOfCases(selectID string, numOfCases uint64) {
	codeStats.Set(selectID, numOfCases)
}

func ToFile(dstFile string) error {
	return stats.ToFile(codeStats, dstFile)
}
