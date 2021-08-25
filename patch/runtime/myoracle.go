package runtime

import "sync/atomic"

func init() {
	MapChToChanInfo = make(map[interface{}]PrimInfo)
}

var GlobalEnableOracle = true

// during benchmark, we don't need to print bugs to stdout
var BoolReportBug = gogetenv("GF_BENCHMARK") != "1"
var BoolDelayCheck = true

var MuReportBug mutex

type PrimInfo interface {
	Lock()
	Unlock()
	MapRef() map[*GoInfo]struct{}
	AddGoroutine(*GoInfo)
	RemoveGoroutine(*GoInfo)
	StringDebug() string
	SetMonitor(uint32)
	LoadMonitor() uint32
}

//// Part 1.1: data struct for each channel

// ChanInfo is 1-to-1 with every channel. It tracks a list of goroutines that hold the reference to the channel
type ChanInfo struct {
	Chan            *hchan               // Stores the channel. Can be used as ID of channel
	IntBuffer       int                  // The buffer capability of channel. 0 if channel is unbuffered
	MapRefGoroutine map[*GoInfo]struct{} // Stores all goroutines that still hold reference to this channel
	StrDebug        string
	OKToCheck       bool  // Disable oracle for not instrumented channels
	BoolInSDK       bool  // Disable oracle for channels in SDK
	IntFlagFoundBug int32 // Use atomic int32 operations to mark if a bug is reported
	Mu              mutex
	SpecialFlag     int8
	Uint32Monitor   uint32 // Default: 0. When a bug is based on the assumption that this primitive won't execute again, set to 1.
	// When it is 1 and still executed, print information to withdraw bugs containing it
}

const (
	TimeTicker int8 = 1
)

type OpType uint8

const (
	ChSend   OpType = 1
	ChRecv   OpType = 2
	ChClose  OpType = 3
	ChSelect OpType = 4
)

var MapChToChanInfo map[interface{}]PrimInfo
var MuMapChToChanInfo mutex

//var DefaultCaseChanInfo = &ChanInfo{}

var strSDKPath string = gogetenv("GOROOT")

func PrintSDK() {
	println(strSDKPath)
}

// Initialize a new ChanInfo with a given channel
func NewChanInfo(ch *hchan) *ChanInfo {
	_, strFile, intLine, _ := Caller(2)
	strLoc := strFile + ":" + Itoa(intLine)
	newChInfo := &ChanInfo{
		Chan:            ch,
		IntBuffer:       int(ch.dataqsiz),
		MapRefGoroutine: make(map[*GoInfo]struct{}),
		StrDebug:        strLoc,
		OKToCheck:       false,
		BoolInSDK:       Index(strLoc, strSDKPath) < 0,
		IntFlagFoundBug: 0,
		SpecialFlag:     0,
	}
	if BoolPrintDebugInfo {
		println("===Debug Info:")
		println("\tMake of a new channel. The creation site is:", strLoc)
		println("\tSDK path is:", strSDKPath, "\tBoolMakeNotInSDK is:", newChInfo.BoolInSDK)
	}
	AddRefGoroutine(newChInfo, CurrentGoInfo())

	// If this channel is in some special API, set special flag
	//creationFunc := MyCaller(1)
	if Index(strLoc, "time") >= 0 {
		newChInfo.SpecialFlag = TimeTicker
	}

	return newChInfo
}

func (chInfo *ChanInfo) StringDebug() string {
	if chInfo == nil {
		return ""
	}
	return chInfo.StrDebug
}

func (chInfo *ChanInfo) SetMonitor(i uint32) {
	atomic.StoreUint32(&chInfo.Uint32Monitor, i)
}

func (chInfo *ChanInfo) LoadMonitor() uint32 {
	return atomic.LoadUint32(&chInfo.Uint32Monitor)
}

func okToCheck(c *hchan) bool {
	if c.chInfo != nil {
		switch c.chInfo.SpecialFlag {
		case TimeTicker:
			return false
		}
	}
	return true
}

func (chInfo *ChanInfo) Lock() {
	if chInfo == nil {
		return
	}
	lock(&chInfo.Mu)
}

func (chInfo *ChanInfo) Unlock() {
	if chInfo == nil {
		return
	}
	unlock(&chInfo.Mu)
}

// Must be called with chInfo.Mu locked
func (chInfo *ChanInfo) MapRef() map[*GoInfo]struct{} {
	if chInfo == nil {
		return make(map[*GoInfo]struct{})
	}
	return chInfo.MapRefGoroutine
}

