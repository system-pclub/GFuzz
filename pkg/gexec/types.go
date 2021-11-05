package gexec

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path"
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

func (g *GoBinTest) GetCmd(ctx context.Context) (*exec.Cmd, error) {
	return exec.CommandContext(ctx, g.Bin, "-test.timeout", "1m", "-test.count=1", "-test.parallel=1", "-test.v", "-test.run", g.Func), nil
}

func (g *GoBinTest) String() string {
	return fmt.Sprintf("%s-%s", g.Bin, g.Func)
}

type GoPkgTest struct {
	// Test function name
	Func string

	// Which package the test function located (if bin has set, package is bin file name since test bin is compiled by package level)
	Package string

	// Where is folder contains the go.mod
	GoModDir string
}

func (g *GoPkgTest) GetCmd(ctx context.Context) (*exec.Cmd, error) {
	var pkg = g.Package
	if pkg == "" {
		pkg = "./..."
	}
	cmd := exec.CommandContext(ctx, "go", "test", "-timeout=1m", "-count=1", "-v", "-run", g.Func, pkg)
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
