package api

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
)

type InterestList struct {
	// interestedInputs contains a list of interested input
	interestedInputs []*InterestInput
	// initInputs contains only init stage
	initInputs []*InterestInput
	rw         sync.RWMutex
	Dirty      bool   // inputs changed before HandleEach
	looping    uint32 // atomic boolean
}

func (i *InterestList) Add(input *InterestInput) error {
	i.rw.Lock()
	defer i.rw.Unlock()
	if !i.Dirty {
		i.Dirty = true
	}
	if input == nil || input.Input == nil {
		return fmt.Errorf("input is nil")
	}
	if input.Input.Stage == InitStage {
		i.initInputs = append(i.initInputs, input)
	} else {
		i.interestedInputs = append(i.interestedInputs, input)
	}
	return nil
}

func (i *InterestList) FindInit(input *Input) *InterestInput {
	// only init stage's interest need/should be updated/fetched
	i.rw.Lock()
	defer i.rw.Unlock()
	for _, ii := range i.initInputs {
		if ii.Input == input {
			return ii
		}
	}
	return nil
}

func (i *InterestList) IsLooping() bool {
	return atomic.LoadUint32(&i.looping) == 1
}

func (i *InterestList) GetInterestingLength() int {
	return len(i.interestedInputs)
}

// Each loops the whole interest queue once
func (i *InterestList) Each(handler InterestHandler) (ret bool) {
	i.rw.Lock()
	var currInterests []*InterestInput

	if len(i.interestedInputs) == 0 {
		currInterests = i.initInputs
	} else {
		currInterests = make([]*InterestInput, len(i.interestedInputs))
		copy(currInterests, i.interestedInputs)
		// clear interest inputs list each time after copying all of them
		i.interestedInputs = nil
	}

	// if current interest queue is too short, loop init also
	if len(currInterests) < len(i.initInputs)/2 {
		currInterests = append(currInterests, i.initInputs...)
		log.Printf("handling interest: loop init becuase of short interest list")
	}
	i.Dirty = false
	i.looping += 1
	i.rw.Unlock()
	log.Printf("interesting list length: %d", len(currInterests))
	for _, i := range currInterests {
		if i.Timeout {
			log.Printf("handling interest %s: skip because of it marked with timeout before", i.Input.ID)
			continue
		}
		handled, err := handler.HandleInterest(i)
		if err != nil {
			log.Printf("handling interest %s: %s", i.Input.ID, err)
		}
		ret = handled || ret
	}

	i.looping -= 1
	return ret
}