// FindChanInfo can retrieve a initialized ChanInfo for a given channel
func FindChanInfo(ch interface{}) *ChanInfo {
	lock(&MuMapChToChanInfo)
	chInfo := MapChToChanInfo[ch]
	unlock(&MuMapChToChanInfo)
	if chInfo == nil {
		return nil
	} else {
		return chInfo.(*ChanInfo)
	}
}

func LinkChToLastChanInfo(ch interface{}) {
	lock(&MuMapChToChanInfo)
	primInfo := LoadLastPrimInfo()
	MapChToChanInfo[ch] = primInfo
	if chInfo, ok := primInfo.(*ChanInfo); ok {
		chInfo.OKToCheck = true
	}
	unlock(&MuMapChToChanInfo)
}

// After the creation of a new channel, or at the head of a goroutine that holds a reference to a channel,
// or whenever a goroutine obtains a reference to a channel, call this function
// AddRefGoroutine links a channel with a goroutine, meaning the goroutine holds the reference to the channel
func AddRefGoroutine(chInfo PrimInfo, goInfo *GoInfo) {
	if chInfo == nil || goInfo == nil {
		return
	}
	chInfo.AddGoroutine(goInfo)
	goInfo.AddPrime(chInfo)
}

func RemoveRefGoroutine(chInfo PrimInfo, goInfo *GoInfo) {
	if chInfo == nil || goInfo == nil {
		return
	}
	chInfo.RemoveGoroutine(goInfo)
	goInfo.RemovePrime(chInfo)
}

// When this goroutine is checking bug, set goInfo.BitCheckBugAtEnd to be 1
func SetCurrentGoCheckBug() {
	getg().goInfo.SetCheckBug()
}

func (goInfo *GoInfo) SetCheckBug() {
	atomic.StoreUint32(&goInfo.BitCheckBugAtEnd, 1)
}

func (goInfo *GoInfo) SetNotCheckBug() {
	atomic.StoreUint32(&goInfo.BitCheckBugAtEnd, 0)
}

// This means the goroutine mapped with goInfo holds the reference to chInfo.Chan
func (chInfo *ChanInfo) AddGoroutine(goInfo *GoInfo) {
	if chInfo == nil {
		return
	}
	chInfo.Lock()
	if chInfo.MapRefGoroutine == nil {
		chInfo.Unlock()
		return
	}
	chInfo.MapRefGoroutine[goInfo] = struct{}{}
	chInfo.Unlock()
}

func (chInfo *ChanInfo) RemoveGoroutine(goInfo *GoInfo) {
	if chInfo == nil {
		return
	}
	chInfo.Lock()
	if chInfo.MapRefGoroutine == nil {
		chInfo.Unlock()
		return
	}
	delete(chInfo.MapRefGoroutine, goInfo)
	chInfo.Unlock()
}

// Only when BoolDelayCheck is true, this struct is used
// CheckEntry contains information needed for a CheckBlockBug
type CheckEntry struct {
	CS              []PrimInfo
	Uint32NeedCheck uint32 // if 0, delete this CheckEntry; if 1, check this CheckEntry
}

var VecCheckEntry []*CheckEntry
var MuCheckEntry mutex
var FnCheckCount = func(*uint32) {} // this is defined in gooracle/gooracle.go
var PtrCheckCounter *uint32

func LockCheckEntry() {
	lock(&MuCheckEntry)
}

func UnlockCheckEntry() {
	unlock(&MuCheckEntry)
}

func DequeueCheckEntry() *CheckEntry {
	lock(&MuCheckEntry)
	if len(VecCheckEntry) == 0 {
		unlock(&MuCheckEntry)
		return nil
	} else {
		result := VecCheckEntry[0]
		VecCheckEntry = VecCheckEntry[1:]
		unlock(&MuCheckEntry)
		return result
	}
}

