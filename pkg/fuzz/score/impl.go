package score

import (
	"gfuzz/pkg/fuzz/api"
	"math"
)

const (
	ScoreTupleCountLog2 = 1
	ScoreCh             = 10
	ScoreNewClosed      = 10
	ScoreNewNotClosed   = 10
	ScoreBuf            = 10
)

type ScoreStrategyImpl struct {
	fctx *api.Context
}

func NewScoreStrategyImpl(fctx *api.Context) api.ScoreStrategy {
	return &ScoreStrategyImpl{
		fctx: fctx,
	}
}

func (s *ScoreStrategyImpl) Score(i *api.Input, o *api.Output) (int, error) {
	var totalScore int = 0
	var tupleScore int = 0
	//var chNotclosedScore int = 0
	var chClosedScore int = 0
	var bufferScore int = 0
	var channelScore int = 0

	/* Calculate tuple score */
	for _, count := range o.OracleRtOutput.Tuples {
		tupleScore += int(math.Log2(float64(count))) * ScoreTupleCountLog2
	}

	for _, chRecord := range o.OracleRtOutput.Channels {

		/* Calculate score for first time closed or not closed channel */
		//if chRecord.NotClosed {
		//	chNotclosedScore += ScoreNewNotClosed
		//}

		if chRecord.Closed {
			chClosedScore += ScoreNewClosed
		}

		//TODO:: Missing first time closed?

		/* Calculate score for score buffer */
		if chRecord.PeakBuf > 0 && chRecord.CapBuf != 0 {
			bufferPer := float64(chRecord.PeakBuf) / float64(chRecord.CapBuf)
			bufferScore += int(float64(ScoreBuf) * bufferPer)
		}

		/* Calculate score for each encountered channel */
		channelScore += ScoreCh
	}

	totalScore += tupleScore
	totalScore += chClosedScore
	totalScore += bufferScore
	totalScore += channelScore

	return totalScore, nil
}
