package fuzzer

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	gExec "gfuzz/pkg/exec"
	"gfuzz/pkg/oraclert"
)

type ExecTask struct {
	ID             string
	OracleRtConfig *oraclert.Config
	Exec           gExec.Executable
}

type ExecOutput struct {
	OracleRtOutput *oraclert.Output
}

func getInputFilePath(outputDir string) (string, error) {
	return filepath.Abs(path.Join(outputDir, "input"))
}

func getOutputFilePath(outputDir string) (string, error) {
	return filepath.Abs(path.Join(outputDir, "stdout"))
}

func getOpCovFilePath(outputDir string) (string, error) {
	return filepath.Abs(path.Join(outputDir, "opcov"))
}

func getRecordFilePath(outputDir string) (string, error) {
	return filepath.Abs(path.Join(outputDir, "record"))
}

func Run(ctx context.Context, fuzzCtx *FuzzContext, task *RunTask) (*RunResult, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[Task %s] recovered from panic in fuzzer", task.id)
		}
	}()

	logger := getWorkerLogger(ctx)

	var err error
	input := task.input

	// Setting up related file paths

	runOutputDir := path.Join(OutputDir, task.id)
	err = createDir(runOutputDir)
	if err != nil {
		return nil, err
	}

	gfInputFp, err := getInputFilePath(runOutputDir)
	if err != nil {
		return nil, err
	}

	gfOutputFp, err := getOutputFilePath(runOutputDir)
	if err != nil {
		return nil, err
	}
	gfRecordFp, err := getRecordFilePath(runOutputDir)
	if err != nil {
		return nil, err
	}
	// gfErrFp, err := getErrFilePath(runOutputDir)
	// if err != nil {
	// 	return nil, err
	// }
	gfOpCovFp, err := getOpCovFilePath(runOutputDir)
	if err != nil {
		return nil, err
	}

	boolFirstRun := input.Note == NotePrintInput

	// Create the input file into disk
	err = SerializeInput(input, gfInputFp)
	if err != nil {
		return nil, err
	}

	var globalTuple string
	if GlobalTuple {
		globalTuple = "1"
	} else {
		globalTuple = "0"
	}

	// prepare timeout context
	runCtx, cancel := context.WithTimeout(ctx, time.Duration(3)*time.Minute)
	defer cancel()

	var cmd *exec.Cmd
	if task.input.GoTestCmd != nil {
		if task.input.GoTestCmd.Bin != "" {
			// Since golang's compiled test can only be one per package, so we just assume the test func must exist in the given binary
			cmd = exec.CommandContext(runCtx, task.input.GoTestCmd.Bin, "-test.timeout", "1m", "-test.parallel", "1", "-test.v", "-test.run", input.GoTestCmd.Func)
		} else {
			var pkg = input.GoTestCmd.Package
			if pkg == "" {
				pkg = "./..."
			}
			cmd = exec.CommandContext(runCtx, "go", "test", "-timeout", "1m", "-v", "-run", input.GoTestCmd.Func, pkg)
		}
	} else if task.input.CustomCmd != "" {
		cmds := strings.SplitN(task.input.CustomCmd, " ", 2)
		cmd = exec.CommandContext(runCtx, cmds[0], cmds[1])
	} else {
		return nil, fmt.Errorf("either testname or custom command is required")
	}
	cmd.Dir = TargetGoModDir

	// setting up environment variables
	env := os.Environ()
	env = append(env, fmt.Sprintf("GF_RECORD_FILE=%s", gfRecordFp))
	env = append(env, fmt.Sprintf("GF_OUTPUT_FILE=%s", gfOutputFp))
	env = append(env, fmt.Sprintf("GF_INPUT_FILE=%s", gfInputFp))
	env = append(env, fmt.Sprintf("GF_OP_COV_FILE=%s", gfOpCovFp))
	env = append(env, fmt.Sprintf("BitGlobalTuple=%s", globalTuple))
	env = append(env, fmt.Sprintf("GF_TIME_DIVIDE=%d", TimeDivide))
	if ScoreSdk {
		env = append(env, "GF_SCORE_SDK=1")
	}
	if ScoreAllPrim {
		env = append(env, "GF_SCORE_TRAD=1")
	}
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

	// Read the newly printed input file if this is the first run
	var retInput *Input
	if boolFirstRun {
		log.Printf("[Worker %s][Task %s] first run, reading input file %s", workerID, task.id, gfInputFp)
		bytes, err := ioutil.ReadFile(gfInputFp)
		if err != nil {
			return nil, err
		}

		retInput, err = ParseInputFile(string(bytes))
		if err != nil {
			return nil, err
		}

		// assign missing parts in input file
		retInput.GoTestCmd = task.input.GoTestCmd
		retInput.CustomCmd = task.input.CustomCmd

	} else {
		retInput = nil
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
