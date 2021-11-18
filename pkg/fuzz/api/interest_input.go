package api

func NewUnexecutedInterestInput(i *Input) *InterestInput {
	return &InterestInput{
		Input:    i,
		Executed: false,
		Output:   nil,
	}
}

func NewExecutedInterestInput(i *Input, o *Output) *InterestInput {
	return &InterestInput{
		Input:    i,
		Executed: true,
		Output:   o,
	}
}
