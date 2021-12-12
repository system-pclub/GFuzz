package bug

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type bugInfo struct {
	id string // could be the blocking location or the creation location of blocking primitive;
	// example:"/data/ziheng/shared/gotest/stubs/etcd/src/go.etcd.io/etcd/pkg/schedule/schedule.go:219"
	strBlockLoc      string   // one of strBlockLoc and strStack may be empty
	vecPrimPtr       []string // All primitives that triggered this bug; Used to withdraw a bug; example:"0xc001985e80"
	vecPrimCreateLoc []string // could be empty for not instrumented primitives
	strStack         string
}

// GetListOfBugIDFromStdoutContent parses gooracle related information(Bug ID) from stdout file content.
// A Bug ID is where channel was been created (identified by file path + line number) and
// that channel will trigger block bug.
//
//     Example of blocking bug from stdout
//-----New Blocking Bug:
//---Blocking location:
///data/ziheng/shared/gotest/stubs/etcd/src/go.etcd.io/etcd/pkg/schedule/schedule.go:219
//---Primitive location:
///usr/local/go/src/context/context.go:363
///data/ziheng/shared/gotest/stubs/etcd/src/go.etcd.io/etcd/pkg/schedule/schedule.go:66
//---Primitive pointer:
//0xc000428480
//0xc0003d1680
//-----End Bug
// In the bug above, "/data/ziheng/shared/gotest/stubs/etcd/src/go.etcd.io/etcd/pkg/schedule/schedule.go:219" is Bug ID
// A new mechanism is introduced: if we later observe "-----Withdraw prim:0xc000428480", then this bug is withdrawed,
// because the primitive is used again after we report the bug and before the unit test completes
//
//     Example of non blocking bug from stdout
//-----New NonBlocking Bug:
//---Stack:
//goroutine 9 [running]:
//runtime.ReportNonBlockingBug(...)
///usr/local/go/src/runtime/myoracle.go:538
//command-line-arguments.(*serverHandlerTransport).do(0xc00005a560, 0x55efd8)
///data/ziheng/shared/gotest/empirical/gobench/gobench/goker/nonblocking/grpc/1687/grpc1687_test.go:58 +0x2e5
//command-line-arguments.(*serverHandlerTransport).Write(0xc00005a560)
///data/ziheng/shared/gotest/empirical/gobench/gobench/goker/nonblocking/grpc/1687/grpc1687_test.go:75 +0x37
//command-line-arguments.TestGrpc1687.func1(0xc00005a570)
///data/ziheng/shared/gotest/empirical/gobench/gobench/goker/nonblocking/grpc/1687/grpc1687_test.go:177 +0x45
//created by command-line-arguments.testHandlerTransportHandleStreams.func1
///data/ziheng/shared/gotest/empirical/gobench/gobench/goker/nonblocking/grpc/1687/grpc1687_test.go:169 +0x3b
//
//-----End Bug

