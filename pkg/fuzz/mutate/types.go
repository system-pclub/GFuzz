package mutate

import "gfuzz/pkg/oraclert/config"

type RtConfigMutateStrategy interface {
	Mutate(old *config.Config) (*config.Config, error)
}
