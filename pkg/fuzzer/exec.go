package fuzzer

import (
	"context"
	"errors"
	"fmt"
	"gfuzz/pkg/fuzz/api"
	"gfuzz/pkg/fuzz/config"
	"gfuzz/pkg/fuzzer/bug"
	ortCfg "gfuzz/pkg/oraclert/config"
	ortEnv "gfuzz/pkg/oraclert/env"
	ortOut "gfuzz/pkg/oraclert/output"
	"gfuzz/pkg/utils/hash"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

func Run(ctx context.Context, cfg *config.Config, input *api.Input) (*api.Output, error) {

	logger := getWorkerLogger(ctx)

	var err error

	// Setting up related file paths
	err = os.MkdirAll(input.OutputDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	ortCfgFp, err := input.GetOrtConfigFilePath()
	if err != nil {
		return nil, err
	}

	gfOutputFp, err := input.GetOutputFilePath()
	if err != nil {
		return nil, err
	}
	ortOutputFp, err := input.GetOrtOutputFilePath()
	if err != nil {
		return nil, err
	}

	// Create the input file into disk
	iBytes, err := ortCfg.Serialize(input.OracleRtConfig)
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(ortCfgFp, iBytes, os.ModePerm)
	if err != nil {
		return nil, err
	}

	// prepare timeout context
	// cmd will have 30 seconds timeout, 30 seconds more is for go test to compile/find correct
	// test to run if no package is given
	cmdCtx, cancel := context.WithTimeout(ctx, time.Duration(1)*time.Minute)
	defer cancel()

	cmd, err := input.Exec.GetCmd(cmdCtx)
	if err != nil {
		return nil, err
	}

	// setting up environment variables
	env := os.Environ()
	env = append(env, fmt.Sprintf("%s=%s", ortEnv.ORACLERT_CONFIG_FILE, ortCfgFp))
	env = append(env, fmt.Sprintf("%s=%s", ortEnv.ORACLERT_STDOUT_FILE, gfOutputFp))
	env = append(env, fmt.Sprintf("%s=%s", ortEnv.ORACLERT_OUTPUT_FILE, ortOutputFp))
	if cfg.OracleRtDebug {
		env = append(env, fmt.Sprintf("%s=1", ortEnv.ORACLERT_DEBUG))
	}
	if cfg.TimeDivideBy > 1 {
		env = append(env, fmt.Sprintf("%s=%d", ortEnv.ORACLERT_TIME_DIVIDE, cfg.TimeDivideBy))
	}
	if cfg.NoOracle {
		env = append(env, fmt.Sprintf("%s=%d", ortEnv.ORACLERT_NOORACLE, 1))
	}
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

		if errors.Is(runErr, context.DeadlineExceeded) {
			timeout = true
		}
	}

	// Read oracle runtime output
	ortOutBytes, err := ioutil.ReadFile(ortOutputFp)
	var ortOutput *ortOut.Output
	if err != nil {
		logger.Printf("[input %s][ignored] cannot read file %s: %v", input.ID, ortOutputFp, err)

	} else {
		ortOutput, err = ortOut.Deserialize(ortOutBytes)
		if err != nil {
			logger.Printf("[input %s][ignored] failed to parse file %s: %v", input.ID, ortOutputFp, err)
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

	execOutput := &api.Output{
		BugIDs:         bugIDs,
		Timeout:        timeout,
		OracleRtOutput: ortOutput,
	}

	return execOutput, nil
}

// HandleExec is called immediately after running execution.
// It should take care of following things:
// 1. ignore this exec if timeout (fixme: more flexible options like increasing timeout and try next time?)
// 2. check if any unique bugs found
// 3. update if any new select records found
// 4. update/add interest input
func HandleExec(ctx context.Context, i *api.Input, o *api.Output, fctx *api.Context, interestHdl api.InterestHandler) error {
	if o.OracleRtOutput == nil {
		return fmt.Errorf("cannot handle an exec without oracle runtime output")
	}

	entry := fctx.GetQueueEntryByGExecID(i.Exec.String())
	if entry == nil {
		return fmt.Errorf("cannot find queue entry for %s", i.Exec.String())
	}

	logger := getWorkerLogger(ctx)
	// 1. ignore this exec if timeout
	if o.Timeout {

		if i.OracleRtConfig != nil && i.OracleRtConfig.SelEfcm.Efcms != nil {
			cfgHash := hash.AsSha256(i.OracleRtConfig)
			entry.RecordTimeoutOrtCfgHash(cfgHash)
		}

		// update init record in interest queue
		if i.Stage == api.InitStage {
			ii := fctx.Interests.FindInit(i)
			ii.Timeout = true
		}
		return fmt.Errorf("ignore because of timeout")
	}

	// 2. check if any unique bugs found
	// Check any unique bugs found
	numOfBugs := 0
	for _, bugID := range o.BugIDs {
		if !fctx.HasBugID(bugID) {
			stdout, _ := i.GetOutputFilePath()
			fctx.AddBugID(bugID, stdout)
			numOfBugs += 1
			logger.Printf("found unique bug: %s", bugID)
		}
	}

	if numOfBugs != 0 {
		logger.Printf("found %d unique bug(s)\n", numOfBugs)

	}

	// 3. update if any new select records found
	isFoundNewSelect := false

	newSelects := entry.UpdateSelectRecordsIfNew(o.OracleRtOutput.Selects)
	if newSelects != 0 {
		logger.Printf("found %d new selects", newSelects)
		isFoundNewSelect = true
	}

	// 4. update/add interest input
	if i.Stage == api.InitStage {
		// if input is init, since init stage by default is interested input, so no need to check interest
		// simply update output and hash
		ii := fctx.Interests.FindInit(i)
		ii.Output = o

		// Use SELECT Record as feedback.
		oSelectCopy := make([]ortOut.SelectRecord, 0)
		for _, v := range o.OracleRtOutput.Selects {
			v_copy := v
			v_copy.Cases = 0
			oSelectCopy = append(oSelectCopy, v_copy)
		}
		fctx.UpdateOrtOutputHash(hash.AsSha256(oSelectCopy))
		return nil
	}

	if interestHdl != nil {
		isInteresed, reason, err := interestHdl.IsInterested(i, o, isFoundNewSelect)
		if err != nil {
			return nil
		}
		if isInteresed {
			logger.Printf("%s is interesting", i.ID)
			ii := api.NewExecutedInterestInput(i, o)
			ii.Reason = reason
			fctx.Interests.Add(ii)
		}
	}

	return nil
}
