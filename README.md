# GFuzz
Fuzzing concurrent Go programs

- [GFuzz](#gfuzz)
  - [Architecture](#architecture)
    - [GFuzz Fuzzer](#gfuzz-fuzzer)
    - [GFuzz Oracle](#gfuzz-oracle)
  - [Packages](#packages)
    - [pkg/gooracle](#pkggooracle)
    - [pkg/selefcm (select enforcement)](#pkgselefcm-select-enforcement)
    - [pkg/inst (instrumentation)](#pkginst-instrumentation)
      - [Built-in Passes](#built-in-passes)
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

### pkg/selefcm (select enforcement)

Package `selefcm` provides a list of strategies for application to choose proper select case by given a list of select choices (optional)

### pkg/inst (instrumentation)

Package `inst` provides modifying golang source code framework and utilities. It provides `InstPass` interface to easily write your own pass to instrument/modify/analysis golang source code.

#### Built-in Passes


<table>
<tr>
<th> Pass </th>
 <th> Description </th> 
 <th>Example</th>
</tr>

<tr>
<td>channel-record</td>
<td>record channel related operations like make, send, recv, close</td>
<td>

```go
//before
ch := make(chan int)

//after

ch := make(chan int
gooracle.StoreChMakeInfo(ch, <some random number>)
```
</td>
</tr>

<tr>
<td>mutex-record</td>
<td>record mutex related operations </td>
<td></td>
</tr>

<tr>
<td>wg-record</td>
<td>record WaitGroup related operations</td>
<td></td>
</tr>

<tr>
<td>cv-record</td>
<td>record Conditional Variable related operations</td>
<td></td>
</tr>

<tr>
<td>select-enforce</td>
<td>transform select into select with integer case (each case is one of original case and timeout)</td>
<td></td>
</tr>
</table>

## Dev
Since large parts of GFuzz are required instrumented Golang environment, we would suggest develop/test in universal Docker environment.

```bash

// The script will 
// 1. build a container with instrumented Golang environment 
// 2. mapping current directly and run the container
// 3. try `make test` after the container bring up!
$ ./script/dev.sh

```