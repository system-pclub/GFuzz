package terminal

import (
	"fmt"
	"gfuzz/pkg/fuzz/api"
	"log"
	"os"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type TerminalReport struct {
	Bugs     int
	Duration time.Duration
	Runs     uint64
}

func (t *TerminalReport) GetTerminalRows() []string {
	return []string{
		fmt.Sprintf("# of Runs: %d", t.Runs),
		fmt.Sprintf("# of Bugs: %d", t.Bugs),
		fmt.Sprintf("Duration: %f minute(s)", t.Duration.Minutes()),
	}
}

func Render(ch chan *TerminalReport) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize terminal UI: %v", err)
	}
	defer ui.Close()
	ui.Clear()

	l := widgets.NewList()
	l.Title = "GFuzz"
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false
	l.SetRect(0, 0, 25, 8)
	uiEvents := ui.PollEvents()

	go func() {
		for {
			e := <-uiEvents
			switch e.ID {
			case "q", "<C-c>":
				ui.Close()
				os.Exit(0)
			}
		}
	}()

	for report := range ch {
		l.Rows = report.GetTerminalRows()
		ui.Render(l)
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
