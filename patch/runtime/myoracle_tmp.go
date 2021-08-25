package runtime

// A temporary oracle for fuzzer
// Deprecated:
//func TmpDumpBlockingInfo() (retStr string, foundBug bool) {
//	retStr = ""
//	foundBug = false
//	SleepMS(500)
//	lock(&muMap)
//outer:
//	for gid, strInfo := range mys.mpGoID2Info {
//		if gid != 1 { // No need to print the main goroutine
//			str := string(strInfo)
//			switch true {
//			case Index(str, "tools/cache/shared_informer.go:628") >= 0:
//				continue
//			case Index(str, "k8s.io/client-go/tools/cache/reflector.go:373") >= 0:
//				continue
//			case Index(str, "k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:1169") >= 0:
//				continue
//			case Index(str, "k8s.io/apimachinery/pkg/watch/mux.go:247") >= 0:
//				continue
//			case Index(str, "k8s.io/apimachinery/pkg/util/wait/wait.go:167") >= 0:
//				continue
//			case Index(str, "scheduler/internal/cache/debugger/debugger.go:63") >= 0:
//				continue
//			case Index(str, "k8s.io/client-go/tools/cache/shared_informer.go:772") >= 0:
//				continue
//			case Index(str, "vendor/k8s.io/client-go/tools/record/event.go:301") >= 0:
//				continue
//			case Index(str, "vendor/k8s.io/client-go/tools/cache/shared_informer.go:742 ") >= 0:
//				continue
//			case Index(str, "k8s.io/client-go/tools/cache/reflector.go:463 ") >= 0:
//				continue
//			case Index(str, "/vendor/") >= 0 && Index(str, "Lock(") == -1:
//				continue
//			case Index(str, "k8s.io/klog:1169") >= 0:
//				continue
//			case Index(str, "/home/luy70/go/src") >= 0:
//				continue
//			case Index(str, "=====") >= 0:
//				continue
//			case Index(str, "=====") >= 0:
//				continue
//			case Index(str, "=====") >= 0:
//				continue
//			case Index(str, "=====") >= 0:
//				continue
//			case Index(str, "=====") >= 0:
//				continue
//			case Index(str, "=====") >= 0:
//				continue
//			case Index(str, "=====") >= 0:
//				continue
//			case Index(str, "=====") >= 0:
//				continue
//			case Index(str, "=====") >= 0:
//				continue
//			case Index(str, "=====") >= 0:
//				continue
//			case Index(str, "=====") >= 0:
//				continue
//			default:
//			}
//
//			stackSingleGo := ParseStackStr(str)
//			if len(stackSingleGo.VecFuncName) == 0 {
//				retStr += "Warning in TmpDumpBlockingInfo: empty VecFunc*\n"
//				continue
//			}
//			firstFuncName := stackSingleGo.VecFuncName[0]
//
//			if Index(firstFuncName, "runtime.TmpBeforeBlock") > -1 {
//				// get the next func
//				if len(stackSingleGo.VecFuncFile) < 2 || len(stackSingleGo.VecFuncLine) < 2 { // unexpected problem: no func after TmpBeforeBlock
//					retStr += "Warning in TmpDumpBlockingInfo: no func after TmpBeforeBlock\n"
//					continue outer
//				}
//				nextFuncFile := stackSingleGo.VecFuncFile[1]
//				nextFuncLine := stackSingleGo.VecFuncLine[1]
//
//
//				if nextFuncFile != strSDKPath + "/src/sync/mutex.go" {
//					// case 1: from channel op
//					if _, reported := ReportedPlace[nextFuncFile + nextFuncLine]; reported {
//						continue outer
//					} else {
//						ReportedPlace[nextFuncFile + nextFuncLine] = struct{}{}
//					}
//					if indexSDK := Index(nextFuncFile, strSDKPath); indexSDK > -1 {
//						if indexSync := Index(nextFuncFile, strSDKPath + "/src/sync"); indexSync == -1 { // In SDK and not in our hacked sync, don't report
//							continue outer
//						}
//					}
//				} else {
//					// case 2: from Lock op
//					if len(stackSingleGo.VecFuncFile) < 3 || len(stackSingleGo.VecFuncLine) < 3 { // unexpected problem: no func after Lock
//						retStr += "Warning in TmpDumpBlockingInfo: no func after mutex.Lock\n"
//						continue outer
//					}
//					nextFuncFile := stackSingleGo.VecFuncFile[2]
//					nextFuncLine := stackSingleGo.VecFuncLine[2]
//					if _, reported := ReportedPlace[nextFuncFile + nextFuncLine]; reported {
//						continue outer
//					} else {
//						ReportedPlace[nextFuncFile + nextFuncLine] = struct{}{}
//					}
//				}
//
//			} else {
//				retStr += "Warning in TmpDumpBlockingInfo: the first func is not TmpBeforeBlock\n"
//			}
//			// delete the BeforeBlock function
//			for {
//				indexTBB := Index(str, "runtime.TmpBeforeBlock()\n\t" + strSDKPath + "/src/runtime/myoracle_tmp.go:")
//				if indexTBB == -1 {
//					break
//				} else {
//					str = str[:indexTBB] + str[indexTBB + 77:]
//				}
//			}
//
//			retStr += "-----New Blocking Bug:\n" + str + "\n"
//			print(retStr)
//			foundBug = true
//
//		}
//	}
//	unlock(&muMap)
//	return
//}

