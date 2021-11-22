package score

import (
	"gfuzz/pkg/fuzz/api"
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
	return 101, nil
}
