package api

// ScoreStrategy is used in interest
type ScoreStrategy interface {
	// Score will score an exec input
	Score(i *Input, o *Output) (int, error)
}

// InterestInput is created if the input is interested(identify by score strategy) or it is init stage input
type InterestInput struct {
	Executed   bool
	HandledCnt uint32
	Input      *Input
	Output     *Output // Output will be nil if executed is false
}

// InterestHandler is used to handle InterestInput in following:
// sends deter or rand mutate inputs to fuzz context's channel
type InterestHandler interface {
	// HandleInterest will be called when handling.
	// Return true if something handled, false otherwise(such as no selects to mutate, no oracle runtime output to handle, etc.)
	HandleInterest(i *InterestInput) (bool, error)
	// IsInterested will decide if this execution should be added into interest list
	IsInterested(i *Input, o *Output, isFoundNewSelect bool) (bool, error)
}
