# GFuzz
Fuzzing concurrent Go programs

- [GFuzz](#gfuzz)
  - [Architecture](#architecture)
    - [GFuzz Fuzzer](#gfuzz-fuzzer)
    - [GFuzz Oracle](#gfuzz-oracle)
  - [Packages](#packages)
    - [pkg/gooracle](#pkggooracle)
    - [pkg/selectstrtg](#pkgselectstrtg)
  - [Dev](#dev)



## Architecture
GFuzz is composed by two parts: **Fuzzer** and **Oracle**.

### GFuzz Fuzzer
GFuzz Fuzzer generates different combination of `select` choices (forcing application to go with certain case).  It requires
help from instrumentation.

### GFuzz Oracle
GFuzz Oracle detects blocking & non-blocking issues during application is running. It requires help from both instrumentation
and golang package `runtime` patched in advance.

## Packages

### pkg/gooracle

Package `gooracle` is part of GFuzz Oracle. This package requires patched golang environment to work properly. It provides
1. Detecting blocking/non-blocking issue happened during application runtime.

### pkg/selectstrtg

Package `selectstrtg` provides a list of strategies for application to choose proper select case by given a list of select choices (optional)

## Dev
Since large parts of GFuzz are required instrumented Golang environment, we would suggest develop/test in universal Docker environment.

```bash

// The script will 
// 1. build a container with instrumented Golang environment 
// 2. mapping current directly and run the container
// 3. try `make test` after the container bring up!
$ ./script/dev.sh

```