func EnqueueCheckEntry(CS []PrimInfo) *CheckEntry {
	lock(&MuCheckEntry)

	if len(CS) == 1 {
		if CS[0].StringDebug() == "/data/ziheng/shared/gotest/stubs/etcd/pkg/mod/github.com/prometheus/client_golang@v1.0.0/prometheus/registry.go:266" {
			print()
		}
	}

	FnCheckCount(PtrCheckCounter)
	newCheckEntry := &CheckEntry{
		CS:              CS,
		Uint32NeedCheck: 1,
	}
	for _, entry := range VecCheckEntry {
		if BoolCheckEntryEqual(entry, newCheckEntry) {
			unlock(&MuCheckEntry)
			return nil // It's OK to return nil
		}
	}
	if BoolDebug {
		if len(CS) == 1 {
			if CS[0].StringDebug() == "/data/ziheng/shared/gotest/stubs/etcd/pkg/mod/github.com/prometheus/client_golang@v1.0.0/prometheus/registry.go:266" {
				print()
			}
		}
		print("Enqueueing:")
		for _, c := range CS {
			print("\t", c.StringDebug())
		}
		println()
	}

	VecCheckEntry = append(VecCheckEntry, newCheckEntry)
	unlock(&MuCheckEntry)
	return newCheckEntry
}

func BoolCheckEntryEqual(a, b *CheckEntry) bool {
	if len(a.CS) != len(b.CS) {
		return false
	}
	for _, primInfo1 := range a.CS {
		boolFound := false
		for _, primInfo2 := range b.CS {
			if primInfo2 == primInfo1 {
				boolFound = true
				break
			}
		}
		if boolFound == false {
			return false
		}
	}
	return true
}

// A blocking bug is detected, if all goroutines that hold the reference to a channel are blocked at an operation of the channel
// finished is true when we are sure that CS doesn't need to be checked again
func CheckBlockBug(CS []PrimInfo) (finished bool) {
	mapCS := make(map[PrimInfo]struct{})
	mapGS := make(map[*GoInfo]struct{}) // all goroutines that hold reference to primitives in mapCS
	finished = false

	if BoolDebug {
		print("Checking primtives:")
		for _, chI := range CS {
			print("\t", chI.StringDebug())
		}
		println()
	}

	for _, primI := range CS {
		if primI == (*ChanInfo)(nil) {
			continue
		}
		if chI, ok := primI.(*ChanInfo); ok {
			if chI.OKToCheck == false {
				return true
			}
		}
		if Index(primI.StringDebug(), strSDKPath) >= 0 {
			if BoolDebug {
				println("Abort checking because this prim is in SDK:", primI.StringDebug())
			}
			finished = true
			return
		}
		primI.Lock()
		for goInfo, _ := range primI.MapRef() {
			mapGS[goInfo] = struct{}{}
		}
		primI.Unlock()
		mapCS[primI] = struct{}{}
	}

	boolAtLeastOneBlocking := false
loopGS:
	if BoolDebug {
		println("Going through mapGS whose length is:", len(mapGS))
	}
	for goInfo, _ := range mapGS {
		if atomic.LoadUint32(&goInfo.BitCheckBugAtEnd) == 1 { // The goroutine is checking bug at the end of unit test
			if BoolDebug {
				println("\tGoID", goInfo.G.goid, "is checking bug at the end of unit test")
			}
			continue
		}
		lock(&goInfo.Mu)
		if len(goInfo.VecBlockInfo) == 0 { // The goroutine is executing non-blocking operations
			if BoolDebug {
				println("\tGoID", goInfo.G.goid, "is executing non-blocking operations")
				println("Aborting checking")
			}
			unlock(&goInfo.Mu)
			return
		}

		boolAtLeastOneBlocking = true

		if BoolDebug {
			println("\tGoID", goInfo.G.goid, "is executing blocking operations")
		}

		for _, blockInfo := range goInfo.VecBlockInfo { // if it is blocked at select, VecBlockInfo contains multiple primitives

			primI := blockInfo.Prim
			if _, exist := mapCS[primI]; !exist {
				if BoolDebug {
					println("\t\tNot existing prim in CS:", blockInfo.Prim.StringDebug())
				}
				mapCS[primI] = struct{}{} // update CS
				primI.Lock()
				for goInfo, _ := range primI.MapRef() { // update GS
					mapGS[goInfo] = struct{}{}
				}
				primI.Unlock()
				unlock(&goInfo.Mu)
				if BoolDebug {
					println("Goto mapGS loop again")
				}
				goto loopGS // since mapGS is updated, we should run this loop once again
			} else {
				if BoolDebug {
					println("\t\tExisting prim in CS:", blockInfo.Prim.StringDebug())
				}
			}
		}
		unlock(&goInfo.Mu)
	}

	if boolAtLeastOneBlocking {
		ReportBug(mapCS)
		finished = true
	}

	return
}

