package fuzzer

import (
	"context"
	"errors"
	"fmt"
	fexec "gfuzz/pkg/fuzz/exec"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

func Run(ctx context.Context, input *fexec.Input) (*fexec.Output, error) {

	logger := getWorkerLogger(ctx)

	var err error

	// Setting up related file paths
	err = os.MkdirAll(input.OutputDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	gfInputFp, err := input.GetInputFilePath()
	if err != nil {
		return nil, err
	}

	gfOutputFp, err := input.GetOutputFilePath()
	if err != nil {
		return nil, err
	}
	gfRtOutputFp, err := input.GetOracleRtOutputFilePath()
	if err != nil {
		return nil, err
	}

	// Create the input file into disk
	iBytes, err := fexec.Serialize(input)
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(gfInputFp, iBytes, os.ModePerm)
	if err != nil {
		return nil, err
	}

	// prepare timeout context
	cmdCtx, cancel := context.WithTimeout(ctx, time.Duration(3)*time.Minute)
	defer cancel()

	cmd, err := input.Exec.GetCmd(cmdCtx)
	if err != nil {
		return nil, err
	}

	// setting up environment variables
	env := os.Environ()
	env = append(env, fmt.Sprintf("GF_RECORD_FILE=%s", gfRecordFp))
	env = append(env, fmt.Sprintf("GF_OUTPUT_FILE=%s", gfOutputFp))
	env = append(env, fmt.Sprintf("GF_INPUT_FILE=%s", gfInputFp))

	if GoRoot != "" {
		env = append(env, fmt.Sprintf("GOROOT=%s", GoRoot))
	}

	cmd.Env = env

	// Save output to the file
	outF, err := os.Create(gfOutputFp)
	if err != nil {
		return nil, fmt.Errorf("create stdout: %s", err)
	}
	defer outF.Close()

	//var buf bytes.Buffer
	cmd.Stdout = outF
	cmd.Stderr = outF

	runErr := cmd.Run()

	var timeout bool = false
	if runErr != nil {
		// Go test failed might be intentional
		log.Printf("[Worker %s][Task %s][ignored] go test command failed: %v", workerID, task.id, runErr)

		if errors.Is(runErr, context.DeadlineExceeded) {
			timeout = true
		}
	}

	// Read the printed record file
	var record *Record
	b, err := ioutil.ReadFile(gfRecordFp)
	if err != nil {
		log.Printf("[Worker %s][Task %s][ignored] cannot read record file %s: %v", workerID, task.id, gfRecordFp, err)
	} else {
		record, err = ParseRecordFile(string(b))

		if err != nil {
			log.Printf("[Worker %s][Task %s][ignored] failed to parse record file %s: %v", workerID, task.id, gfRecordFp, err)
		}
	}

	// Read bug IDs from stdout
	var bugIDs []string
	b, err = ioutil.ReadFile(gfOutputFp)
	if err != nil {
		// if error happened, log and ignore
		log.Printf("[Task %s][ignored] cannot read output file %s: %v", task.id, gfOutputFp, err)
	} else {
		output := string(b)
		bugIDs, err = GetListOfBugIDFromStdoutContent(output)
		if err != nil {
			log.Printf("[Task %s][ignored] failed to find bug from output %s: %v", task.id, gfOutputFp, err)
		}

		if strings.Contains(output, "panic: test timed out after") {
			timeout = true
		}
	}

	// Read executed operations from gfOpCovFp
	b, err = ioutil.ReadFile(gfOpCovFp)
	if err != nil {
		// if error happened, log and ignore
		log.Printf("[Task %s][ignored] cannot read operation coverage file %s: %v", task.id, gfOpCovFp, err)
	}
	ids := strings.Split(string(b), "\n")

	retOutput := &RunResult{
		RetInput:       retInput,
		RetRecord:      record,
		BugIDs:         bugIDs,
		StdoutFilepath: gfOutputFp,
		opIDs:          ids,
		IsTimeout:      timeout,
	}

	// Increment number of runs after a successfully run
	fuzzCtx.IncNumOfRun()
	return retOutput, nil
}
