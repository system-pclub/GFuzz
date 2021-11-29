package bug

import (
	"gfuzz/pkg/utils/arr"
	"testing"
)

func TestGetListOfBugIDFromStdoutContentHappyNonBlocking(t *testing.T) {
	content := `=== RUN   TestGrpc1687
-----New NonBlocking Bug:
---Stack:
goroutine 9 [running]:
runtime.ReportNonBlockingBug(...)
	/usr/local/go/src/runtime/myoracle.go:538
command-line-arguments.(*serverHandlerTransport).do(0xc00005a560, 0x55efd8)
	/data/ziheng/shared/gotest/empirical/gobench/gobench/goker/nonblocking/grpc/1687/grpc1687_test.go:58 +0x2e5
command-line-arguments.(*serverHandlerTransport).Write(0xc00005a560)
	/data/ziheng/shared/gotest/empirical/gobench/gobench/goker/nonblocking/grpc/1687/grpc1687_test.go:75 +0x37
command-line-arguments.TestGrpc1687.func1(0xc00005a570)
	/data/ziheng/shared/gotest/empirical/gobench/gobench/goker/nonblocking/grpc/1687/grpc1687_test.go:177 +0x45
created by command-line-arguments.testHandlerTransportHandleStreams.func1
	/data/ziheng/shared/gotest/empirical/gobench/gobench/goker/nonblocking/grpc/1687/grpc1687_test.go:169 +0x3b

-----End Bug
panic: send on closed channel

goroutine 9 [running]:
command-line-arguments.(*serverHandlerTransport).do(0xc00005a560, 0x55efd8)
	/data/ziheng/shared/gotest/empirical/gobench/gobench/goker/nonblocking/grpc/1687/grpc1687_test.go:58 +0x2e5
command-line-arguments.(*serverHandlerTransport).Write(0xc00005a560)
	/data/ziheng/shared/gotest/empirical/gobench/gobench/goker/nonblocking/grpc/1687/grpc1687_test.go:75 +0x37
command-line-arguments.TestGrpc1687.func1(0xc00005a570)
	/data/ziheng/shared/gotest/empirical/gobench/gobench/goker/nonblocking/grpc/1687/grpc1687_test.go:177 +0x45
created by command-line-arguments.testHandlerTransportHandleStreams.func1
	/data/ziheng/shared/gotest/empirical/gobench/gobench/goker/nonblocking/grpc/1687/grpc1687_test.go:169 +0x3b

Process finished with exit code 1`

	bugIds, err := GetListOfBugIDFromStdoutContent(content)
	if err != nil {
		t.Fail()
	}
	if len(bugIds) != 1 {
		t.Fail()
	}
	if !arr.Contains(bugIds, "/data/ziheng/shared/gotest/empirical/gobench/gobench/goker/nonblocking/grpc/1687/grpc1687_test.go:58") {
		t.Fail()
	}
}

func TestGetListOfBugIDFromStdoutContentHappyBlocking(t *testing.T) {
	content := `=== RUN   TestBalancerUnderBlackholeNoKeepAliveDelete
Check bugs after 2 minutes
Failed to create file: 
-----New Blocking Bug:
---Blocking location:
/data/ziheng/shared/gotest/stubs/etcd/src/go.etcd.io/etcd/etcdserver/v3_server.go:838
---Primitive location:
/data/ziheng/shared/gotest/stubs/etcd/src/go.etcd.io/etcd/etcdserver/server.go:806
/data/ziheng/shared/gotest/stubs/etcd/src/go.etcd.io/etcd/etcdserver/server.go:809
/data/ziheng/shared/gotest/stubs/etcd/src/go.etcd.io/etcd/etcdserver/server.go:1215
---Primitive pointer:
0xc000442540
0xc000442580
0xc001d8bf40
-----End Bug
-----New Blocking Bug:
---Blocking location:
/data/ziheng/shared/gotest/stubs/etcd/pkg/mod/github.com/soheilhy/cmux@v0.1.4/cmux.go:229
---Primitive location:
/data/ziheng/shared/gotest/stubs/etcd/pkg/mod/github.com/soheilhy/cmux@v0.1.4/cmux.go:135
---Primitive pointer:
0xc001affc40
-----End Bug
{"level":"warn","ts":"2021-07-26T23:25:16.731-0400","caller":"clientv3/retry_interceptor.go:62","msg":"retrying of unary invoker failed","target":"endpoint://client-c437f924-48d8-4d97-a6e9-ed535f5f4462/localhost:37946441868244577890","attempt":0,"error":"rpc error: code = DeadlineExceeded desc = context deadline exceeded"}
    black_hole_test.go:277: #1: current error expected error
-----Withdraw prim:0xc000442540
{"level":"warn","ts":"2021-07-26T23:25:16.731-0400","caller":"clientv3/retry_interceptor.go:62","msg":"retrying of unary invoker failed","target":"endpoint://client-c437f924-48d8-4d97-a6e9-ed535f5f4462/localhost:37946441868244577890","attempt":0,"error":"rpc error: code = DeadlineExceeded desc = context deadline exceeded"}
    black_hole_test.go:277: #1: current error expected error
-----New Blocking Bug:
---Blocking location:
/data/ziheng/shared/gotest/stubs/etcd/src/go.etcd.io/etcd/etcdserver/v3_server.go:123456
---Primitive location:
/data/ziheng/shared/gotest/stubs/etcd/src/go.etcd.io/etcd/etcdserver/server.go:806
---Primitive pointer:
0xc000442540
-----End Bug
--- PASS: TestBalancerUnderBlackholeNoKeepAliveDelete (10.26s)
PASS`

	bugIds, err := GetListOfBugIDFromStdoutContent(content)
	if err != nil {
		t.Fail()
	}
	if len(bugIds) != 2 {
		t.Fail()
	}
	if !arr.Contains(bugIds, "/data/ziheng/shared/gotest/stubs/etcd/pkg/mod/github.com/soheilhy/cmux@v0.1.4/cmux.go:229") {
		t.Fail()
	}
	if !arr.Contains(bugIds, "/data/ziheng/shared/gotest/stubs/etcd/src/go.etcd.io/etcd/etcdserver/v3_server.go:123456") {
		t.Fail()
	}
}

