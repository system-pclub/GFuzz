package interest

import (
	"fmt"
	"gfuzz/pkg/fuzz/api"
	"gfuzz/pkg/fuzz/gexecfuzz"
	"gfuzz/pkg/fuzz/mutate"
	"gfuzz/pkg/fuzz/score"
	"gfuzz/pkg/utils/bits"
	"gfuzz/pkg/utils/hash"
	"gfuzz/pkg/utils/rand"
	"log"
	"strconv"
	"strings"
)

type InterestHandlerImpl struct {
	fctx *api.Context
}

func NewInterestHandlerImpl(fctx *api.Context) api.InterestHandler {
	return &InterestHandlerImpl{
		fctx: fctx,
	}
}

//todo: move new select detection from exec.go to here
func (h *InterestHandlerImpl) IsInterested(i *api.Input, o *api.Output, isFoundNewSelect bool) (bool, api.InterestReason, error) {

	// If isIgnoreFeedback is true, we treat every feedback as interesting and directly return.
	if h.fctx.Cfg.IsIgnoreFeedback {
		return true, api.NoInterest, nil
	}

	var reason api.InterestReason

	isInteresting := false

	if isFoundNewSelect {
		isInteresting = true
		reason = api.InterestReason(bits.Set(bits.Bits(reason), bits.Bits(api.NewSelectFound)))
	}

	if !h.fctx.Cfg.NoSelEfcm && !IsEfcmCovered(i.OracleRtConfig.SelEfcm.Efcms, o.OracleRtOutput.Selects) {
		reason = api.InterestReason(bits.Set(bits.Bits(reason), bits.Bits(api.SelEfcmNotCovered)))
		isInteresting = true
	}

	// Check new tuple
	entry := h.fctx.GetQueueEntryByGExecID(i.Exec.String())
	if entry != nil && entry.UpdateTupleRecordsIfNew(o.OracleRtOutput.Tuples) > 0 {
		// Has new tuples, interesting
		isInteresting = true
		reason = api.InterestReason(bits.Set(bits.Bits(reason), bits.Bits(api.NewTuple)))
	}

	// See if we found new channel or new channel state.
	if entry != nil && entry.UpdateChannelRecordsIfNew(o.OracleRtOutput.Channels) > 0 {
		// Has new channels, interesting
		isInteresting = true
		reason = api.InterestReason(bits.Set(bits.Bits(reason), bits.Bits(api.NewChannel)))
	}

	// Using SELECT record as feedback
	oSelectMap := make(map[string][]uint)
	for _, v := range o.OracleRtOutput.Selects {
		vSelectId := v.ID
		vSelectChosen := v.Chosen

		if val, ok := oSelectMap[vSelectId]; ok {
			val = append(val, vSelectChosen)
			oSelectMap[vSelectId] = val
		} else {
			val := []uint{vSelectChosen}
			oSelectMap[vSelectId] = val
		}
	}

	if len(oSelectMap) == 0 {
		/* No selects detected. Return not interesting.  */
		return false, api.NoInterest, nil
	}

	selectHash := hash.AsSha256(oSelectMap)
	// fixme: should be entrypoint based right?
	if h.fctx.UpdateOrtOutputHash(selectHash) {
		isInteresting = true
		reason = api.InterestReason(bits.Set(bits.Bits(reason), bits.Bits(api.Other)))
	}
	if isInteresting {
		return true, reason, nil
	} else {
		return false, api.NoInterest, nil
	}
}

func (h *InterestHandlerImpl) HandleInterest(i *api.InterestInput) (ret bool, err error) {

	// if init input has not been executed, execute first
	if i.Input.Stage == api.InitStage && !i.Executed {
		i.HandledCnt += 1
		i.Executed = true
		h.fctx.ExecInputCh <- i.Input
		return true, nil
	}

	if i.Output == nil {
		// if executed is true but output is nil
		// it could be still in queue, pending to run
		return false, nil
	}

	// if interested input has been executed, then try to mutate and send to execution according to its stage
	switch i.Input.Stage {
	case api.InitStage:
		// we are handling the output from the input with init stage

		//if i.HandledCnt == 1 {
		//	// if handle init first time, treat it as normal init stage
		//	err = handleInitStageInput(h.fctx, i.Input, i.Output)
		//} else {
		//	// if not the first time, treat it as random
		//}

		// Yu:: Directly jump to RandStage if we see the init_stage for the second time.
		// Yu :: The first time of execution should be covered by !i.Executed.
		ret, err = handleRandStageInput(h.fctx, i)
	case api.DeterStage:
		// we are handling the output from the input with deter stage
		//err = handleDeterStageInput(h.fctx, i.Input, i.Output)

		// Yu:: Should not be seen in the exec. But if seen, treat it as rand.
		ret, err = handleRandStageInput(h.fctx, i)
	case api.CalibStage:
		// we are handling the output from the input with calib stage
		//err = handleCalibStageInput(h.fctx,  i)

		// Yu:: Should not be seen in the exec. But if seen, treat it as rand.
		ret, err = handleRandStageInput(h.fctx, i)
	case api.RandStage:
		// we are handling the output from the input with rand stage
		ret, err = handleRandStageInput(h.fctx, i)
	case api.ReplayStage:
		// no need to handle replay

	default:
		err = fmt.Errorf("unexpected stage: %s", i.Input.Stage)
	}

	if err != nil {
		return false, err
	}

	i.HandledCnt += 1

	return ret, nil

}