//func (chInfo *ChanInfo) CheckBlockBug() {
//	if atomic.LoadInt32(&chInfo.intFlagFoundBug) != 0 {
//		return
//	}
//
//	if chInfo.intBuffer == 0 {
//		countRefGo := 0 // Number of goroutines that hold the reference to the channel
//		countBlockAtThisChanGo := 0 // Number of goroutines that are blocked at an operation of this channel
//		f := func(key interface{}, value interface{}) bool {
//			goInfo, _ := key.(*GoInfo)
//
//			boolIsBlock, _ := goInfo.IsBlockAtGivenChan(chInfo)
//			if boolIsBlock {
//				countBlockAtThisChanGo++
//			}
//			countRefGo++
//			return true // continue Range
//		}
//		chInfo.mapRefGoroutine.Range(f)
//
//		if countRefGo == countBlockAtThisChanGo {
//			if countRefGo == 0 {
//				// debug code
//				countRefGo2 := 0 // Number of goroutines that hold the reference to the channel
//				countBlockAtThisChanGo2 := 0 // Number of goroutines that are blocked at an operation of this channel
//				f := func(key interface{}, value interface{}) bool {
//					goInfo, _ := key.(*GoInfo)
//
//					boolIsBlock, _ := goInfo.IsBlockAtGivenChan(chInfo)
//					if boolIsBlock {
//						countBlockAtThisChanGo2++
//					}
//					countRefGo2++
//					return true // continue Range
//				}
//				chInfo.mapRefGoroutine.Range(f)
//				fmt.Print()
//
//				return
//			}
//			ReportBug(chInfo)
//			atomic.AddInt32(&chInfo.intFlagFoundBug, 1)
//		}
//
//	} else { // Buffered channel
//		if reflect.ValueOf(chInfo.Chan).Len() == chInfo.intBuffer { // Buffer is full
//			// Check if all ref goroutines are blocked at send
//			boolAllBlockAtSend := true
//			countRefGo := 0
//			countBlockAtThisChanGo := 0
//			f := func(key interface{}, value interface{}) bool {
//				goInfo, _ := key.(*GoInfo)
//
//				boolIsBlock, strOp := goInfo.IsBlockAtGivenChan(chInfo)
//				if boolIsBlock {
//					countBlockAtThisChanGo++
//				}
//				if strOp != Send {
//					boolAllBlockAtSend = false
//				}
//				countRefGo++
//				return true // continue Range
//			}
//			chInfo.mapRefGoroutine.Range(f)
//
//			if countRefGo == countBlockAtThisChanGo && boolAllBlockAtSend {
//				ReportBug(chInfo)
//				atomic.AddInt32(&chInfo.intFlagFoundBug, 1)
//			}
//
//		} else if reflect.ValueOf(chInfo.Chan).Len() == 0 { // Buffer is empty
//			// Check if all ref goroutines are blocked at receive
//			boolAllBlockAtRecv := true
//			countRefGo := 0
//			countBlockAtThisChanGo := 0
//			f := func(key interface{}, value interface{}) bool {
//				goInfo, _ := key.(*GoInfo)
//
//				boolIsBlock, strOp := goInfo.IsBlockAtGivenChan(chInfo)
//				if boolIsBlock {
//					countBlockAtThisChanGo++
//				}
//				if strOp != Recv {
//					boolAllBlockAtRecv = false
//				}
//				countRefGo++
//				return true // continue Range
//			}
//			chInfo.mapRefGoroutine.Range(f)
//
//			if countRefGo == countBlockAtThisChanGo && boolAllBlockAtRecv {
//				ReportBug(chInfo)
//				atomic.AddInt32(&chInfo.intFlagFoundBug, 1)
//			}
//
//		} else { // Buffer is not full or empty. Then it is not possible to block
//			// do nothing
//		}
//	}
//}

func ReportBug(mapCS map[PrimInfo]struct{}) {

	//for chInfo, _ := range mapCS {
	//	atomic.AddInt32(&chInfo.IntFlagFoundBug, 1)
	//}
	//return
	if BoolReportBug == false {
		return
	}
	lock(&MuReportBug)
	print("-----New Blocking Bug:\n")
	const size = 64 << 10
	buf := make([]byte, size)
	buf = buf[:Stack(buf, false)]
	print("---Primitive location:\n")
	for primInfo, _ := range mapCS {
		print(primInfo.StringDebug() + "\n")
		primInfo.SetMonitor(1)
	}
	print("---Primitive pointer:\n")
	for primInfo, _ := range mapCS {
		print(FnPointer2String(primInfo) + "\n")
	}
	print("---Stack:\n", string(buf), "\n")
	print("-----End Bug\n")
	unlock(&MuReportBug)
}

