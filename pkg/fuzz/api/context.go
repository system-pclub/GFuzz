package api

import (
	"gfuzz/pkg/fuzz/config"
	"gfuzz/pkg/fuzz/gexecfuzz"
	"gfuzz/pkg/gexec"
	"sync"
	"sync/atomic"
	"time"
)

// Context record all necessary information for help fuzzer to prioritize input and record process.
type Context struct {
	ExecInputCh chan *Input
	fm          map[string]*gexecfuzz.GExecFuzz // map of gexec ID and GExecFuzz

	lock sync.RWMutex   // lock for Context
	Cfg  *config.Config // fuzz configuration

	Interests    InterestList // interested inputs
	recordHashes map[string]struct{}
	// autoIncID is for unique execution run each time
	autoIncID uint32
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
func NewContext(
	execs []gexec.Executable,
	cfg *config.Config,
) *Context {
	fm := make(map[string]*gexecfuzz.GExecFuzz)

	// Create QueueEntry for each Go executable(gexec)
	for _, e := range execs {
		entry := gexecfuzz.NewGExecFuzz(e)
		fm[e.String()] = entry
	}
	return &Context{
		ExecInputCh:    make(chan *Input),
		fm:             fm,
		Cfg:            cfg,
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

func (c *Context) GetQueueEntryByGExecID(gexecID string) *gexecfuzz.GExecFuzz {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.fm[gexecID]
}

func (c *Context) EachGExecFuzz(f func(*gexecfuzz.GExecFuzz)) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	for _, g := range c.fm {
		f(g)
	}
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
	return atomic.AddUint32(&c.autoIncID, 1)
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
