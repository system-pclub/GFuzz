package api

// ScoreStrategy is used in interest
type ScoreStrategy interface {
	// Score will score an exec input
	Score(i *Input, o *Output) (int, error)
	// InterestScore will return the a score. If exec input whose score above the returned value, this input will be added into interested exec input list
	InterestScore() int
}

// InterestInput is created if the input is interested(identify by score strategy) or it is init stage input
type InterestInput struct {
	Executed bool
	Input    *Input
	Output   *Output // Output will be nil if executed is false
}

// InterestHandler is used to handle InterestInput in following:
// sends deter or rand mutate inputs to fuzz context's channel
type InterestHandler interface {
	// HandleInterest will be called when handling
	HandleInterest(i *InterestInput) error
}
