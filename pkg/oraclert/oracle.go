package oraclert

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	NotePrintInput string = "PrintInput"
	InputFileName  string = "myinput.txt"
	RecordFileName        = "myrecord.txt"
	OutputFileName        = "myoutput.txt"
	ErrFileName           = "myerror.txt"
	RecordSplitter        = "-----"
)

var StrTestPath string
var StrTestMod string
var StrTestName string
var StrTestFile string

const (
	DelayCheckModPerTime int = 0 // Check bugs every DelayCheckMS Milliseconds
	DelayCheckModCount   int = 1 // Check bugs when runtime.EnqueueCheckEntry is called DelayCheckCountMax times
)

// config
var DelayCheckMod int = DelayCheckModPerTime
var DelayCheckMS int = 1000
var DelayCheckCountMax uint32 = 10

type OracleEntry struct {
	WgCheckBug              *sync.WaitGroup
	ChEnforceCheck          chan struct{}
	Uint32DelayCheckCounter uint32
}

func init() {
	runtime.FnPointer2String = StrPointer
}

func StrPointer(v interface{}) string {
	return fmt.Sprintf("%p", v)
}

var BoolOracleStarted bool = false // This variable is used to avoid this problem: a test invokes multiple tests, and so our
// BeforeRun and AfterRun is also invoked multiple times, bringing unexpected problems to fuzzer
var globalEntry *OracleEntry
var afterRunCalled bool
var entryMtx sync.Mutex

func BeforeRun() *OracleEntry {
	StrTestMod = os.Getenv("TestMod")
	switch StrTestMod {
	case "TestOnce": // Run all unit tests once, and print a file containing each test's name, # of select visited
		return BeforeRunTestOnce()
	default: // Normal fuzzing
		return BeforeRunFuzz()
	}
}

func BeforeRunTestOnce() *OracleEntry {
	StrTestPath = os.Getenv("TestPath")
	StrTestName = runtime.MyCaller(1)
	if indexDot := strings.Index(StrTestName, "."); indexDot > -1 {
		StrTestName = StrTestName[indexDot+1:]
	}
	_, StrTestFile, _, _ = runtime.Caller(2)
	runtime.BoolSelectCount = true
	return &OracleEntry{
		WgCheckBug:              &sync.WaitGroup{},
		ChEnforceCheck:          make(chan struct{}),
		Uint32DelayCheckCounter: 0,
	}
}

func BeforeRunFuzz() (result *OracleEntry) {
	entryMtx.Lock()
	defer entryMtx.Unlock()
	if globalEntry != nil {
		println("[oraclert] already started, return existing entry.")
		return globalEntry
	} else {
		BoolOracleStarted = true
		println("[oraclert] started")
	}

	result = &OracleEntry{
		WgCheckBug:              &sync.WaitGroup{},
		ChEnforceCheck:          make(chan struct{}),
		Uint32DelayCheckCounter: 0,
	}
	globalEntry = result
	StrBitGlobalTuple := os.Getenv("BitGlobalTuple")
	if StrBitGlobalTuple == "1" {
		runtime.BoolRecordPerCh = false
	} else {
		runtime.BoolRecordPerCh = true
	}
	StrTestPath = os.Getenv("TestPath")
	//StrTestPath ="/data/ziheng/shared/gotest/gotest/src/gotest/testdata/toyprogram"

	// read input file
	if ortBenchmark {
		runtime.RecordSelectChoice = false
	}
	CheckBugStart(result)
	return
}

// Only enables oracle
func LightBeforeRun() *OracleEntry {
	if BoolOracleStarted {
		return nil
	} else {
		BoolOracleStarted = true
	}
	entry := &OracleEntry{
		WgCheckBug:              &sync.WaitGroup{},
		ChEnforceCheck:          make(chan struct{}),
		Uint32DelayCheckCounter: 0,
	}
	CheckBugStart(entry)
	return entry
}

// Start the endless loop that checks bug. Should be called at the beginning of unit test
func CheckBugStart(entry *OracleEntry) {
	// go CheckBugLate()
	// cancelled CheckBugLate since we are going to ignore timeout result

	if runtime.BoolDelayCheck {
		if DelayCheckMod == DelayCheckModCount {
			runtime.FnCheckCount = DelayCheckCounterFN
			runtime.PtrCheckCounter = &entry.Uint32DelayCheckCounter // TODO: potential data race here
		}
		entry.WgCheckBug.Add(1)
		go CheckBugRun(entry)
	}
}

// An endless loop that checks bug. Exits when the unit test ends
func CheckBugRun(entry *OracleEntry) {
	defer entry.WgCheckBug.Done()

	boolBreakLoop := false
	for {
		switch DelayCheckMod {
		case DelayCheckModPerTime:
			select {
			case <-time.After(time.Millisecond * time.Duration(DelayCheckMS)):
			case <-entry.ChEnforceCheck:
				if runtime.BoolDebug {
					fmt.Printf("Check bugs at the end of unit test\n")
				}
				boolBreakLoop = true
			}
		case DelayCheckModCount:
			if atomic.LoadUint32(&entry.Uint32DelayCheckCounter) >= DelayCheckCountMax {
				atomic.StoreUint32(&entry.Uint32DelayCheckCounter, 0)
				// check
			} else {
				select {
				case <-time.After(time.Millisecond * time.Duration(DelayCheckMS)): // set a timeout to check Uint32DelayCheckCounter again later
					//don't check, go back to see if counter >= DelayCheckCountMax
					continue
				case <-entry.ChEnforceCheck:
					if runtime.BoolDebug {
						fmt.Printf("Check bugs at the end of unit test\n")
					}
					boolBreakLoop = true
					// check
				}
			}
		}

		enqueueAgain := [][]runtime.PrimInfo{}
		for {
			runtime.LockCheckEntry()
			if len(runtime.VecCheckEntry) == 0 {
				runtime.UnlockCheckEntry()
				break
			}
			runtime.UnlockCheckEntry()
			checkEntry := runtime.DequeueCheckEntry()
			if checkEntry == nil {
				continue
			}
			if runtime.BoolDebug {
				print("Dequeueing:")
				for _, C := range checkEntry.CS {
					if ch, ok := C.(*runtime.ChanInfo); ok {
						print("\t", ch.StrDebug)
					}
				}
				println()
			}
			if atomic.LoadUint32(&checkEntry.Uint32NeedCheck) == 1 {
				if runtime.CheckBlockBug(checkEntry.CS) == false { // CS needs to be checked again in the future
					enqueueAgain = append(enqueueAgain, checkEntry.CS)
				}
			}
		}
		for _, CS := range enqueueAgain {
			runtime.EnqueueCheckEntry(CS)
		}
		if boolBreakLoop {
			break
		}
	}
}