func GetListOfBugIDFromStdoutContent(c string) ([]string, error) {
	lines := strings.Split(c, "\n")
	ids := make(map[string]struct{})
	result := []string{}
	mapBlockingBug := make(map[*bugInfo]struct{})
	numOfLines := len(lines)
	for idx, line := range lines {
		if line == "" {
			continue
		}

		// trim space and tab
		line = strings.TrimLeft(line, " \t")
		if strings.HasPrefix(line, "-----New NonBlocking Bug:") ||
			strings.HasPrefix(line, "panic: ") || strings.HasPrefix(line, "fatal error") {

			if strings.HasPrefix(line, "panic: test timed out after") {
				continue
			}
			idLineIdx := idx + 1

			// Skip file location(s) that belongs to my*.go until find the bug root cause
			// if this line is not from our my*.go files, then it is where bug happened
			for {
				if idLineIdx >= numOfLines {
					return nil, fmt.Errorf("total line %d, target bug ID line at %d", numOfLines, idLineIdx)
				}

				if lines[idLineIdx] == "---Stack:" {
					idLineIdx += 3
					continue
				}

				if strings.Contains(lines[idLineIdx], "/go/src/runtime") || strings.Contains(lines[idLineIdx], "src/sync") || strings.Contains(lines[idLineIdx], "/go/src/testing") {
					idLineIdx += 2
					continue
				}

				if strings.HasPrefix(lines[idLineIdx], "\t") && !strings.Contains(lines[idLineIdx], "panic: ") {
					// first line with filename + linenumer
					break
				}

				idLineIdx += 1
			}

			targetLine := lines[idLineIdx]
			println(targetLine)
			id, err := getFileAndLineFromStacktraceLine(targetLine)

			if err != nil {
				log.Printf("getFileAndLineFromStacktraceLine failed: %s", err)
				continue
			}
			ids[id] = struct{}{}

			if strings.HasPrefix(line, "fatal error") {
				// if this is a fatal bug, ignore all other bugs found in this unit test, because fatal program can cause lots of strange things
				oneID := []string{}

				// let's watch for two places indicating that our oracle is causing fatal. If they exist, don't record this bug
				// Why the first place will fatal is unsure. See https://github.com/golang/go/issues/41285
				boolFoundOracleFatal := false
				for i := idx; i < numOfLines; i++ {
					sucLine := lines[i]
					if strings.Contains(sucLine, "runtime/map_fast64.go:291") || strings.Contains(sucLine, "runtime/string.go:63") {
						boolFoundOracleFatal = true
						break
					} else if strings.Contains(sucLine, "runtime/myoracle.go:") {
						if isContainMayFatalOracle(sucLine) {
							boolFoundOracleFatal = true
							break
						}
					}
				}
				if !boolFoundOracleFatal {
					oneID = []string{id}
				}

				return oneID, nil
			}

		} else if strings.HasPrefix(line, "-----New Blocking Bug:") {
			// find "-----End Bug"
			idEnd, err := findSucLineEqual(idx, lines, "-----End Bug")

			if err != nil {
				log.Printf("find \"-----End Bug\" failed: %s", err)
				continue
			}

			// then only need to inspect lines betwee idx and idEnd
			linesThisBug := lines[idx : idEnd+1]
			bug := &bugInfo{}

			// find "---Blocking location:", if any
			idBlockLoc, err := findSucLineEqual(0, linesThisBug, "---Blocking location:")

			if err != nil {
				// this is OK, some bug reports don't contain "---Blocking location:"
			} else {
				if idBlockLoc+1 >= len(linesThisBug) {
					log.Printf("the next line of \"---Blocking location:\" doesn't exist")
				} else {
					bug.strBlockLoc = linesThisBug[idBlockLoc+1]
				}
			}

			// find "---Primitive location:"
			idPrimLoc, err := findSucLineEqual(0, linesThisBug, "---Primitive location:")

			if err != nil {
				log.Printf("find \"---Primitive location:\" failed: %s", err)
				continue
			}

			idPrimLocEnd, err := findSucLinePrefix(idPrimLoc, linesThisBug, "---", "")

			if err != nil {
				log.Printf("find \"---\" after \"---Primitive location:\" failed: %s", err)
				continue
			}

			for i := idPrimLoc + 1; i < idPrimLocEnd; i++ {
				if strings.HasPrefix(linesThisBug[i], "[oraclert]") {
					continue
				}
				bug.vecPrimCreateLoc = append(bug.vecPrimCreateLoc, linesThisBug[i])
			}

			// find "---Primitive pointer:"
			idPrimPtr, err := findSucLineEqual(0, linesThisBug, "---Primitive pointer:")

			if err != nil {
				log.Printf("find \"---Primitive pointer:\" failed: %s", err)
				continue
			}

			idPrimPtrEnd, err := findSucLinePrefix(idPrimPtr, linesThisBug, "---", "-----Withdraw")

			if err != nil {
				log.Printf("find \"---\" after \"---Primitive pointer:\" failed: %s", err)
				continue
			}

			for i := idPrimPtr + 1; i < idPrimPtrEnd; i++ {
				if !strings.HasPrefix(linesThisBug[i], "0x") {
					continue
				}
				bug.vecPrimPtr = append(bug.vecPrimPtr, linesThisBug[i])
			}

			// find "---Stack:", if any
			idStack, err := findSucLineEqual(0, linesThisBug, "---Stack:")

			if err != nil {
				// this is OK, some bug reports don't contain "---Stack:"
			} else {
				idStackEnd, err := findSucLinePrefix(idStack, linesThisBug, "---", "")

				if err != nil {
					log.Printf("find \"---\" after \"---Stack:\" failed: %s", err)
					continue
				}

				for i := idStack + 1; i < idStackEnd; i++ {
					if strings.HasPrefix(linesThisBug[i], "[oraclert]") {
						continue
					}
					bug.strStack += linesThisBug[i]
				}
			}

			// if strBlockLoc is not empty, bug.id is strBlockLoc; else bug.id is the first not empty primLoc; else bug.id is strStack
			if bug.strBlockLoc != "" {
				bug.id = bug.strBlockLoc
			} else {
				for _, primLoc := range bug.vecPrimCreateLoc {
					if primLoc != "" {
						bug.id = primLoc
						break
					}
				}
				//if bug.id == "" { // let's not do this, strStack can't provide any useful information
				//	bug.id = bug.strStack
				//}
			}

			if bug.id != "" {
				mapBlockingBug[bug] = struct{}{}
			}

		} else if strings.HasPrefix(line, "-----Withdraw prim:") {
			// We need to delete all bugs in mapBlockingBug whose vecPrimStr contains this ptr
			strPtr := strings.TrimPrefix(line, "-----Withdraw prim:")
			for bug := range mapBlockingBug {
				for _, primPtr := range bug.vecPrimPtr {
					if strPtr == primPtr {
						delete(mapBlockingBug, bug)
						break
					}
				}
			}
		}
	}

	// for survival bugs in this map, we believe they are real bugs
	for bug := range mapBlockingBug {
		ids[bug.id] = struct{}{}
	}

	for id := range ids {
		result = append(result, id)
	}

	return result, nil
}

