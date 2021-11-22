# GFuzz
Fuzzing concurrent Go programs

- [GFuzz](#gfuzz)
  - [Architecture](#architecture)
    - [GFuzz Fuzzer](#gfuzz-fuzzer)
    - [GFuzz Oracle](#gfuzz-oracle)
  - [Packages](#packages)
    - [pkg/oraclert](#pkgoraclert)
    - [pkg/selefcm (select enforcement)](#pkgselefcm-select-enforcement)
    - [pkg/inst (instrumentation)](#pkginst-instrumentation)
    - [pkg/fuzz](#pkgfuzz)
    - [pkg/fuzzer](#pkgfuzzer)
    - [pkg/gexec](#pkggexec)
    - [pkg/inst](#pkginst)
    - [pkg/stats](#pkgstats)
    - [pkg/utils/**](#pkgutils)
    - [pkg/inst/pass (built-in passes)](#pkginstpass-built-in-passes)
  - [Dev](#dev)
    - [Prerequistes](#prerequistes)
    - [Build](#build)
    - [Useful Scripts](#useful-scripts)
      - [Manually Run a Test/Program after Instrumentation](#manually-run-a-testprogram-after-instrumentation)
  - [Executable `bin/inst`](#executable-bininst)
  - [Executable `bin/fuzzer`](#executable-binfuzzer)



## Architecture
GFuzz is composed by two parts: **Fuzzer** and **Oracle**.

### GFuzz Fuzzer
GFuzz Fuzzer generates different combination of `select` choices (forcing application to go with certain case).  It requires
help from instrumentation.

### GFuzz Oracle
GFuzz Oracle detects blocking & non-blocking issues during application is running. It requires help from both instrumentation
and golang package `runtime` patched in advance.

## Packages

### pkg/oraclert

Package `oraclert` is part of GFuzz Oracle. This package requires patched golang environment to work properly. It provides
1. Detecting blocking/non-blocking issue happened during application runtime.

### pkg/selefcm (select enforcement)

Package `selefcm` provides a list of strategies for application to choose proper select case by given a list of select choices (optional)

### pkg/inst (instrumentation)

Package `inst` provides modifying golang source code framework and utilities. It provides `InstPass` interface to easily write your own pass to instrument/modify/analysis golang source code.

### pkg/fuzz

Package `fuzz` provides 

### pkg/fuzzer

Package `fuzzer` provides 

### pkg/gexec

Package `gexec` provides 

### pkg/inst

Package `inst` provides 

### pkg/stats

Package `stats` provides 

### pkg/utils/**

Packages under `utils/**` provides 

### pkg/inst/pass (built-in passes)


<table>
<tr>
<th> Pass </th>
 <th> Description </th> 
 <th> Source</th>
 <th>Example</th>
</tr>

<tr>
<td>chrec</td>
<td>record channel related operations like make, send, recv, close</td>
<td><a href="pkg/inst/pass/chrec.go">pkg/inst/pass/chrec.go</a></td>
<td><a href="_examples/inst/chrec">_examples/inst/chrec</a></td>
</tr>

<tr>
<td>mtxrec</td>
<td>record mutex related operations </td>
<td><a href="pkg/inst/pass/mtxrec.go">pkg/inst/pass/mtxrec.go</a></td>
<td><a href="_examples/inst/mtxrec">_examples/inst/mtxrec</a></td>
</tr>

<tr>
<td>wgrec</td>
<td>record WaitGroup related operations</td>
<td><a href="pkg/inst/pass/wgrec.go">pkg/inst/pass/wgrec.go</a></td>
<td><a href="_examples/inst/wgrec">_examples/inst/wgrec</a></td>
</tr>

<tr>
<td>cvrec</td>
<td>record Conditional Variable related operations</td>
<td><a href="pkg/inst/pass/cvrec.go">pkg/inst/pass/cvrec.go</a></td>
<td><a href="_examples/inst/cvrec">_examples/inst/cvrec</a></td>
</tr>

<tr>
<td>selefcm</td>
<td>transform select into select with integer case (each case is one of original case and timeout)</td>
<td><a href="pkg/inst/pass/selefcm.go">pkg/inst/pass/selefcm.go</a></td>
<td><a href="_examples/inst/selefcm">_examples/inst/selefcm</a></td>
</tr>

<tr>
<td>oracle</td>
<td>insert function call to trigger oracle at the beginning of Test function or main program (TODO)</td>
<td><a href="pkg/inst/pass/oracle.go">pkg/inst/pass/oracle.go</a></td>
<td><a href="_examples/inst/oracle">_examples/inst/oracle</a></td>
</tr>

</table>

## Dev

### Prerequistes
Since large parts of GFuzz are required instrumented Golang environment, we would suggest develop/test in universal Docker environment.

```bash

// The script will 
// 1. build a container with instrumented Golang environment 
// 2. mapping current directly and run the container
// 3. try `make test` after the container bring up!
$ ./script/dev.sh

```

### Build

```bash
$ make
```

### Useful Scripts

#### Manually Run a Test/Program after Instrumentation
```
# For example, after instrument the gRPC code base by bin/inst, we can manully trigger a test in patched go runtime environment.
# Usually it is run in Docker so that patched go runtime will not effect your local environment.

export ORACLERT_CONFIG_FILE=/workspaces/GFuzz/tmp/grpc-go/ort_config
export ORACLERT_OUTPUT_FILE=/workspaces/GFuzz/tmp/grpc-go/ort_outputs
go test -run=TestGRPCLBStatsUnaryFailedToSend google.golang.org/grpc/balancer/grpclb
```

## Executable `bin/inst`

```
Usage:
  inst [OPTIONS] [Globs...]

Application Options:
      --pass=  A list of passes you want to use in this instrumentation
      --dir=   Instrument all go source files under this directory
      --file=  Instrument single go source file
      --out=   Output instrumented golang source file to the given file. Only allow when instrumenting single golang source file

Help Options:
  -h, --help   Show this help message
```

## Executable `bin/fuzzer`
```
Usage:

```