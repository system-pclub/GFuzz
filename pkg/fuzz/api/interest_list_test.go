package api

import (
	"testing"
	"time"
)

type TestInterestHandler struct {
	shoulTouch map[*InterestInput]bool
	t          *testing.T
}

func (m *TestInterestHandler) IsInterested(i *Input, o *Output) (bool, error) {
	return false, nil
}

func (m *TestInterestHandler) HandleInterest(i *InterestInput) (bool, error) {
	if touched, existed := m.shoulTouch[i]; touched || !existed {
		m.t.Fail()
	}
	m.shoulTouch[i] = true
	return true, nil
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

	input1 := &InterestInput{}
	input2 := &InterestInput{}
	input3 := &InterestInput{}

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
		it.Add(&InterestInput{})
		it.Add(&InterestInput{})
		it.Add(&InterestInput{})
		done <- struct{}{}
	}()
	it.Each(th)
	th.Check()
	<-done
	if len(it.interestedInputs) != 6 {
		t.Fail()
	}
}
