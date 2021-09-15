package gotest

import (
	"bytes"
	"fmt"
	"gfuzz/pkg/utils/fs"
	"io"
	"log"
	"os/exec"
	"path"
	"strings"
)

func ListGoTestsFromTestBinGlobs(globs []string) ([]*GoTest, error) {
	var tests []*GoTest
	for _, glob := range globs {
		ts, err := ListGoTestsFromTestBinGlob(glob)
		if err != nil {
			log.Printf("ListGoTestsFromTestBinGlob '%s' failed: %v", glob, err)
		} else {
			tests = append(tests, ts...)
		}
	}
	return tests, nil
}

func ListGoTestsFromTestBinGlob(glob string) ([]*GoTest, error) {
	files, err := fs.ListFilesByGlob(glob)
	if err != nil {
		return nil, err
	}
	var tests []*GoTest
	for _, file := range files {
		testsInFile, err := ListGoTestsFromTestBin(file)
		if err != nil {
			log.Printf("ListGoTestsFromTestBin '%s' failed: %v", file, err)
		} else {
			tests = append(tests, testsInFile...)
		}
	}
	return tests, nil
}

func ListGoTestsFromTestBin(testBin string) ([]*GoTest, error) {
	cmd := exec.Command(testBin, "-test.list", ".*")
	var out bytes.Buffer
	w := io.MultiWriter(&out, log.Writer())
	cmd.Stdout = w
	cmd.Stderr = w

	log.Printf("%s -test.list .*", testBin)

	err := cmd.Run()

	if err != nil {
		return nil, fmt.Errorf("[%s -test.list .*] failed: %v", testBin, err)
	}

	testFuncs, err := parseGoCmdTestListOutput(out.String())
	if err != nil {
		return nil, err
	}

	goTests := make([]*GoTest, 0, len(testFuncs))
	binName := path.Base(testBin)
	for _, testFunc := range testFuncs {
		goTests = append(goTests, &GoTest{
			Func:    testFunc,
			Bin:     testBin,
			Package: binName,
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

		if line != "" && strings.HasPrefix(line, "Test") && line != "Test" && !strings.ContainsAny(line, " ") {
			filtered = append(filtered, line)
		}
	}
	return filtered, nil
}
