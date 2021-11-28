package gexec

import (
	"bytes"
	"gfuzz/pkg/utils/fs"
	"io"
	"log"
	"os/exec"
)

func ListExecutablesFromTestBinGlobs(globs []string) ([]Executable, error) {
	var tests []Executable
	for _, glob := range globs {
		ts, err := ListExecutablesFromTestBinGlob(glob)
		if err != nil {
			log.Printf("ListExecutablesFromTestBinGlob '%s' failed: %v", glob, err)
		} else {
			tests = append(tests, ts...)
		}
	}
	return tests, nil
}

func ListExecutablesFromTestBinGlob(glob string) ([]Executable, error) {
	files, err := fs.ListFilesByGlob(glob)
	if err != nil {
		return nil, err
	}
	var tests []Executable
	for _, file := range files {
		testsInFile, err := ListExecutablesFromTestBin(file)
		if err != nil {
			log.Printf("ListExecutablesFromTestBin '%s' failed: %v", file, err)
		} else {
			tests = append(tests, testsInFile...)
		}
	}
	return tests, nil
}

func ListExecutablesFromTestBin(testBin string) ([]Executable, error) {
	cmd := exec.Command(testBin, "-test.list", ".*")
	var out bytes.Buffer
	w := io.MultiWriter(&out, log.Writer())
	cmd.Stdout = w
	cmd.Stderr = w

	log.Printf("%s -test.list .*", testBin)

	err := cmd.Run()

	if err != nil {
		log.Printf("[%s -test.list .*] failed: %v", testBin, err)
	}

	testFuncs, err := parseGoCmdTestListOutput(out.String())
	if err != nil {
		return nil, err
	}

	goTests := make([]Executable, 0, len(testFuncs))
	for _, testFunc := range testFuncs {
		goTests = append(goTests, &GoBinTest{
			Func: testFunc,
			Bin:  testBin,
		})
	}
	return goTests, nil

}
