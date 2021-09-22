package fuzz

import (
	"container/list"
	"gfuzz/pkg/exec"
	"sync"
	"sync/atomic"
	"time"
)

// Context record all necessary information for help fuzzer to prioritize input and record process.
type Context struct {
	q                *list.List
	qLock            sync.RWMutex // lock for fuzzingQueue
	mainRecord       *Record
	allRecordHashMap map[string]struct{}

	// Metrics
	numOfBugsFound      uint64
	numOfRuns           uint64
	numOfFuzzQueryEntry uint64
	numOfTargets        uint64
	startAt             time.Time

	// timeout counter: src => how many times timeout when running this src
	// if more than 3, drop it in the queue
	timeoutTargets     map[string]uint32
	timeoutTargetsLock sync.RWMutex
}

// NewContext returns a new FuzzerContext
func NewContext(execs []exec.Executable) *Context {
	q := list.New()
	for e := range execs {
		q.PushBack(e)
	}
	return &Context{
		q:                q,
		mainRecord:       EmptyRecord(),
		allRecordHashMap: make(map[string]struct{}),
		timeoutTargets:   make(map[string]uint32),
		startAt:          time.Now(),
	}
}

func (c *Context) IncNumOfRun() {
	atomic.AddUint64(&c.numOfRuns, 1)
}

func (c *Context) IncNumOfBugsFound(num uint64) {
	atomic.AddUint64(&c.numOfBugsFound, num)
}

func (c *Context) NewFuzzQueryEntryIndex() uint64 {
	return atomic.AddUint64(&c.numOfFuzzQueryEntry, 1)
}

func (c *Context) HasBugID(id string) bool {
	c.bugID2FpLock.RLock()
	_, exists := c.allBugID2Fp[id]
	c.bugID2FpLock.RUnlock()
	return exists
}

func (c *Context) AddBugID(bugID string, filepath string) {
	c.bugID2FpLock.Lock()
	defer c.bugID2FpLock.Unlock()
	c.allBugID2Fp[bugID] = &BugMetrics{
		FoundAt: time.Now(),
		Stdout:  filepath,
	}

}
func (c *Context) UpdateTargetMaxCaseCov(target string, caseCov float32) {
	c.targetStagesLock.Lock()
	defer c.targetStagesLock.Unlock()

	var track *TargetMetrics
	var exist bool
	if track, exist = c.targetStages[target]; !exist {
		track = &TargetMetrics{
			At: make(map[FuzzStage]time.Time),
		}
		c.targetStages[target] = track
	}

	if caseCov > track.MaxCaseCov {
		track.MaxCaseCov = caseCov
	}
}

func (c *Context) RecordTargetTimeoutOnce(target string) {
	c.timeoutTargetsLock.Lock()
	defer c.timeoutTargetsLock.Unlock()

	c.timeoutTargets[target] += 1
}