func ReportNonBlockingBug() {
	print("-----New NonBlocking Bug:\n")
	const size = 64 << 10
	buf := make([]byte, size)
	buf = buf[:Stack(buf, false)]
	print("---Stack:\n", string(buf), "\n")
	print("-----End Bug\n")
}

// Part 1.2 Data structure for waitgroup

// WgInfo is 1-to-1 with every WaitGroup.
type WgInfo struct {
	WgCounter       uint32
	MapRefGoroutine map[*GoInfo]struct{}
	StrDebug        string
	EnableOracle    bool  // Disable oracle for channels in SDK
	IntFlagFoundBug int32 // Use atomic int32 operations to mark if a bug is reported
	Mu              mutex // Protects MapRefGoroutine
	Uint32Monitor   uint32
}

func NewWgInfo() *WgInfo {
	_, strFile, intLine, _ := Caller(2)
	strLoc := strFile + ":" + Itoa(intLine)
	wg := &WgInfo{
		WgCounter:       0,
		MapRefGoroutine: make(map[*GoInfo]struct{}),
		StrDebug:        strLoc,
		EnableOracle:    Index(strLoc, strSDKPath) < 0,
		IntFlagFoundBug: 0,
		Mu:              mutex{},
	}
	return wg
}

func (w *WgInfo) StringDebug() string {
	if w == nil {
		return ""
	}
	return w.StrDebug
}

func (w *WgInfo) SetMonitor(i uint32) {
	atomic.StoreUint32(&w.Uint32Monitor, i)
}

func (w *WgInfo) LoadMonitor() uint32 {
	return atomic.LoadUint32(&w.Uint32Monitor)
}

// FindChanInfo can retrieve a initialized ChanInfo for a given channel
func FindWgInfo(wg interface{}) *WgInfo {
	lock(&MuMapChToChanInfo)
	wgInfo := MapChToChanInfo[wg]
	unlock(&MuMapChToChanInfo)
	return wgInfo.(*WgInfo)
}

func LinkWgToLastWgInfo(wg interface{}) {
	lock(&MuMapChToChanInfo)
	MapChToChanInfo[wg] = LoadLastPrimInfo()
	unlock(&MuMapChToChanInfo)
}

func (w *WgInfo) Lock() {
	lock(&w.Mu)
}

func (w *WgInfo) Unlock() {
	unlock(&w.Mu)
}

// Must be called with lock
func (w *WgInfo) MapRef() map[*GoInfo]struct{} {
	return w.MapRefGoroutine
}

// This means the goroutine mapped with goInfo holds the reference to chInfo.Chan
// Must be called when chInfo.Chan.lock is held
func (w *WgInfo) AddGoroutine(goInfo *GoInfo) {
	w.MapRefGoroutine[goInfo] = struct{}{}
}

// Must be called when chInfo.Chan.lock is held
func (w *WgInfo) RemoveGoroutine(goInfo *GoInfo) {
	delete(w.MapRefGoroutine, goInfo)
}

func (w *WgInfo) IamBug() {

}

// Part 1.3 Data structure for mutex

// MuInfo is 1-to-1 with every sync.Mutex.
type MuInfo struct {
	MapRefGoroutine map[*GoInfo]struct{}
	StrDebug        string
	EnableOracle    bool  // Disable oracle for channels in SDK
	IntFlagFoundBug int32 // Use atomic int32 operations to mark if a bug is reported
	Mu              mutex // Protects MapRefGoroutine
	Uint32Monitor   uint32
}

func NewMuInfo() *MuInfo {
	_, strFile, intLine, _ := Caller(2)
	strLoc := strFile + ":" + Itoa(intLine)
	mu := &MuInfo{
		MapRefGoroutine: make(map[*GoInfo]struct{}),
		StrDebug:        strLoc,
		EnableOracle:    Index(strLoc, strSDKPath) < 0,
		IntFlagFoundBug: 0,
		Mu:              mutex{},
	}
	return mu
}

func (m *MuInfo) StringDebug() string {
	if m == nil {
		return ""
	}
	return m.StrDebug
}

func (m *MuInfo) SetMonitor(i uint32) {
	atomic.StoreUint32(&m.Uint32Monitor, i)
}

