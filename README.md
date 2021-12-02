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
    - [Example](#example)
  - [Executable `bin/fuzzer`](#executable-binfuzzer)
    - [Example](#example-1)



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
      --pass=               A list of passes you want to use in this instrumentation
      --dir=                Instrument all go source files under this directory
      --file=               Instrument single go source file
      --out=                Output instrumented golang source file to the given file. Only allow when instrumenting single golang source
                            file
      --statsOut=           Output statistics
      --version             Print version and exit
      --check-syntax-err
      --recover-syntax-err
      --parallel=
      --cpuprofile=

Help Options:
  -h, --help                Show this help message
```

### Example
```bash
# Suppose /abc/def is directory contains go.mod
$ inst --dir /abc/def
```

## Executable `bin/fuzzer`
```
Usage:
  fuzzer [OPTIONS]

Application Options:
      --gomod=            Directory contains go.mod
      --func=             Only run specific test function in the test
      --pkg=              Only run test functions in the specific package
      --bin=              A list of globs for Go test bins.
      --ortconfig=        Only run once with given ortconfig
      --out=              Directory for fuzzing output
      --parallel=         Number of workers to fuzz parallel (default: 5)
      --instStats=        This parameter consumes a file path to a statistics file generated by isnt.
      --version           Print version and exit
      --globalTuple       Whether prev_location is global or per channel
      --scoreSdk          Recording/scoring if channel comes from Go SDK
      --scoreAllPrim      Recording/scoring other primitives like Mutex together with channel
      --timeDivideBy=     Durations in time/sleep.go will be divided by this int number
      --oraclertdebug
      --isIgnoreFeedback  Is ignoring the feedback, and save every mutated seed into the fuzzing queue
      --randMutateEnergy= Determine the energy of random mutations. If == 100 (default), then each seed would mutate 100 times in the rand mutation stage
      --isUsingScore      Is using score to priority testing case.

Help Options:
  -h, --help              Show this help message
```

### Example
```bash
# Suppose /abc/def is directory contains go.mod

# Fuzz whole module, it implies every tests in every packages will be fuzzed
$ fuzzer --gomod /abc/def

# Fuzz selected package(s)
$ fuzzer --gomod /abc/def --pkg gomodulename/aaa --pkg gomodulename/bbb

# Fuzz selected func(s)
$ fuzzer --gomod /abc/def --pkg gomodulename/aaa --func TestABC

# If you have compiled a list of test binary files
$ fuzzer --bin /abc/*.test
```

---------
Asplos Artifact Tutorial

1. asplos 710 table 2 benchmark

[benchmark/README.md](benchmark/README.md)

2. asplos 710 table 2 bug

In short, there are two main scripts for achieving out-of-box fuzzing:

scripts/fuzz-git.sh
	This script is suitable for those Golang repositories that can directly be used with fuzzing, without extra changes/replacements before fuzzing.
Usage:
./scripts/fuzz-git.sh <GIT URL> <GIT COMMIT> <OUTPUT DIR> [optional flags for fuzzer] 


scripts/fuzz-mount.sh
	This script is suitable for those Golang repositories that require extra changes before fuzzing. For example, gRPC has lots of Test Suite(single test entrypoint but triggered lots of tests)

Usage:
	./scripts/fuzz-mount.sh <REPO DIR> <OUTPUT DIR> [optional flags for fuzzer] 

3. To reproduce Figure 5, Contributions of GFuzz components:

We evaluate GFuzz on grpc in Figure 5. 

``` bash
# First of all, setup grpc
cd ~
git clone https://github.com/grpc/grpc-go.git
cd grpc-go

# Checkout the version we are evaluating on
git checkout 9280052d36656451dd7568a18a836c2a74edaf6c 
```

It is required to use a fresh grpc folder with the correct version every time when we are begin the GFuzz fuzzing. 
To run GFuzz on the grpc library:
``` bash
cp -r /path/to/grpc/ /path/to/grpc_0/
sudo ./script/fuzz-mount.sh /path/to/grpc_0/ /path/to/output/folder/GFuzz_out/
```

Additionally, to run GFuzz on grpc, but ignore the fuzzer feedback:

``` bash
cp -r /path/to/grpc/ /path/to/grpc_1/
sudo ./script/fuzz-mount.sh /path/to/grpc_1/ /path/to/output/folder/GFuzz_no_feedback/ --isIgnoreFeedback 1
```

For fuzzing without mutations:

``` bash
cp -r /path/to/grpc/ /path/to/grpc_2/
sudo ./script/fuzz-mount.sh /path/to/grpc_2/ /path/to/output/folder/GFuzz_no_mutation/ --isNoMutation 1
```

For fuzzing without oracle:

``` bash
# Copy grpc into the GFuzz benchmark folder. The grpc code must be exactly in the ./benchmark/tmp/builder folder
cp -r /path/to/grpc/ ./benchmark/tmp/builder/grpc/

# If you have ran this script before, skip it
./benchmark/build.sh

# Build an uninstrumented grpc
# /builder is the mapped directory of host directory 'tmp/builder', which is output of ./build.sh
./benchmark.sh benchmark --dir /builder/grpc/native --mode native

# Run GFuzz with the compiled grpc
sudo ./script/fuzz-testbins.sh ./benchmark/tmp/builder/grpc/native/ /path/to/output/folder/GFuzz_no_oracle/
```

After fuzzing for 3 hours for each config, we can plot Figure 5: 

``` bash
# Install python3 dependent libraries
pip3 install matplotlib click datetime

python3 ./script/plot_Figure_5.py --with-feedback-path /path/to/output/folder/GFuzz_out/ --no-feedback-path /path/to/output/folder/GFuzz_no_feedback/ --no-mutation-path grpc_no_feedback_all_stage_0 --no-oracle-path /path/to/output/folder/GFuzz_no_oracle/
```

4. Reproduce GFuzz bugs using GCatch:

