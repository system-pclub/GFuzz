package gexec

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path"
	"strings"
)

// Executable represents a target that can be triggered by the fuzzer.
// It usually is the program after instrumentation (so that it has oracle at runtime)
type Executable interface {
	// Return a cmd that can be executed
	GetCmd(ctx context.Context) (*exec.Cmd, error)

	// Return an ID/string representation of this executable
	String() string
}

// GoBinTest represents a test from go test binary file
type GoBinTest struct {
	// Test function name
	Func string

	// test should be triggered from compiled binary file
	Bin string
}

// strictGoTestRun will insert ^ at the front and $ at the end if they not present
func strictGoTestRun(run string) string {
	var final string
	if !strings.HasPrefix(run, "^") {
		final = "^" + run
	} else {
		final = run
	}

	if !strings.HasSuffix(final, "$") {
		final = final + "$"
	}
	return final
}

func (g *GoBinTest) GetCmd(ctx context.Context) (*exec.Cmd, error) {
	final := strictGoTestRun(g.Func)
	return exec.CommandContext(ctx, g.Bin, "-test.timeout", "30s", "-test.count=1", "-test.parallel=1", "-test.v", "-test.run", final), nil
}

func (g *GoBinTest) String() string {
	filename := path.Base(g.Bin)
	return fmt.Sprintf("%s-%s", filename, g.Func)
}

type GoPkgTest struct {
	// Test function name
	Func string

	// Which package the test function located
	Package string

	// Where is folder contains the go.mod
	GoModDir string
}

func (g *GoPkgTest) GetCmd(ctx context.Context) (*exec.Cmd, error) {
	var pkg = g.Package
	if pkg == "" {
		pkg = "./..."
	}
	final := strictGoTestRun(g.Func)

	cmd := exec.CommandContext(ctx, "go", "test", "-timeout=30s", "-count=1", "-v", "-run", final, pkg)
	cmd.Dir = g.GoModDir
	return cmd, nil
}

func (g *GoPkgTest) String() string {
	// abc.com/def => def
	basePkg := path.Base(g.Package)
	return fmt.Sprintf("%s-%s", basePkg, g.Func)
}

type Bin string

func (g *Bin) GetCmd(ctx context.Context) (*exec.Cmd, error) {
	if g == nil {
		return nil, errors.New("get command from empty bin")
	}
	return exec.CommandContext(ctx, string(*g)), nil
}

func (g *Bin) String() string {
	if g == nil {
		return ""
	}
	return string(*g)
}
