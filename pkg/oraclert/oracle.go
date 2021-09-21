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

// GFuzzBenchmark will omit any operations related to fuzzing
// including:
// 1. writing input, record file
var GFuzzBenchmark bool = os.Getenv("GF_BENCHMARK") == "1"

var StrTestPath string
var BoolFirstRun bool = true
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

var BoolOracleStarted bool = false // This variable is used to avoid this problem: a test invokes multiple tests, and so our
// BeforeRun and AfterRun is also invoked multiple times, bringing unexpected problems to fuzzer

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
	if BoolOracleStarted {
		return nil
	} else {
		BoolOracleStarted = true
	}
	var err error
	baseStr := os.Getenv("GF_TIME_DIVIDE")
	if baseStr == "" {
		baseStr = "1"
	}
	time.DurDivideBy, err = strconv.Atoi(baseStr)
	if err != nil {
		fmt.Println("Failed to set time.DurDivideBy. time.DurDivideBy is set to 1. Err:", err)
	}

	result = &OracleEntry{
		WgCheckBug:              &sync.WaitGroup{},
		ChEnforceCheck:          make(chan struct{}),
		Uint32DelayCheckCounter: 0,
	}
	StrBitGlobalTuple := os.Getenv("BitGlobalTuple")
	if StrBitGlobalTuple == "1" {
		runtime.BoolRecordPerCh = false
	} else {
		runtime.BoolRecordPerCh = true
	}
	StrTestPath = os.Getenv("TestPath")
	//StrTestPath ="/data/ziheng/shared/gotest/gotest/src/gotest/testdata/toyprogram"

	// Create an output file and bound os.Stdout to it
	//OpenOutputFile() // No need

	// read input file
	if GFuzzBenchmark {
		BoolFirstRun = false
		runtime.RecordSelectChoice = false
	} else {
		file, err := os.Open(FileNameOfInput())
		if err != nil {
			fmt.Println("Failed to open input file:", FileNameOfInput())
			return
		}
		defer file.Close()

		var text []string

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			text = append(text, scanner.Text())
		}

		if len(text) > 0 && text[0] == NotePrintInput {
			runtime.RecordSelectChoice = true
		} else {
			BoolFirstRun = false
		}

		MapInput = ParseInputStr(text)
		if MapInput == nil {
			fmt.Println("Error when parsing input during text start: MapInput is nil")
		}
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
	go CheckBugLate()
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
	time.Sleep(45 * time.Second) // Before the deadline we set for unit test in fuzzer/run.go, check once again

	fmt.Printf("Check bugs after 45 seconds\n")

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

	if GFuzzBenchmark {
		return
	}
	// print bug info
	str, _ := runtime.CheckBlockEntry()

	// print stdout
	out, err := os.OpenFile(os.Getenv("GF_OUTPUT_FILE"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println("Failed to create file:", os.Getenv("GF_OUTPUT_FILE"), err)
		print(str)
		return
	}
	defer out.Close()

	w := bufio.NewWriter(out)
	defer w.Flush()

	w.WriteString(str)
	w.WriteString(runtime.StrWithdraw)

	// print record
	// create output file using runtime's global variable
	CreateRecordFile()

	// print op-cov
	// dump operation records
	opFile := os.Getenv("GF_OP_COV_FILE")
	if opFile != "" {
		err := dumpOpRecordsToFile(opFile, opRecords)
		if err != nil {
			// print to error
			println(err)
		}
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

		if GFuzzBenchmark {
			return
		}
		// print stdout
		out, err := os.OpenFile(os.Getenv("GF_OUTPUT_FILE"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			fmt.Println("Failed to create file:", os.Getenv("GF_OUTPUT_FILE"), err)
			print(str)
			return
		}
		defer out.Close()

		w := bufio.NewWriter(out)
		defer w.Flush()

		w.WriteString(str)
		w.WriteString(runtime.StrWithdraw)
	}
}

func DelayCheckCounterFN(ptrCounter *uint32) {
	if DelayCheckMod == DelayCheckModCount {
		atomic.AddUint32(ptrCounter, 1) // no need to worry about data race, since runtime.MuCheckEntry is held
	}
}

func AfterRun(entry *OracleEntry) {
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
	if entry == nil {
		return
	}

	// if this is the first run, create input file using runtime's global variable
	if BoolFirstRun {
		CreateInput()
	}

	// create output file using runtime's global variable
	CreateRecordFile()

	CheckBugEnd(entry)

	// dump operation records
	opFile := os.Getenv("GF_OP_COV_FILE")
	if opFile != "" {
		err := dumpOpRecordsToFile(opFile, opRecords)
		if err != nil {
			// print to error
			println(err)
		}
	}
}

// Only enables oracle
func LightAfterRun(entry *OracleEntry) {
	if entry == nil {
		return
	}
	CheckBugEnd(entry)
}