func init() {
	MapSelectInfo = make(map[string]SelectInfo)
	//MapInput = make(map[string]SelectInput)
}

var BoolDebug = false


var MuBlockEntry mutex
var MapBlockEntry map[*BlockEntry]struct{} = make(map[*BlockEntry]struct{})

type BlockEntry struct {
	VecPrim []PrimInfo
	StrOpPosition string
	CurrentGoInfo *GoInfo
}

func EnqueueBlockEntry(vecPrim []PrimInfo, op string) *BlockEntry {
	entry := &BlockEntry{
		VecPrim:       []PrimInfo{},
		StrOpPosition: "",
		CurrentGoInfo: CurrentGoInfo(),
	}
	for _, prim := range vecPrim {
		entry.VecPrim = append(entry.VecPrim, prim)
	}
	var layer int
	switch op {
	case Recv, Send:
		layer = 3
	case Select:
		layer = 2
	case MuLock:
		layer = 2
	}

	_, strFile, intLine, _ := Caller(layer)
	entry.StrOpPosition = strFile + ":" + Itoa(intLine)
	if entry.StrOpPosition == "/data/ziheng/shared/gotest/stubs/etcd/src/go.etcd.io/etcd/client/client.go:708" {
		//println("enqueue blockentry for 708:", entry)
	}

	lock(&MuBlockEntry)
	MapBlockEntry[entry] = struct{}{}
	unlock(&MuBlockEntry)

	return entry
}

func DequeueBlockEntry(entry *BlockEntry) {
	lock(&MuBlockEntry)
	if entry.StrOpPosition == "/data/ziheng/shared/gotest/stubs/etcd/src/go.etcd.io/etcd/client/client.go:708" {
		//println("dequeue blockentry for 708:", entry)
	}
	delete(MapBlockEntry, entry)
	unlock(&MuBlockEntry)
}

// If PrimInfo.LoadMonitor is 1, a bug has been reported based on the assumption that this prim won't be reached again
// However, this Monitor() is still invoked, meaning our assumption is incorrect. Withdraw the bug
// If this is really a bug, it will be reported later again
func Monitor(prim PrimInfo) {
	if prim.LoadMonitor() == 1 {
		prim.SetMonitor(0)
		str := "-----Withdraw prim:" + FnPointer2String(prim) + "\n"
		print(str)
		lock(&MuWithdraw)
		StrWithdraw += str
		unlock(&MuWithdraw)
	}
}

var MuWithdraw mutex
var StrWithdraw string

func CheckBlockEntry() (strReturn string, foundBug bool) {
	strReturn = ""
	foundBug = false
	// for all BlockEntry
	//		print it
	// 		for all prim in it
	//			set PrimInfo.Watch = true
	lock(&MuBlockEntry)
	for entry, _ := range MapBlockEntry {
		if Index(entry.StrOpPosition, strSDKPath) >= 0 { // don't report bugs in SDK. testing is often blocking
			continue
		}
		foundBug = true
		lock(&MuReportBug)
		strReturn += "-----New Blocking Bug:\n"
		strReturn += "---Blocking location:\n" + entry.StrOpPosition + "\n"
		strReturn += "---Primitive location:\n"
		for _, prim := range entry.VecPrim {
			strReturn += prim.StringDebug() + "\n"
			prim.SetMonitor(1)
		}
		strReturn += "---Primitive pointer:\n"
		for _, prim := range entry.VecPrim {
			strReturn += FnPointer2String(prim) + "\n"
		}
		strReturn += "-----End Bug\n"
		unlock(&MuReportBug)
	}
	unlock(&MuBlockEntry)
	return
}

var FnPointer2String func(interface{}) string
