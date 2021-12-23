package interest

import (
	"fmt"
	"gfuzz/pkg/fuzz/api"
	"gfuzz/pkg/fuzz/mutate"
	"gfuzz/pkg/fuzz/score"
	"gfuzz/pkg/utils/hash"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type InterestHandlerImpl struct {
	fctx *api.Context
}

func NewInterestHandlerImpl(fctx *api.Context) api.InterestHandler {
	return &InterestHandlerImpl{
		fctx: fctx,
	}
}
func (h *InterestHandlerImpl) IsInterested(i *api.Input, o *api.Output, isFoundNewSelect bool) (bool, error) {

	// If isIgnoreFeedback is true, we treat every feedback as interesting and directly return.
	if h.fctx.Cfg.IsIgnoreFeedback {
		return true, nil
	}

	isInteresting := false

	// Check new tuple
	entry := h.fctx.GetQueueEntryByGExecID(i.Exec.String())
	if entry != nil && entry.UpdateTupleRecordsIfNew(o.OracleRtOutput.Tuples) > 0 {
		// Has new tuples, interesting
		isInteresting = true
	}

	// See if we found new channel or new channel state.
	if entry != nil && entry.UpdateChannelRecordsIfNew(o.OracleRtOutput.Channels) > 0 {
		// Has new channels, interesting
		isInteresting = true
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
		return false, nil
	}

	ortCfgHash := hash.AsSha256(oSelectMap)
	if h.fctx.UpdateOrtOutputHash(ortCfgHash) {
		isInteresting = true
	}
	if isInteresting || isFoundNewSelect {
		return true, nil
	} else {
		return false, nil
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

	// If IsNoMuation, then we do not mutate the seeds. Directly runs the seed.
	if h.fctx.Cfg.IsNoMutation {
		i.HandledCnt += 1
		h.fctx.ExecInputCh <- i.Input
		return true, nil
	}

	if i.Output == nil {
		// if executed is true but output is nil
		// it could be still in channel pending to run
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
	var randInputs []*api.Input
	var mts mutate.OrtConfigMutateStrategy = &mutate.RandomMutateStrategy{}

	randMutateEnergy := fctx.Cfg.RandMutateEnergy

	if !fctx.Cfg.IsDisableScore && fctx.GlobalBestScore >= 100 {
		var scoreFunc = score.NewScoreStrategyImpl(fctx)
		curScore, _ := scoreFunc.Score(i, o)
		if curScore > fctx.GlobalBestScore {
			fctx.GlobalBestScore = curScore
		}
		origChance := int(100.0 * (float64(curScore) / float64(fctx.GlobalBestScore)))
		randMutateChance := segmentChance(origChance)
		log.Printf("handle %d, current score %d, max score %d, execution chance %d%%(%d%%)", execID, curScore, fctx.GlobalBestScore, randMutateChance, origChance)
		rand.Seed(time.Now().UnixNano())
		if rand.Intn(100) >= randMutateChance {
			// Skip the test case based on rand possibilities.
			log.Printf("handle %d, skip because of score", execID)
			// add it back to interest queue since right now queue is not persistant
			if ii.Input.Stage != api.InitStage {
				fctx.Interests.Add(ii)
			}
			return true, nil
		}
	}

	cfgs, err := mts.Mutate(g, i.OracleRtConfig, o.OracleRtOutput, randMutateEnergy)

	if err != nil {
		return false, err
	}

	if len(cfgs) == 0 {
		// if no configuration, it implies there is no selects for fuzzer to mutate
		// if no feedback mode,  we should not consider it and simply rerun it
		if fctx.Cfg.IsIgnoreFeedback {
			cfgs = append(cfgs, i.OracleRtConfig)
		}
	}
	for _, cfg := range cfgs {

		if !fctx.Cfg.IsIgnoreFeedback {
			if g.HasTimeoutEfcm(hash.AsSha256(cfg.SelEfcm.Efcms)) {
				log.Printf("handle %d, skip a generated config becuase of timeout", execID)
				continue
			}

			if g.HasOrtCfgHash(hash.AsSha256(cfg)) {
				log.Printf("handle %d, skip a generated config becuase of duplication", execID)
				continue
			}

			/* Check select enforcements, see if they are redundant cases. */
			selectEfcms := cfg.SelEfcm.Efcms
			if g.UpdateInputSelectEfcmsIfNew(selectEfcms) > 0 {
				cfgs = append(cfgs, cfg)
			} else if getRandomWithMax(10) < 1 {
				/* If this is redundant case, give 10% chance to rerun. */
				cfgs = append(cfgs, cfg)
			} else if mutate_idx == (energy-1) && len(cfgs) == 0 {
				/* Keep one case */
				//cfgs = append(cfgs, cfg)
				return nil, nil
			}
		}

		
		g.RecordOrtCfgHash(hash.AsSha256(cfg))
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