func handleInitStageInput(fctx *api.Context, i *api.Input, o *api.Output) (bool, error) {

	//g := fctx.GetQueueEntryByGExecID(i.Exec.String())
	//execID, err := getExecIDFromInputID(i.ID)
	//if err != nil {
	//	return false, err
	//}
	//var deterInputs []*api.Input
	//var mts mutate.OrtConfigMutateStrategy = &mutate.DeterMutateStrategy{}
	//
	//if o.OracleRtOutput == nil {
	//	return false, nil
	//}
	//cfgs, err := mts.Mutate(g, i.OracleRtConfig, o.OracleRtOutput, 0)
	//if err != nil {
	//	return false, err
	//}
	//
	//for _, cfg := range cfgs {
	//	deterInputs = append(deterInputs, api.NewExecInput(fctx.GetAutoIncGlobalID(), execID, fctx.Cfg.OutputDir, g.Exec, cfg, api.DeterStage))
	//}
	//
	//if len(deterInputs) == 0 {
	//	return false, nil
	//}
	//for _, input := range deterInputs {
	//	fctx.ExecInputCh <- input
	//}

	return true, nil
}

func handleDeterStageInput(fctx *api.Context, i *api.Input, o *api.Output) (bool, error) {
	g := fctx.GetQueueEntryByGExecID(i.Exec.String())
	execID, err := getExecIDFromInputID(i.ID)
	if err != nil {
		return false, err
	}

	input := api.NewExecInput(fctx.GetAutoIncGlobalID(), execID, fctx.Cfg.OutputDir, g.Exec, i.OracleRtConfig, api.CalibStage)
	fctx.ExecInputCh <- input
	return true, nil
}

func handleCalibStageInput(fctx *api.Context, i *api.Input, o *api.Output) (bool, error) {
	g := fctx.GetQueueEntryByGExecID(i.Exec.String())
	execID, err := getExecIDFromInputID(i.ID)
	if err != nil {
		return false, err
	}

	input := api.NewExecInput(fctx.GetAutoIncGlobalID(), execID, fctx.Cfg.OutputDir, g.Exec, i.OracleRtConfig, api.RandStage)
	fctx.ExecInputCh <- input
	return true, nil
}

// segment chance into 4 interval
func segmentChance(chance int) int {
	if chance <= 20 {
		return 20
	} else if chance <= 40 {
		return 40
	} else if chance <= 60 {
		return 60
	} else if chance <= 80 {
		return 80
	}

	return 100
}