func (m *MuInfo) LoadMonitor() uint32 {
	return atomic.LoadUint32(&m.Uint32Monitor)
}

// FindChanInfo can retrieve a initialized ChanInfo for a given channel
func FindMuInfo(mu interface{}) *MuInfo {
	lock(&MuMapChToChanInfo)
	muInfo := MapChToChanInfo[mu]
	unlock(&MuMapChToChanInfo)
	return muInfo.(*MuInfo)
}

func LinkMuToLastMuInfo(mu interface{}) {
	lock(&MuMapChToChanInfo)
	MapChToChanInfo[mu] = LoadLastPrimInfo()
	unlock(&MuMapChToChanInfo)
}

func (mu *MuInfo) Lock() {
	lock(&mu.Mu)
}

func (mu *MuInfo) Unlock() {
	unlock(&mu.Mu)
}

// Must be called with lock
func (mu *MuInfo) MapRef() map[*GoInfo]struct{} {
	return mu.MapRefGoroutine
}

// This means the goroutine mapped with goInfo holds the reference to chInfo.Chan
// Must be called when chInfo.Chan.lock is held
func (mu *MuInfo) AddGoroutine(goInfo *GoInfo) {
	mu.MapRefGoroutine[goInfo] = struct{}{}
}

// Must be called when chInfo.Chan.lock is held
func (mu *MuInfo) RemoveGoroutine(goInfo *GoInfo) {
	delete(mu.MapRefGoroutine, goInfo)
}

// Part 1.4 Data structure for rwmutex

// RWMuInfo is 1-to-1 with every sync.RWMutex.
type RWMuInfo struct {
	MapRefGoroutine map[*GoInfo]struct{}
	StrDebug        string
	EnableOracle    bool  // Disable oracle for channels in SDK
	IntFlagFoundBug int32 // Use atomic int32 operations to mark if a bug is reported
	Mu              mutex // Protects MapRefGoroutine
	Uint32Monitor   uint32
}

func NewRWMuInfo() *RWMuInfo {
	_, strFile, intLine, _ := Caller(2)
	strLoc := strFile + ":" + Itoa(intLine)
	mu := &RWMuInfo{
		MapRefGoroutine: make(map[*GoInfo]struct{}),
		StrDebug:        strLoc,
		EnableOracle:    Index(strLoc, strSDKPath) < 0,
		IntFlagFoundBug: 0,
		Mu:              mutex{},
	}
	return mu
}

func (m *RWMuInfo) StringDebug() string {
	if m == nil {
		return ""
	}
	return m.StrDebug
}

func (m *RWMuInfo) SetMonitor(i uint32) {
	atomic.StoreUint32(&m.Uint32Monitor, i)
}

func (m *RWMuInfo) LoadMonitor() uint32 {
	return atomic.LoadUint32(&m.Uint32Monitor)
}

// FindChanInfo can retrieve a initialized ChanInfo for a given channel
func FindRWMuInfo(rwmu interface{}) *RWMuInfo {
	lock(&MuMapChToChanInfo)
	muInfo := MapChToChanInfo[rwmu]
	unlock(&MuMapChToChanInfo)
	return muInfo.(*RWMuInfo)
}

func LinkRWMuToLastRWMuInfo(rwmu interface{}) {
	lock(&MuMapChToChanInfo)
	MapChToChanInfo[rwmu] = LoadLastPrimInfo()
	unlock(&MuMapChToChanInfo)
}

func (mu *RWMuInfo) Lock() {
	lock(&mu.Mu)
}

func (mu *RWMuInfo) Unlock() {
	unlock(&mu.Mu)
}

// Must be called with lock
func (mu *RWMuInfo) MapRef() map[*GoInfo]struct{} {
	return mu.MapRefGoroutine
}

// This means the goroutine mapped with goInfo holds the reference to chInfo.Chan
// Must be called when chInfo.Chan.lock is held
func (mu *RWMuInfo) AddGoroutine(goInfo *GoInfo) {
	mu.MapRefGoroutine[goInfo] = struct{}{}
}

// Must be called when chInfo.Chan.lock is held
func (mu *RWMuInfo) RemoveGoroutine(goInfo *GoInfo) {
	delete(mu.MapRefGoroutine, goInfo)
}

// Part 1.5 Data structure for conditional variable