func TestGetListOfBugIDFromStdoutContentEmpty(t *testing.T) {
	content := `
goroutine 3855 [running]:
github.com/prometheus/prometheus/tsdb/wal.(*WAL).run(0xc0002e7c20)
	/Users/xsh/code/prometheus/tsdb/wal/wal.go:372 +0x47a
created by github.com/prometheus/prometheus/tsdb/wal.NewSize
	/Users/xsh/code/prometheus/tsdb/wal/wal.go:302 +0x325
	`

	bugIds, err := GetListOfBugIDFromStdoutContent(content)
	if err != nil {
		t.Fail()
	}
	if len(bugIds) != 0 {
		t.Fail()
	}
}

func TestGetListOfBugIDFromStdoutCausedByPanic(t *testing.T) {
	content := `
panic: send on closed channel

goroutine 7 [running]:
fuzzer-toy/blocking/grpc/1353.(*roundRobin).watchAddrUpdates(0xc00001c810)
	/fuzz/target/blocking/grpc/1353/grpc1353_test.go:84 +0x10f
fuzzer-toy/blocking/grpc/1353.(*roundRobin).Start.func1(0xc00001c810)
	/fuzz/target/blocking/grpc/1353/grpc1353_test.go:52 +0x35
created by fuzzer-toy/blocking/grpc/1353.(*roundRobin).Start
	/fuzz/target/blocking/grpc/1353/grpc1353_test.go:50 +0x91
	`

	bugIds, err := GetListOfBugIDFromStdoutContent(content)
	if err != nil {
		t.Fail()
	}
	if bugIds == nil {
		t.Fail()
	}

	if !arr.Contains(bugIds, "/fuzz/target/blocking/grpc/1353/grpc1353_test.go:84") {
		t.Fail()
	}

}

func TestGetListOfBugIDFromStdoutSkipTimeout(t *testing.T) {
	content := `
-----New Blocking Bug:
---Blocking location:
/data/ziheng/shared/gotest/stubs/etcd/pkg/mod/github.com/soheilhy/cmux@v0.1.4/cmux.go:229
---Primitive location:
/data/ziheng/shared/gotest/stubs/etcd/pkg/mod/github.com/soheilhy/cmux@v0.1.4/cmux.go:135
---Primitive pointer:
0xc001affc40
-----End Bug

panic: test timed out after 1m0s

goroutine 468 [running]:
testing.(*M).startAlarm.func1()
	/usr/local/go/src/testing/testing.go:1700 +0xe5
created by time.goFunc
	/usr/local/go/src/time/sleep.go:180 +0x45
	`

	bugIds, err := GetListOfBugIDFromStdoutContent(content)
	if err != nil {
		t.Fail()
	}
	if bugIds == nil {
		t.Fail()
	}

	if !arr.Contains(bugIds, "/data/ziheng/shared/gotest/stubs/etcd/pkg/mod/github.com/soheilhy/cmux@v0.1.4/cmux.go:229") {
		t.Fail()
	}

}

