package api

import (
	"log"
	"sync"
	"sync/atomic"
)

type InterestList struct {
	// interestedInputs contains a list of interested input
	interestedInputs []*InterestInput
	rw               sync.RWMutex
	Dirty            bool   // inputs changed before HandleEach
	looping          uint32 // atomic boolean
}

func (i *InterestList) Add(input *InterestInput) {
	i.rw.Lock()
	defer i.rw.Unlock()
	if !i.Dirty {
		i.Dirty = true
	}
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

func (i *InterestList) IsLooping() bool {
	return atomic.LoadUint32(&i.looping) == 1
}

func (i *InterestList) Each(handler InterestHandler) {
	i.rw.Lock()
	currInterests := make([]*InterestInput, len(i.interestedInputs))
	i.Dirty = false
	copy(currInterests, i.interestedInputs)
	i.looping += 1
	i.rw.Unlock()
	log.Printf("interesting list length: %d", len(currInterests))
	for _, i := range currInterests {
		handler.HandleInterest(i)
	}

	i.looping -= 1

}