// CondInfo is 1-to-1 with every sync.Cond.
type CondInfo struct {
	MapRefGoroutine map[*GoInfo]struct{}
	StrDebug        string
	EnableOracle    bool  // Disable oracle for channels in SDK
	IntFlagFoundBug int32 // Use atomic int32 operations to mark if a bug is reported
	Mu              mutex // Protects MapRefGoroutine
	Uint32Monitor   uint32
}

func NewCondInfo() *CondInfo {
	_, strFile, intLine, _ := Caller(2)
	strLoc := strFile + ":" + Itoa(intLine)
	cond := &CondInfo{
		MapRefGoroutine: make(map[*GoInfo]struct{}),
		StrDebug:        strLoc,
		EnableOracle:    Index(strLoc, strSDKPath) < 0,
		IntFlagFoundBug: 0,
		Mu:              mutex{},
	}
	return cond
}

func (c *CondInfo) StringDebug() string {
	if c == nil {
		return ""
	}
	return c.StrDebug
}

func (c *CondInfo) SetMonitor(i uint32) {
	atomic.StoreUint32(&c.Uint32Monitor, i)
}

func (c *CondInfo) LoadMonitor() uint32 {
	return atomic.LoadUint32(&c.Uint32Monitor)
}

// FindChanInfo can retrieve a initialized ChanInfo for a given channel
func FindCondInfo(cond interface{}) *CondInfo {
	lock(&MuMapChToChanInfo)
	condInfo := MapChToChanInfo[cond]
	unlock(&MuMapChToChanInfo)
	return condInfo.(*CondInfo)
}

func LinkCondToLastCondInfo(cond interface{}) {
	lock(&MuMapChToChanInfo)
	MapChToChanInfo[cond] = LoadLastPrimInfo()
	unlock(&MuMapChToChanInfo)
}

func (cond *CondInfo) Lock() {
	lock(&cond.Mu)
}

func (cond *CondInfo) Unlock() {
	unlock(&cond.Mu)
}

// Must be called with lock
func (cond *CondInfo) MapRef() map[*GoInfo]struct{} {
	return cond.MapRefGoroutine
}

// This means the goroutine mapped with goInfo holds the reference to chInfo.Chan
// Must be called when chInfo.Chan.lock is held
func (cond *CondInfo) AddGoroutine(goInfo *GoInfo) {
	cond.MapRefGoroutine[goInfo] = struct{}{}
}

// Must be called when chInfo.Chan.lock is held
func (cond *CondInfo) RemoveGoroutine(goInfo *GoInfo) {
	delete(cond.MapRefGoroutine, goInfo)
}

//// Part 2.1: data struct for each goroutine

// GoInfo is 1-to-1 with each goroutine.
// Go language doesn't allow us to acquire the ID of a goroutine, because they want goroutines to be anonymous.
// Normally, Go programmers use runtime.Stack() to print all IDs of all goroutines, but this function is very inefficient
//, since it calls stopTheWorld()
// Currently we use a global atomic int64 to differentiate each goroutine, and a variable currentGo to represent each goroutine
// This is not a good practice because the goroutine need to pass currentGo to its every callee
type GoInfo struct {
	G            *g
	VecBlockInfo []BlockInfo // Nil when normally running. When blocked at an operation of ChanInfo, store
	// one ChanInfo and the operation. When blocked at select, store multiple ChanInfo and
	// operation. Default in select is also also stored in map, which is DefaultCaseChanInfo
	BitCheckBugAtEnd uint32                // 0 when normally running. 1 when this goroutine is checking bug.
	MapPrimeInfo     map[PrimInfo]struct{} // Stores all channels that this goroutine still hold reference to
	Mu               mutex                 // protects VecBlockInfo and MapPrimeInfo
}

type BlockInfo struct {
	Prim  PrimInfo
	StrOp string
}

const (
	Send   = "Send"
	Recv   = "Recv"
	Close  = "Close"
	Select = "Select"

	MuLock   = "MuLock"
	MuUnlock = "MuUnlock"

	WgWait = "WgWait"

	CdWait      = "CdWait"
	CdSignal    = "CdSignal"
	CdBroadcast = "CdBroadcast"
)

// Initialize a GoInfo
func NewGoInfo(goroutine *g) *GoInfo {
	newGoInfo := &GoInfo{
		G:            goroutine,
		VecBlockInfo: []BlockInfo{},
		MapPrimeInfo: make(map[PrimInfo]struct{}),
	}
	return newGoInfo
}

func CurrentGoInfo() *GoInfo {
	return getg().goInfo
}

func StoreLastPrimInfo(chInfo PrimInfo) {
	getg().lastPrimInfo = chInfo
}

