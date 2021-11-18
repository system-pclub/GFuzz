package api

import (
	"sync"
)

type InterestList struct {
	// interestedInputs contains a list of interested input
	interestedInputs []*InterestInput
	rw               sync.RWMutex
}

func (i *InterestList) Add(input *InterestInput) {
	i.rw.Lock()
	defer i.rw.Unlock()
	i.interestedInputs = append(i.interestedInputs, input)
}

func (i *InterestList) Each(handler InterestHandler) {
	i.rw.RLock()
	defer i.rw.RUnlock()
	for _, input := range i.interestedInputs {
		handler.HandleInterest(input)
	}
}
