package api

import (
	"log"
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

func (i *InterestList) Find(input *Input) *InterestInput {
	i.rw.Lock()
	defer i.rw.Unlock()
	for _, ii := range i.interestedInputs {
		if ii.Input == input {
			return ii
		}
	}
	return nil
}

func (i *InterestList) Each(handler InterestHandler) {
	i.rw.RLock()
	currInterests := make([]*InterestInput, len(i.interestedInputs))
	copy(currInterests, i.interestedInputs)
	i.rw.RUnlock()
	log.Printf("interest list length: %d\n", len(currInterests))
	for _, i := range currInterests {
		handler.HandleInterest(i)
	}
}