// getFileAndLineFromStacktraceLine returns only <file>:<line>
// from string with format <file>:<line> [<stack offset>]
func getFileAndLineFromStacktraceLine(line string) (string, error) {
	targetLine := strings.TrimLeft(line, " \t")
	parts := strings.Split(targetLine, " ")
	fileAndLine := parts[0]

	if len(strings.Split(fileAndLine, ":")) != 2 {
		return "", fmt.Errorf("malformed stacktrace, expected format: <file>:<line> [<stack offset>], got %s", targetLine)
	}

	return fileAndLine, nil
}

func findSucLineEqual(idx int, lines []string, strTarget string) (int, error) {
	for i := idx + 1; i < len(lines); i++ {
		if lines[i] == strTarget {
			return i, nil
		}
	}
	return -1, fmt.Errorf("malformed bug, can't find string %s since line %d", strTarget, idx)
}

func findSucLinePrefix(idx int, lines []string, strTarget string, skipPrefix string) (int, error) {
	for i := idx + 1; i < len(lines); i++ {
		if skipPrefix != "" && strings.HasPrefix(lines[i], skipPrefix) {
			continue
		}
		if strings.HasPrefix(lines[i], strTarget) {
			return i, nil
		}
	}
	return -1, fmt.Errorf("malformed bug, can't find prefix %s since line %d", strTarget, idx)
}

// isContainMayFatalOracle checks whether the line contains myoracle.go:200~225.
//
func isContainMayFatalOracle(str string) bool {
	for i := 200; i < 226; i++ {
		if strings.Contains(str, "runtime/myoracle.go:"+strconv.Itoa(i)) {
			return true
		}
	}
	return false
}