func TestGetListOfBugIDFromPanicWithPanicStackTrace(t *testing.T) {

	content := `
panic: runtime error: invalid memory address or nil pointer dereference [recovered]
	panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x58 pc=0x900610]

goroutine 19 [running]:
testing.tRunner.func1.2(0x9df6c0, 0xe4a060)
	/usr/local/go/src/testing/testing.go:1143 +0x335
testing.tRunner.func1(0xc00034d180)
	/usr/local/go/src/testing/testing.go:1146 +0x4c2
panic(0x9df6c0, 0xe4a060)
	/usr/local/go/src/runtime/panic.go:965 +0x1b9
github.com/docker/docker/client.(*Client).ClientVersion(...)
	/go/src/github.com/docker/docker/client/client.go:197
github.com/docker/docker/client.TestNewClientWithOpsFromEnv(0xc00034d180)
	/go/src/github.com/docker/docker/client/client_test.go:100 +0x690
testing.tRunner(0xc00034d180, 0xaa34b8)
	/usr/local/go/src/testing/testing.go:1193 +0xef
created by testing.(*T).Run
	/usr/local/go/src/testing/testing.go:1238 +0x2b5`

	bugIds, err := GetListOfBugIDFromStdoutContent(content)
	if err != nil {
		t.Fail()
	}
	if bugIds == nil {
		t.Fail()
	}

	if !arr.Contains(bugIds, "/go/src/github.com/docker/docker/client/client.go:197") {
		t.Fail()
	}

}

func TestGetListOfBugIDFromStdoutContentOracleFatal(t *testing.T) {
	content := `End of unit test. Check bugs
fatal error: unexpected signal during runtime execution
[signal SIGSEGV: segmentation violation code=0x1 addr=0x1df7630 pc=0x47d165]
goroutine 135 [running]:
runtime.throw(0xee3442, 0x2a)
    /usr/local/go/src/runtime/panic.go:1117 +0x72 fp=0xc000d95888 sp=0xc000d95858 pc=0x441f92
runtime.sigpanic()
    /usr/local/go/src/runtime/signal_unix.go:718 +0x2e5 fp=0xc000d958c0 sp=0xc000d95888 pc=0x45a225
runtime.memmove(0xc000d96076, 0xefcc38, 0xefaa78)
    /usr/local/go/src/runtime/memmove_amd64.s:392 +0x485 fp=0xc000d958c8 sp=0xc000d958c0 pc=0x47d165
runtime.concatstrings(0x0, 0xc000d959a8, 0x3, 0x3, 0xc0006b6300, 0x76)
    /usr/local/go/src/runtime/string.go:52 +0x19e fp=0xc000d95960 sp=0xc000d958c8 pc=0x45f0be
runtime.concatstring3(0x0, 0xc0006b6300, 0x76, 0xefcc38, 0xefaa78, 0xefaa98, 0x1, 0xefaa98, 0x1)
    /usr/local/go/src/runtime/string.go:63 +0x47 fp=0xc000d959a0 sp=0xc000d95960 pc=0x45f2e7
runtime.CheckBlockEntry(0xc00047ccc0, 0x1d, 0x0)
    /usr/local/go/src/runtime/myoracle_tmp.go:226 +0x2fe fp=0xc000d95ab8 sp=0xc000d959a0 pc=0x439c1e
gooracle.CheckBugEnd(0xc0005d0870)
    /usr/local/go/src/gooracle/gooracle.go:325 +0x118 fp=0xc000d95bc8 sp=0xc000d95ab8 pc=0x5735f8
gooracle.AfterRunFuzz(0xc0005d0870)
    /usr/local/go/src/gooracle/gooracle.go:395 +0x59 fp=0xc000d95c20 sp=0xc000d95bc8 pc=0x573dd9
gooracle.AfterRun(0xc0005d0870)
    /usr/local/go/src/gooracle/gooracle.go:358 +0x4b fp=0xc000d95c38 sp=0xc000d95c20 pc=0x57394b
runtime.call16(0x0, 0xefc3c8, 0xc000b8a228, 0x800000008)
    /usr/local/go/src/runtime/asm_amd64.s:550 +0x3e fp=0xc000d95c58 sp=0xc000d95c38 pc=0x47a49e
runtime.reflectcallSave(0xc000d95d50, 0xefc3c8, 0xc000b8a228, 0x8)
    /usr/local/go/src/runtime/panic.go:877 +0x58 fp=0xc000d95c88 sp=0xc000d95c58 pc=0x441558
runtime.runOpenDeferFrame(0xc000535c20, 0xc000b8a1e0, 0xc0006efd98)
    /usr/local/go/src/runtime/panic.go:851 +0x62d fp=0xc000d95d10 sp=0xc000d95c88 pc=0x44122d
runtime.Goexit()
    /usr/local/go/src/runtime/panic.go:613 +0x1e5 fp=0xc000d95d98 sp=0xc000d95d10 pc=0x4406e5
testing.(*common).FailNow(0xc00035a780)
    /usr/local/go/src/testing/testing.go:741 +0x3c fp=0xc000d95db0 sp=0xc000d95d98 pc=0x523ebc
testing.(*common).Fatal(0xc00035a780, 0xc0006efee0, 0x1, 0x1)
    /usr/local/go/src/testing/testing.go:809 +0x78 fp=0xc000d95de8 sp=0xc000d95db0 pc=0x5247f8
google.golang.org/grpc/xds/internal/test_test.TestServerSideXDS_SecurityConfigChange(0xc00035a780)
    /fuzz/target/xds/internal/test/xds_server_integration_test.go:389 +0x8f7 fp=0xc000d95f80 sp=0xc000d95de8 pc=0xc52f77
testing.tRunner(0xc00035a780, 0xefb590)
    /usr/local/go/src/testing/testing.go:1193 +0xef fp=0xc000d95fd0 sp=0xc000d95f80 pc=0x525daf
runtime.goexit()
    /usr/local/go/src/runtime/asm_amd64.s:1371 +0x1 fp=0xc000d95fd8 sp=0xc000d95fd0 pc=0x47bf01
created by testing.(*T).Run
    /usr/local/go/src/testing/testing.go:1238 +0x2b5
goroutine 1 [chan receive]:
testing.(*T).Run(0xc00035a780, 0xedf2cb, 0x26, 0xefb590, 0x49f701)
    /usr/local/go/src/testing/testing.go:1239 +0x2dc
testing.runTests.func1(0xc00035a500)
    /usr/local/go/src/testing/testing.go:1511 +0x78
testing.tRunner(0xc00035a500, 0xc000237de0)
    /usr/local/go/src/testing/testing.go:1193 +0xef
testing.runTests(0xc0005d0858, 0x16c66e0, 0x4, 0x4, 0xc03c34318ea821cd, 0xdf924ed58, 0x16e0440, 0xec93c3)
    /usr/local/go/src/testing/testing.go:1509 +0x305
testing.(*M).Run(0xc000362680, 0x0)
    /usr/local/go/src/testing/testing.go:1417 +0x1eb
main.main()
    _testmain.go:49 +0x13`

	bugIds, err := GetListOfBugIDFromStdoutContent(content)
	if err != nil {
		// The text above are all in runtime, so there will be err reported at bug.go:80
		//t.Fail()
	}
	if len(bugIds) != 0 {
		t.Fail()
	}
}