func CheckBugLate() {
	time.Sleep(20 * time.Second) // Before the deadline we set for unit test in fuzzer/run.go, check once again

	fmt.Printf("Check bugs after 20 seconds\n")

	for {
		runtime.LockCheckEntry()
		if len(runtime.VecCheckEntry) == 0 {
			runtime.UnlockCheckEntry()
			break
		}
		runtime.UnlockCheckEntry()
		checkEntry := runtime.DequeueCheckEntry()
		if runtime.BoolDebug {
			print("Dequeueing:")
			for _, C := range checkEntry.CS {
				if ch, ok := C.(*runtime.ChanInfo); ok {
					print("\t", ch.StrDebug)
				}
			}
			println()
		}
		if atomic.LoadUint32(&checkEntry.Uint32NeedCheck) == 1 {
			if runtime.CheckBlockBug(checkEntry.CS) == false { // CS needs to be checked again in the future
			}
		}
	}

	if ortBenchmark {
		return
	}
	// print bug info
	str, _ := runtime.CheckBlockEntry()

	// print stdout
	if ortStdoutFile != "" {
		out, err := os.OpenFile(ortStdoutFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			fmt.Println("Failed to create file:", ortStdoutFile, err)
			print(str)
			return
		}
		defer out.Close()

		w := bufio.NewWriter(out)
		defer w.Flush()

		w.WriteString(str)
		w.WriteString(runtime.StrWithdraw)
		w.WriteString("---Stack:\n")
		const size = 64 << 10
		buf := make([]byte, size)
		buf = buf[:runtime.Stack(buf, true)]
		w.Write(buf)
	}

	// print record
	// create output file using runtime's global variable
	err := DumpOracleRtOutput(ortConfig, ortOutputFile)
	if err != nil {
		println("DumpOracleRtOutput", err)
	}
}

// When unit test ends, do all delayed bug detect, and wait for the checking process to end
func CheckBugEnd(entry *OracleEntry) {
	if runtime.BoolDelayCheck {
		runtime.SetCurrentGoCheckBug()
		println("End of unit test. Check bugs")
		close(entry.ChEnforceCheck)
		entry.WgCheckBug.Wait() // let's not use send of channel, to make the code clearer

		str, _ := runtime.CheckBlockEntry()

		if ortBenchmark {
			return
		}

		if str == "" {
			println("\n-----NO BLOCKING\n")
		} else {
			println("\n-----FOUND BLOCKING\n")
		}
		//println(str) ignore this since we have stack trace
		println(runtime.StrWithdraw)

		// print stdout
	}
}

func DelayCheckCounterFN(ptrCounter *uint32) {
	if DelayCheckMod == DelayCheckModCount {
		atomic.AddUint32(ptrCounter, 1) // no need to worry about data race, since runtime.MuCheckEntry is held
	}
}

func AfterRun(entry *OracleEntry) {
	println("[oraclert]: AfterRun")
	switch StrTestMod {
	case "TestOnce": // Run all unit tests once, and print a file containing each test's name, # of select visited
		AfterRunTestOnce(entry)
	default: // Normal fuzzing
		AfterRunFuzz(entry)
	}
}

func AfterRunTestOnce(entry *OracleEntry) {
	strOutputPath := os.Getenv("OutputFullPath")
	out, err := os.OpenFile(strOutputPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println("Failed to create file:", strOutputPath, err)
		return
	}
	defer out.Close()

	w := bufio.NewWriter(out)
	defer w.Flush()

	w.WriteString(StrTestNameAndSelectCount())
}

func StrTestNameAndSelectCount() string {
	return "\n" + StrTestFile + ":" + StrTestName + ":" + strconv.Itoa(int(runtime.ReadSelectCount()))
}

func AfterRunFuzz(entry *OracleEntry) {
	entryMtx.Lock()
	if afterRunCalled {
		println("[oraclert]: called. ignored.")
		entryMtx.Unlock()
		return
	}
	afterRunCalled = true
	println("[oraclert]: AfterRunFuzz")
	entryMtx.Unlock()

	if entry == nil {
		println("[oraclert]: entry is nil. return.")
		return
	}

	err := DumpOracleRtOutput(ortConfig, ortOutputFile)
	if err != nil {
		println("DumpOracleRtOutput", err)
	}

	println("[oraclert]: CheckBugEnd...")

	CheckBugEnd(entry)
	runtime.DumpAllStack()
}

// Only enables oracle
func LightAfterRun(entry *OracleEntry) {
	if entry == nil {
		return
	}
	CheckBugEnd(entry)
}
