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
	q            *list.List
	lock         sync.RWMutex // lock for Context
	cfg          *Config
	mainRecord   *Record
	recordHashes map[string]struct{}
	autoIncID    uint32
	// Bugs
	bugs2InputID map[string]string

	// Metrics
	numOfBugsFound uint64
	numOfRuns      uint64
	startAt        time.Time

	// timeout counter: src => how many times timeout when running this src
	// if more than 3, drop it in the queue
	timeoutTargets map[string]uint32
}

// NewContext returns a new FuzzerContext
func NewContext(execs []exec.Executable, cfg *Config) *Context {
	q := list.New()
	for e := range execs {
		q.PushBack(e)
	}
	return &Context{
		q:              q,
		cfg:            cfg,
		mainRecord:     EmptyRecord(),
		recordHashes:   make(map[string]struct{}),
		timeoutTargets: make(map[string]uint32),
		startAt:        time.Now(),
	}
}

func (c *Context) IncNumOfRun() {
	atomic.AddUint64(&c.numOfRuns, 1)
}

func (c *Context) IncNumOfBugsFound(num uint64) {
	atomic.AddUint64(&c.numOfBugsFound, num)
}

func (c *Context) HasBugID(id string) bool {
	c.lock.RLock()
	_, exists := c.bugs2InputID[id]
	c.lock.RUnlock()
	return exists
}

func (c *Context) AddBugID(bugID string, inputID string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.bugs2InputID[bugID] = inputID
}

func (c *Context) GetAutoIncGlobalID() uint32 {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.autoIncID += 1
	return c.autoIncID
}

// func (c *Context) UpdateTargetMaxCaseCov(target string, caseCov float32) {
// 	c.lock.Lock()
// 	defer c.lock.Unlock()

// 	var track *TargetMetrics
// 	var exist bool
// 	if track, exist = c.targetStages[target]; !exist {
// 		track = &TargetMetrics{
// 			At: make(map[FuzzStage]time.Time),
// 		}
// 		c.targetStages[target] = track
// 	}

// 	if caseCov > track.MaxCaseCov {
// 		track.MaxCaseCov = caseCov
// 	}
// }

func (c *Context) RecordTargetTimeoutOnce(target string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.timeoutTargets[target] += 1
}
