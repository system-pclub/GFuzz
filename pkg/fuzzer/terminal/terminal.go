package terminal

import (
	"fmt"
	"gfuzz/pkg/fuzz/api"
	"os"
	"text/tabwriter"
	"time"
)

type TerminalReport struct {
	Bugs     int
	Duration time.Duration
	Runs     uint64
}

func (t *TerminalReport) GetTerminalRows() []string {
	return []string{
		fmt.Sprintf("# of Runs:\t%d", t.Runs),
		fmt.Sprintf("# of Bugs:\t%d", t.Bugs),
		fmt.Sprintf("Duration:\t%f min(s)", t.Duration.Minutes()),
	}
}

func Render(ch chan *TerminalReport) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)

	for report := range ch {
		fmt.Print("\033[H\033[2J")
		for _, row := range report.GetTerminalRows() {
			fmt.Fprintln(w, row)
		}
		w.Flush()

	}
}

func Feed(ch chan *TerminalReport, fctx *api.Context) {
	for {
		report := &TerminalReport{
			Bugs:     fctx.GetNumOfBugs(),
			Duration: fctx.GetDuration(),
			Runs:     fctx.GetNumOfRuns(),
		}
		ch <- report
		time.Sleep(1 * time.Second)
	}
}
