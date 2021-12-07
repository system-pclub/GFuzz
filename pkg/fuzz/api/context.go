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
	ExecInputCh          chan *Input                     // channel for worker to get execute input
	fm                   map[string]*gexecfuzz.GExecFuzz // map of gexec ID and GExecFuzz
	lock                 sync.RWMutex                    // lock for Context
	Cfg                  *config.Config                  // fuzz configuration
	Interests            InterestList                    // interested inputs
	oracleRtOutputHashes map[string]struct{}             // Hashes for oracle runtime configuration. So that it can check if an oracle runtime configuration has been generated or not.

	autoIncID    uint32            // autoIncID is for unique execution run each time
	bugs2InputID map[string]string // Bugs

	// Metrics
	numOfRuns uint64
	startAt   time.Time

	// timeout counter: src => how many times timeout when running this src
	// Skip handling interest input if its target's timeout is more than three
	timeoutTargets map[string]uint32

	// Global best score
	GlobalBestScore int
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
		ExecInputCh:          make(chan *Input),
		fm:                   fm,
		Cfg:                  cfg,
		oracleRtOutputHashes: make(map[string]struct{}),
		bugs2InputID:         make(map[string]string),
		timeoutTargets:       make(map[string]uint32),
		startAt:              time.Now(),
	}
}

func (c *Context) IncNumOfRun() {
	atomic.AddUint64(&c.numOfRuns, 1)
}

func (c *Context) GetNumOfRuns() uint64 {
	return atomic.LoadUint64(&c.numOfRuns)
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
	defer c.lock.RUnlock()
	_, exists := c.bugs2InputID[id]

	return exists
}

func (c *Context) AddBugID(bugID string, inputID string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.bugs2InputID[bugID] = inputID
}

func (c *Context) GetNumOfBugs() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return len(c.bugs2InputID)
}

func (c *Context) GetDuration() time.Duration {
	return time.Now().Sub(c.startAt)
}

func (c *Context) GetAutoIncGlobalID() uint32 {
	return atomic.AddUint32(&c.autoIncID, 1)
}

func (c *Context) RecordTargetTimeoutOnce(target string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.timeoutTargets[target] += 1
}

// UpdateOrtOutputHash will update hash record if not exist. It returns true if it is successfully updated, false otherwise
func (c *Context) UpdateOrtOutputHash(h string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, exist := c.oracleRtOutputHashes[h]; !exist {
		c.oracleRtOutputHashes[h] = struct{}{}
		return true
	}
	return false
}
