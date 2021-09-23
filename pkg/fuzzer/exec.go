package fuzzer

import (
	"context"
	"errors"
	"fmt"
	fexec "gfuzz/pkg/fuzz/exec"
	"gfuzz/pkg/fuzzer/bug"
	ortEnv "gfuzz/pkg/oraclert/env"
	ortOut "gfuzz/pkg/oraclert/output"
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
	env = append(env, fmt.Sprintf("%s=%s", ortEnv.ORACLERT_CONFIG_FILE, gfInputFp))
	env = append(env, fmt.Sprintf("%s=%s", ortEnv.ORACLERT_STDOUT_FILE, gfOutputFp))
	env = append(env, fmt.Sprintf("%s=%s", ortEnv.ORACLERT_OUTPUT_FILE, gfRtOutputFp))

	cmd.Env = env

	// redirect stdout to the file
	outF, err := os.Create(gfOutputFp)
	if err != nil {
		return nil, fmt.Errorf("create stdout: %s", err)
	}
	defer outF.Close()

	cmd.Stdout = outF
	cmd.Stderr = outF

	runErr := cmd.Run()

	var timeout bool = false
	if runErr != nil {
		// Go test failed might be intentional
		logger.Printf("[input %s][ignored] run failed: %v", input.ID, runErr)

		if errors.Is(runErr, context.DeadlineExceeded) {
			timeout = true
		}
	}

	// Read oracle runtime output
	ortOutBytes, err := ioutil.ReadFile(gfRtOutputFp)
	var ortOutput *ortOut.Output
	if err != nil {
		logger.Printf("[input %s][ignored] cannot read file %s: %v", input.ID, gfRtOutputFp, err)

	} else {
		ortOutput, err = ortOut.Deserialize(ortOutBytes)
		if err != nil {
			logger.Printf("[input %s][ignored] failed to parse file %s: %v", input.ID, gfRtOutputFp, err)
		}
	}

	// Read bug IDs from stdout
	var bugIDs []string
	b, err := ioutil.ReadFile(gfOutputFp)
	if err != nil {
		// if error happened, log and ignore
		logger.Printf("[input %s][ignored] cannot read output file %s: %v", input.ID, gfOutputFp, err)
	} else {
		output := string(b)
		bugIDs, err = bug.GetListOfBugIDFromStdoutContent(output)
		if err != nil {
			log.Printf("[input %s][ignored] failed to find bug from output %s: %v", input.ID, gfOutputFp, err)
		}

		if strings.Contains(output, "panic: test timed out after") {
			timeout = true
		}
	}

	execOutput := &fexec.Output{
		BugIDs:         bugIDs,
		IsTimeout:      timeout,
		OracleRtOutput: ortOutput,
	}

	return execOutput, nil
}