func TestRandomOutput1(t *testing.T) {
	content := `
-----New Blocking Bug:
---Primitive location:
/fuzz/target/internal/transport/http2_client.go:287
/fuzz/target/internal/grpcsync/event.go:67
---Primitive pointer:
-----Withdraw prim:0xc0002a6400
0xc0002a69c0
0xc0002a6400
-----End Bug`
	bugIds, err := GetListOfBugIDFromStdoutContent(content)
	if err != nil {
		t.Fail()
	}

	if len(bugIds) != 0 {
		t.Fail()
	}

}

func TestRandomOutput(t *testing.T) {
	content := `
	=== RUN   TestPingPong
    end2end_test.go:3949: Running test in tcp-clear-v1-balancer environment...
    end2end_test.go:3949: Running test in tcp-tls-v1-balancer environment...
    end2end_test.go:3949: Running test in tcp-clear environment...
    end2end_test.go:3949: Running test in tcp-tls environment...
    end2end_test.go:3949: Running test in handler-tls environment...
    end2end_test.go:3949: Running test in no-balancer environment...
-----New Blocking Bug:
---Primitive location:
/fuzz/target/picker_wrapper.go:63
/fuzz/target/clientconn.go:562
---Primitive pointer:
0xc000686180
0xc0004961c0
-----End Bug
-----New Blocking Bug:
---Primitive location:
/fuzz/target/clientconn.go:562
/fuzz/target/picker_wrapper.go:63
---Primitive pointer:
0xc00031c700
0xc0004961c0
-----End Bug
-----New Blocking Bug:
---Primitive location:
/fuzz/target/internal/transport/http2_server.go:213
/fuzz/target/clientconn.go:562
---Primitive pointer:
0xc000604480
0xc0004961c0
-----End Bug
-----New Blocking Bug:
---Primitive location:
/fuzz/target/picker_wrapper.go:63
/fuzz/target/clientconn.go:562
---Primitive pointer:
0xc000155a40
0xc0004961c0
-----End Bug
-----New Blocking Bug:
---Primitive location:
/fuzz/target/picker_wrapper.go:63
/fuzz/target/clientconn.go:562
---Primitive pointer:
0xc000686a80
0xc0004961c0
-----End Bug
-----New Blocking Bug:
---Primitive location:
/fuzz/target/internal/transport/http2_server.go:213
/fuzz/target/clientconn.go:562
---Primitive pointer:
0xc00031ce40
0xc0004961c0
-----End Bug
-----New Blocking Bug:
---Primitive location:
/fuzz/target/picker_wrapper.go:63
/fuzz/target/clientconn.go:562
---Primitive pointer:
0xc0004961c0
0xc000154c40
-----End Bug
-----New Blocking Bug:
---Primitive location:
/fuzz/target/internal/transport/http2_server.go:213
/fuzz/target/clientconn.go:562
---Primitive pointer:
0xc0006040c0
0xc0004961c0
-----End Bug
-----New Blocking Bug:
---Primitive location:
/fuzz/target/clientconn.go:562
/fuzz/target/picker_wrapper.go:63
---Primitive pointer:
0xc00031ca40
0xc0004961c0
-----End Bug
-----New Blocking Bug:
---Primitive location:
/fuzz/target/picker_wrapper.go:63
/fuzz/target/clientconn.go:562
---Primitive pointer:
0xc000414880
0xc0004961c0
-----End Bug
-----New Blocking Bug:
---Primitive location:
/fuzz/target/clientconn.go:562
---Primitive pointer:
0xc0004961c0
-----End Bug
-----Withdraw prim:0xc0004961c0
End of unit test. Check bugs
---Stack:
goroutine 6 [running]:
runtime.DumpAllStack()
    /usr/local/go/src/runtime/myoracle_tmp.go:212 +0x85
gfuzz/pkg/oraclert.AfterRunFuzz(0xc000322000)
    /usr/local/go/src/gfuzz/pkg/oraclert/oracle.go:361 +0x7c
gfuzz/pkg/oraclert.AfterRun(0xc000322000)
    /usr/local/go/src/gfuzz/pkg/oraclert/oracle.go:327 +0x4b
google.golang.org/grpc/test.TestPingPong(0xc00025cc80)
    /fuzz/target/test/end2end_test.go:3945 +0x10d
testing.tRunner(0xc00025cc80, 0xb8ae20)
    /usr/local/go/src/testing/testing.go:1193 +0xef
created by testing.(*T).Run
    /usr/local/go/src/testing/testing.go:1238 +0x2b5
goroutine 1 [chan receive]:
testing.(*T).Run(0xc00025cc80, 0xb5d888, 0xc, 0xb8ae20, 0x49ad01)
    /usr/local/go/src/testing/testing.go:1239 +0x2dc
testing.runTests.func1(0xc00025ca00)
    /usr/local/go/src/testing/testing.go:1511 +0x78
testing.tRunner(0xc00025ca00, 0xc0001c3de0)
    /usr/local/go/src/testing/testing.go:1193 +0xef
testing.runTests(0xc00000c8d0, 0xfc1da0, 0xcb, 0xcb, 0xc05e32cd150bb46d, 0x6fedbb0aa, 0xfc7660, 0xb5e554)
    /usr/local/go/src/testing/testing.go:1509 +0x305
testing.(*M).Run(0xc000275ee0, 0x0)
    /usr/local/go/src/testing/testing.go:1417 +0x1eb
main.main()
    _testmain.go:447 +0x138
goroutine 18 [sleep]:
time.Sleep(0x4a817c800)
    /usr/local/go/src/runtime/time.go:193 +0xd2
gfuzz/pkg/oraclert.CheckBugLate()
    /usr/local/go/src/gfuzz/pkg/oraclert/oracle.go:220 +0x5d
created by gfuzz/pkg/oraclert.CheckBugStart
    /usr/local/go/src/gfuzz/pkg/oraclert/oracle.go:137 +0x39
--- PASS: TestPingPong (1.31s)
PASS
ok      google.golang.org/grpc/test 1.558s	
`
	bugIds, err := GetListOfBugIDFromStdoutContent(content)
	if err != nil {
		t.Fail()
	}

	if len(bugIds) != 0 {
		t.Fail()
	}

}