func LoadLastPrimInfo() PrimInfo {
	return getg().lastPrimInfo
}

func CurrentGoID() int64 {
	return getg().goid
}

// This means the goroutine mapped with goInfo holds the reference to chInfo.Chan
func (goInfo *GoInfo) AddPrime(chInfo PrimInfo) {
	if goInfo.MapPrimeInfo == nil {
		goInfo.MapPrimeInfo = make(map[PrimInfo]struct{})
	}
	goInfo.MapPrimeInfo[chInfo] = struct{}{}
}

func (goInfo *GoInfo) RemovePrime(chInfo PrimInfo) {
	if goInfo.MapPrimeInfo != nil {
		delete(goInfo.MapPrimeInfo, chInfo)
	}
}

func CurrentGoAddCh(ch interface{}) {
	lock(&MuMapChToChanInfo)
	chInfo, exist := MapChToChanInfo[ch]
	unlock(&MuMapChToChanInfo)
	if !exist {
		return
	}
	AddRefGoroutine(chInfo, CurrentGoInfo())
}

// RemoveRef should be called at the end of every goroutine. It will remove goInfo from the reference list of every
// channel it holds the reference to
func (goInfo *GoInfo) RemoveAllRef() {

	if goInfo.MapPrimeInfo == nil {
		return
	}
	for chInfo, _ := range goInfo.MapPrimeInfo {
		RemoveRefGoroutine(chInfo, goInfo)
		if chInfo == nil {
			continue
		}
		CS := []PrimInfo{chInfo}
		if BoolDelayCheck {
			EnqueueCheckEntry(CS)
		} else {
			CheckBlockBug(CS)
		}
	}
}

// SetBlockAt should be called before each channel operation, meaning the current goroutine is about to execute that operation
// Note that we check bug in this function, because it's possible for the goroutine to be blocked forever if it execute that operation
// For example, a channel with no buffer is held by a parent and a child.
//              The parent has already exited, but the child is now about to send to that channel.
//              Then now is our only chance to detect this bug, so we call CheckBlockBug()
func (goInfo *GoInfo) SetBlockAt(prim PrimInfo, strOp string) {
	goInfo.VecBlockInfo = append(goInfo.VecBlockInfo, BlockInfo{
		Prim:  prim,
		StrOp: strOp,
	})
}

// WithdrawBlock should be called after each channel operation, meaning the current goroutine finished execution that operation
// If the operation is select, remember to call this function right after each case of the select
func (goInfo *GoInfo) WithdrawBlock(checkEntry *CheckEntry) {
	goInfo.VecBlockInfo = []BlockInfo{}
	if checkEntry != nil {
		atomic.StoreUint32(&checkEntry.Uint32NeedCheck, 0)
	}
}

func (goInfo *GoInfo) IsBlock() (boolIsBlock bool, strOp string) {
	boolIsBlock, strOp = false, ""
	boolIsSelect := false

	lock(&goInfo.Mu)
	defer unlock(&goInfo.Mu)
	if len(goInfo.VecBlockInfo) == 0 {
		return
	} else {
		boolIsBlock = true
	}

	// Now we compute strOp

	if len(goInfo.VecBlockInfo) > 1 {
		boolIsSelect = true
	} else if len(goInfo.VecBlockInfo) == 0 {
		print("Fatal in GoInfo.IsBlock(): goInfo.VecBlockInfo is not nil but lenth is 0\n")
	}

	if boolIsSelect {
		strOp = Select
	} else {
		for _, blockInfo := range goInfo.VecBlockInfo { // This loop will be executed only one time, since goInfo.VecBlockInfo's len() is 1
			strOp = blockInfo.StrOp
		}
	}

	return
}

// This function checks if the goroutine mapped with goInfo is currently blocking at an operation of chInfo.Chan
// If so, returns true and the string of channel operation
func (goInfo *GoInfo) IsBlockAtGivenChan(chInfo *ChanInfo) (boolIsBlockAtGiven bool, strOp string) {
	boolIsBlockAtGiven, strOp = false, ""

	lock(&goInfo.Mu)
	defer unlock(&goInfo.Mu)
	if goInfo.VecBlockInfo == nil {
		return
	}

	for _, blockInfo := range goInfo.VecBlockInfo {
		if blockInfo.Prim == chInfo {
			boolIsBlockAtGiven = true
			strOp = blockInfo.StrOp
			break
		}
	}

	return
}
