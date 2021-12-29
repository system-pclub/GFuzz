package api

import (
	"gfuzz/pkg/utils/bits"
)

// ScoreStrategy is used in interest
type ScoreStrategy interface {
	// Score will score an exec input
	Score(i *Input, o *Output) (int, error)
}

type InterestReason bits.Bits

const (
	NoInterest        InterestReason = 0
	SelEfcmNotCovered InterestReason = 1 << iota
	NewSelectFound
	NewTuple
	NewChannel
	InitStg // inil stage always interest
	Other   //todo: should be replaced by more detailed description
)

// InterestInput is created if the input is interested(identify by score strategy) or it is init stage input
type InterestInput struct {
	Executed   bool
	Timeout    bool // Timeout usually only used in init stage, since other stage timeout will not be added into interest queue.
	HandledCnt uint32
	Reason     InterestReason
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
	IsInterested(i *Input, o *Output, isFoundNewSelect bool) (bool, InterestReason, error)

	CleanAllGExecsRecords() error
}