func handleRandStageInput(fctx *api.Context, ii *api.InterestInput) (bool, error) {
	i := ii.Input
	o := ii.Output
	g := fctx.GetQueueEntryByGExecID(i.Exec.String())
	execID, err := getExecIDFromInputID(i.ID)

	if err != nil {
		return false, err
	}

	if !fctx.Cfg.IsIgnoreFeedback && !fctx.Cfg.NoSelEfcm {
		// if this interest does not covered the enforcement, rerun
		if bits.Has(bits.Bits(ii.Reason), bits.Bits(api.SelEfcmNotCovered)) {
			newCfg := i.OracleRtConfig.Copy()
			newCfg.SelEfcm.SelTimeout += 3000
			if newCfg.SelEfcm.SelTimeout < 12500 {
				log.Printf("handle %d, new rerun with timeout %d becuase of uncovered efcm", execID, newCfg.SelEfcm.SelTimeout)
				ni := api.NewExecInput(fctx.GetAutoIncGlobalID(), execID, fctx.Cfg.OutputDir, g.Exec, newCfg, api.RandStage)
				fctx.ExecInputCh <- ni
			}
		}

		if !bits.Has(bits.Bits(ii.Reason), bits.Bits(api.NewChannel)) &&
			!bits.Has(bits.Bits(ii.Reason), bits.Bits(api.NewSelectFound)) &&
			!bits.Has(bits.Bits(ii.Reason), bits.Bits(api.NewTuple)) &&
			!bits.Has(bits.Bits(ii.Reason), bits.Bits(api.InitStg)) &&
			!bits.Has(bits.Bits(ii.Reason), bits.Bits(api.Other)) {
			log.Printf("handle %v, skip mutating since no other interest reason", execID)
			return true, nil
		}
	}

	var randInputs []*api.Input
	var mts mutate.OrtConfigMutateStrategy

	if fctx.Cfg.IsIgnoreFeedback {
		mts = &mutate.NfbRandomMutateStrategy{
			SelEfcmTimeout:      fctx.Cfg.SelEfcmTimeout,
			FixedSelEfcmTimeout: fctx.Cfg.FixedSelEfcmTimeout,
			RandomTimeoutIncr:   fctx.Cfg.NfbRandSelEfcmTimeout,
		}
	} else {
		// in feedback mode
		if len(o.OracleRtOutput.Selects) <= 5 && fctx.Cfg.MemRandStrat {
			mts = &mutate.MemRandMutateStrategy{
				SelEfcmTimeout:      fctx.Cfg.SelEfcmTimeout,
				FixedSelEfcmTimeout: fctx.Cfg.FixedSelEfcmTimeout,
			}
		} else if fctx.Cfg.NoSelEfcm {
			mts = &mutate.NoMutateStrategy{}
		} else {
			mts = &mutate.RandomMutateStrategy{
				SelEfcmTimeout:      fctx.Cfg.SelEfcmTimeout,
				FixedSelEfcmTimeout: fctx.Cfg.FixedSelEfcmTimeout,
			}
		}

	}

	mutateEnergy := fctx.Cfg.RandMutateEnergy
	// if nfb, random 1-5 to be our energy
	if fctx.Cfg.IsIgnoreFeedback && fctx.Cfg.NfbRandEnergy {
		mutateEnergy = rand.GetRandomWithMax(5) + 1
	}

	var scoreFunc = score.NewScoreStrategyImpl(fctx)
	curScore, _ := scoreFunc.Score(i, o)
	if curScore > fctx.GlobalBestScore {
		fctx.GlobalBestScore = curScore
	}
	if !fctx.Cfg.IsDisableScore && fctx.GlobalBestScore >= 100 {

		origChance := int(100.0 * (float64(curScore) / float64(fctx.GlobalBestScore)))
		if origChance == 0 {
			// if chance is less than 1%, return
			log.Printf("handle %d, skip because of 0 score", execID)
			return true, nil
		}
		randMutateChance := segmentChance(origChance)
		log.Printf("handle %d, current score %d, max score %d, execution chance %d%%(%d%%)",
			execID, curScore, fctx.GlobalBestScore, randMutateChance, origChance)
		if fctx.Cfg.ScoreBasedEnergy {
			mutateEnergy = int(randMutateChance / 20)
			log.Printf("handle %d, energy %d", execID, mutateEnergy)
		} else if rand.GetRandomWithMax(100) >= randMutateChance {
			// Skip the test case based on rand possibilities.
			log.Printf("handle %d, skip because of score", execID)
			// add it back to interest queue since right now queue is not persistant
			if ii.Input.Stage != api.InitStage {
				fctx.Interests.Add(ii)
			}
			return true, nil
		}
	}

	var caseNums []int
	for _, sel := range o.OracleRtOutput.Selects {
		caseNums = append(caseNums, int(sel.Cases))
	}

	log.Printf("handle %d, cases: %v", execID, caseNums)

	cfgs, err := mts.Mutate(g, i.OracleRtConfig, o.OracleRtOutput, mutateEnergy)

	if err != nil {
		return false, err
	}

	if len(cfgs) == 0 {
		// if no configuration, it implies there is no selects for fuzzer to mutate
		// if no feedback mode,  we should not consider it and simply rerun it
		if fctx.Cfg.IsIgnoreFeedback {
			for idx := 0; idx < mutateEnergy; idx++ {
				cfgs = append(cfgs, i.OracleRtConfig.Copy())
			}
		}
	}
	for _, cfg := range cfgs {
		cfgHash := hash.AsSha256(cfg)
		if !fctx.Cfg.IsIgnoreFeedback {
			if g.HasTimeoutOrtCfgHash(cfgHash) {
				log.Printf("handle %d, skip a generated config because of timeout", execID)
				continue
			}

			if !fctx.Cfg.AllowDupCfg && g.HasOrtCfgHash(cfgHash) {
				log.Printf("handle %d, skip a generated config because of duplication", execID)
				continue
			}
		}

		g.RecordOrtCfgHash(cfgHash)
		randInputs = append(randInputs, api.NewExecInput(fctx.GetAutoIncGlobalID(), execID, fctx.Cfg.OutputDir, g.Exec, cfg, api.RandStage))
	}

	for _, input := range randInputs {
		fctx.ExecInputCh <- input
	}

	return true, nil
}

func getExecIDFromInputID(inputID string) (uint32, error) {
	id, err := strconv.Atoi(strings.Split(inputID, "-")[0])
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}

func (h *InterestHandlerImpl) CleanAllGExecsRecords() (err error) {
	h.fctx.EachGExecFuzz(func(gf *gexecfuzz.GExecFuzz) {
		gf.Clean()
	})
	return nil
}
