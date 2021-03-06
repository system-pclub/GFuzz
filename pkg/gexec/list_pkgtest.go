package gexec

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	oe "os/exec"
	"path"
	"strings"
	"time"
)

// ListPackages lists all packages in the current module
// (Has to be run at the directory contains go.mod)
func ListPackages(goModRootPath string) ([]string, error) {
	cmd := oe.Command("go", "list", "./...")
	if goModRootPath != "" {
		cmd.Dir = goModRootPath
	}
	cmd.Env = os.Environ()

	var out bytes.Buffer
	w := io.MultiWriter(&out, log.Writer())
	cmd.Stdout = w
	cmd.Stderr = w

	log.Printf("go list ./... in %s", goModRootPath)
	err := cmd.Run()

	if err != nil {
		log.Printf("[go list ./...]: %v", err)
	}

	return parseGoCmdListOutput(out.String())

}

func parseGoCmdListOutput(output string) ([]string, error) {
	lines := strings.Split(output, "\n")
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.HasPrefix(line, "go: downloading") {
			continue
		}
		if line != "" {
			filtered = append(filtered, line)
		}
	}
	return filtered, nil
}

// ListTestsInPackage lists all tests in the given package
// pkg can be ./... to search in all packages
func ListExecutablesInPackage(goModDir string, pkg string) ([]Executable, error) {
	if pkg == "" {
		pkg = "./..."
	}

	// prepare timeout context
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Minute)
	defer cancel()

	cmd := oe.CommandContext(ctx, "go", "test", "-list", ".*", pkg)
	if goModDir != "" {
		cmd.Dir = goModDir
	}
	cmd.Env = os.Environ()

	var out bytes.Buffer
	w := io.MultiWriter(&out, log.Writer())
	cmd.Stdout = w
	cmd.Stderr = w

	log.Printf("go test -list .* %s", pkg)

	err := cmd.Run()

	if err != nil {
		return nil, fmt.Errorf("[go test -list .* %s] %v", pkg, err)
	}

	testFuncs, err := parseGoCmdTestListOutput(out.String())
	if err != nil {
		return nil, err
	}

	goTests := make([]Executable, 0, len(testFuncs))
	for _, testFunc := range testFuncs {
		goTests = append(goTests, &GoPkgTest{
			Func:     testFunc,
			Package:  pkg,
			GoModDir: goModDir,
		})
	}
	return goTests, nil
}

func parseGoCmdTestListOutput(output string) ([]string, error) {
	lines := strings.Split(output, "\n")
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		// To filter out output likes
		// ?   	k8s.io/kubernetes/cluster/images/etcd-version-monitor	[no test files]
		// ok      goFuzz/example/simple1  0.218s
		// Only keep output like:
		// TestParseInputFileHappy

		if line != "" && strings.HasPrefix(line, "Test") && !strings.ContainsAny(line, " ") {
			filtered = append(filtered, line)
		}
	}
	return filtered, nil
}

// ListExecutablesFromGoModule will return all executable tests
// under given go module directory
// If forceBinTest is used, then tests will first be compiled to binary file
// and then use binary file to list tests (from GoPkgTest to GoBinTest)
func ListExecutablesFromGoModule(goModDir string,
	pkgs []string, forceBinTest bool, outputDir string) ([]Executable, error) {
	var err error
	if pkgs == nil {
		// Find all tests in all packages
		pkgs, err = ListPackages(goModDir)
		if err != nil {
			return nil, fmt.Errorf("failed to list packages at %s: %v", goModDir, err)
		}
	}

	var execs []Executable

	// ListTestsInPackage utilized command `go test -list` which cannot be run in parallel if they share same go code file.
	// Run parallel will cause `intput/output error` when `go test` tries to open file already opened by previous `go test` command.
	// Using other methold like `find Test | grep` can find test name but cannot find package location
	for _, pkg := range pkgs {
		var testsInPkg []Executable
		if forceBinTest {
			testBinFile := GetTestBinFileName(pkg)
			fullTestBinFile := path.Join(outputDir, testBinFile)
			err = CompileTestBinary(goModDir, pkg, fullTestBinFile)
			if err != nil {
				log.Printf("[ignored] failed to compile package %s: %v", pkg, err)
				continue
			}
			testsInPkg, err = ListExecutablesFromTestBin(fullTestBinFile)
		} else {
			testsInPkg, err = ListExecutablesInPackage(goModDir, pkg)
		}

		if err != nil {
			log.Printf("[ignored] failed to list tests at package %s: %v", pkg, err)
			continue
		}
		execs = append(execs, testsInPkg...)

	}

	return execs, nil
}

func GetTestBinFileName(pkg string) string {
	pkgName := strings.ReplaceAll(pkg, "/", "-")
	return pkgName + ".test"
}

// CompileTestBinary compiles the tests under the given package to the given
// output file.
func CompileTestBinary(goModDir string, pkg string, output string) error {

	cmdParams := []string{
		"go", "test", "-c", "-o", output, pkg,
	}
	// go test -c -o <output> <pkg>
	cmd := oe.Command(cmdParams[0], cmdParams[1:]...)
	if goModDir != "" {
		cmd.Dir = goModDir
	}
	cmd.Env = os.Environ()

	var out bytes.Buffer
	w := io.MultiWriter(&out, log.Writer())
	cmd.Stdout = w
	cmd.Stderr = w

	log.Printf("%s", cmdParams)

	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}
