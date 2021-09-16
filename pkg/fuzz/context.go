package fuzz

import (
	"container/list"
	"sync"
	"sync/atomic"
	"time"
)

// fuzzerContext is a global context during whole fuzzing process.
var fuzzerContext *Context = NewContext()

// Context record all necessary information for help fuzzer to prioritize input and record process.
type Context struct {
	runTaskCh        chan *RunTask // task for worker to run
	fuzzingQueue     *list.List
	fqLock           sync.RWMutex // lock for fuzzingQueue
	mainRecord       *Record
	allRecordHashMap map[string]struct{}

	// A map from bug ID to stdout file contains that bug
	allBugID2Fp  map[string]*BugMetrics
	bugID2FpLock sync.RWMutex

	// Metrics
	numOfBugsFound      uint64
	numOfRuns           uint64
	numOfFuzzQueryEntry uint64
	numOfTargets        uint64
	targetStages        map[string]*TargetMetrics
	targetStagesLock    sync.RWMutex
	startAt             time.Time

	// timeout counter: src => how many times timeout when running this src
	// if more than 3, drop it in the queue
	timeoutTargets     map[string]uint32
	timeoutTargetsLock sync.RWMutex
}

// NewContext returns a new FuzzerContext
func NewContext() *Context {
	return &Context{
		runTaskCh:        make(chan *RunTask),
		fuzzingQueue:     list.New(),
		mainRecord:       EmptyRecord(),
		allRecordHashMap: make(map[string]struct{}),
		allBugID2Fp:      make(map[string]*BugMetrics),
		targetStages:     make(map[string]*TargetMetrics),
		timeoutTargets:   make(map[string]uint32),
		startAt:          time.Now(),
	}
}

func (c *Context) DequeueQueryEntry() (*QueueEntry, error) {
	c.fqLock.RLock()
	if c.fuzzingQueue.Len() == 0 {
		c.fqLock.RUnlock()
		return nil, nil
	}
	elm := c.fuzzingQueue.Front()
	c.fqLock.RUnlock()

	c.fqLock.Lock()
	entry := c.fuzzingQueue.Remove(elm)
	c.fqLock.Unlock()
	return entry.(*QueueEntry), nil

}
func (c *Context) EnqueueQueryEntry(e *QueueEntry) error {
	c.fqLock.Lock()
	c.fuzzingQueue.PushBack(e)
	c.fqLock.Unlock()
	return nil
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

func (c *Context) UpdateTargetStage(target string, stage FuzzStage) {
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

	if _, exist := track.At[stage]; !exist {
		track.At[stage] = time.Now()
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
