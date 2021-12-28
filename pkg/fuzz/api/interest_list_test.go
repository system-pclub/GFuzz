package api

import (
	"gfuzz/pkg/fuzz/api"
	"testing"
	"time"
)

type TestInterestHandler struct {
	shoulTouch map[*InterestInput]bool
	t          *testing.T
}

func (m *TestInterestHandler) IsInterested(i *Input, o *Output, isFoundNewSelect bool) (bool, error) {
	return false, nil
}

func (m *TestInterestHandler) HandleInterest(i *InterestInput) (bool, api.InterestReason, error) {
	if touched, existed := m.shoulTouch[i]; touched || !existed {
		m.t.Fail()
	}
	m.shoulTouch[i] = true
	return true, api.NoInterest, nil
}

func (m *TestInterestHandler) Check() {
	for _, touched := range m.shoulTouch {
		if !touched {
			m.t.Fail()
		}
	}
}

func TestConcurrentInterestListEachAndAdd(t *testing.T) {
	it := &InterestList{}

	input1 := &InterestInput{
		Input: &Input{
			Stage: InitStage,
		},
	}
	input2 := &InterestInput{
		Input: &Input{
			Stage: InitStage,
		},
	}
	input3 := &InterestInput{
		Input: &Input{
			Stage: InitStage,
		},
	}

	it.Add(input1)
	it.Add(input2)
	it.Add(input3)

	th := &TestInterestHandler{
		shoulTouch: make(map[*InterestInput]bool),
		t:          t,
	}

	th.shoulTouch[input1] = false
	th.shoulTouch[input2] = false
	th.shoulTouch[input3] = false

	done := make(chan struct{}, 1)
	go func() {
		// make sure it.Each called first
		time.Sleep(100 * time.Millisecond)
		it.Add(&InterestInput{
			Input: &Input{
				Stage: DeterStage,
			},
		})
		it.Add(&InterestInput{
			Input: &Input{
				Stage: DeterStage,
			},
		})
		it.Add(&InterestInput{
			Input: &Input{
				Stage: DeterStage,
			},
		})
		done <- struct{}{}
	}()
	it.Each(th)
	th.Check()
	<-done
	if len(it.interestedInputs) != 3 {
		t.Fail()
	}

	if len(it.initInputs) != 3 {
		t.Fail()
	}
}

func TestInterestListEach(t *testing.T) {
	it := &InterestList{}

	input1 := &InterestInput{
		Input: &Input{
			Stage: InitStage,
		},
	}
	input2 := &InterestInput{
		Input: &Input{
			Stage: InitStage,
		},
	}
	input3 := &InterestInput{
		Input: &Input{
			Stage: InitStage,
		},
	}

	it.Add(input1)
	it.Add(input2)
	it.Add(input3)

	th := &TestInterestHandler{
		shoulTouch: make(map[*InterestInput]bool),
		t:          t,
	}

	th.shoulTouch[input1] = false
	th.shoulTouch[input2] = false
	th.shoulTouch[input3] = false

	it.Each(th)
	th.Check()
}

func TestInterestListAdd(t *testing.T) {
	it := &InterestList{}

	input1 := &InterestInput{
		Input: &Input{
			Stage: InitStage,
		},
	}
	input2 := &InterestInput{
		Input: &Input{
			Stage: DeterStage,
		},
	}
	th := &TestInterestHandler{
		shoulTouch: make(map[*InterestInput]bool),
		t:          t,
	}
	th.shoulTouch[input1] = false
	th.shoulTouch[input2] = false
	it.Add(input1)

	if it.initInputs[0] != input1 {
		t.Fail()
	}

	if len(it.interestedInputs) != 0 {
		t.Fail()
	}

	it.Each(th)

	it.Add(input2)
	if it.interestedInputs[0] != input2 {
		t.Fail()
	}
	it.Each(th)
	if len(it.interestedInputs) != 0 {
		t.Fail()
	}

}
