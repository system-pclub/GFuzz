# The code, analysis scripts and results for ASPLOS 2022 Artifact Evaluation

Version: 1.1\
Update:  Dec 01, 2021\
Paper:   Who Goes First? Detecting Go Concurrency Bugs via Message Reordering

This document is to help users reproduce the results we reported in our submission. 
It contains the following descriptions:

## 0. Artifact Expectation

The code and the scripts of our built tool are released in this repository. 
The detailed information of our experiments is released in an excel file. 
All the experiments can be executed using Docker 20.10.8. We expect users 
to use a docker of this version or higher to reproduce our experiments. 

## 1. Artifact Overview

Our paper presents GFuzz, a dynamic detector for channel-related concurrency
bugs in Go programs. For artifact evaluation, we release 
- (1) the tool we built, 
- (2) information of evaluated benchmarks, 
- (3) information of detected bugs, 
- (4) execution overhead of GFuzz's sanitizer, 
- (5) study results of whether 
GFuzz can help detect bugs in two public concurrency bug sets,
- and (6) scripts to compare the effectiveness of GFuzz's different features.  


Item (1) can be checked out by executing the following commands

``` bash
$ git clone https://github.com/system-pclub/GFuzz.git
$ cd GFuzz
$ git checkout asplos-artifact
```


Items (2), (3), (4) and (5) are released using a Google Sheet file [asplos-710-artifact](https://docs.google.com/spreadsheets/d/1tLcgsfYlll0g20KMYgDKkAtwZtk426dMSUZ6SvXk04s/edit?usp=sharing). 
Particularly, (2), (3) and (4) are related to Table 2 in the paper, and
(5) is to provide more information for Table 3 in the paper. 
All columns and tabs discussed later are in the Google Sheet file, unless otherwise specified. 

Item (6) is to reproduce Figure 5 in the paper. 


## 2. Tab Table-2-Benchmark

This tab shows the information of our evaluated benchmarks. 
Column F shows the versions we use to count the line numbers (Column D)
and the unit-test numbers (Column E). 

Prerequiste: clone all applications with specific commit version:
```bash
$ cd benchmark
$ clone-repos.sh ./repos
```

The line number of all applications 
can be counted by executing the following command:

``` bash
$ cd benchmark
$ loc.sh ./repos
```

The unit tests of an applications (e.g., etcd) at a particular version 
(e.g., see ./benchmark/clone-repos.sh)
can be counted by executing the following command:

``` bash
$ cd benchmark

# If you run before, skip it.
$ ./clone-repos.sh repos

# If you run before, skip it.
$ ./build.sh

# run script to count specific app. Same way for the others
$ ./benchmark.sh count-tests --dir /builder/etcd/native
```

## 3. Tab Table-2-Bug 

This tab shows the detailed information of each detected bugs, including at which application
version we find the bug (Column B), where we report the bug (Column E), what is the current 
status of the filed bug report (Columns G-J), which category the bug belongs to (Columns L-V), whether 
GCatch can detect the bug (Column X), what is the reason if GCatch fails (Columns Y-AC), 
and the unit test we used to find the bug (columns AE-AF). 

### 3.1. Kubernetes

Users can execute the following command to apply GFuzz to 
fuzz Kubernetes at a particular version 
(e.g., 97d40890d00acf721ecabb8c9a6fec3b3234b74b):

``` bash
$ ./scripts/fuzz-git.sh https://github.com/kubernetes/kubernetes 97d40890d00acf721ecabb8c9a6fec3b3234b74b $(pwd)/tmp/out

```

By default, GFuzz will use all unit tests of a given application. To ease
the reproduction of our results, we enhance GFuzz to only use one unit 
test (columns AE-AF). For example, users can execute the following command
to inspect whether GFuzz can still detect the bug at row 4. 

``` bash
$ ./scripts/fuzz-git.sh https://github.com/kubernetes/kubernetes 97d40890d00acf721ecabb8c9a6fec3b3234b74b $(pwd)/tmp/out --pkg k8s.io/kubernetes/pkg/kubelet/cm/devicemanager --func TestAllocate
```

### 3.2. Docker


GFuzz needs to build unit tests first and then conducts the fuzzing. 
Since etcd is monorepo of many Golang modules, you need to manually build its tests using the commands in [docker/builder/entrypoint.sh](docker/builder/entrypoint.sh). 
Docker (moby) does not support go module, so that you have to turn off GO111MODULE. 
Specifically, you need to use the following commands to build Docker’s tests. 

```bash
OUTPUT_DIR=<output dir>
export GO111MODULE=off
pkg_list=$(go list github.com/docker/docker/... | grep -vE "(integration)")

for pkg in $pkg_list
do
    echo "generating test bin for $pkg"
    name=$(echo "$pkg" | sed "s/\//-/g")
    go test -c -o $OUTPUT_DIR/moby/native/$name.test $pkg
done
export GO111MODULE=on
```

For etcd and Docker, the recommended way is to use benchmark/clone-repos.sh and benchmark/build.sh to build tests firstly (feel free to comment out applications you are not interested in [docker/builder/entrypoint.sh](docker/builder/entrypoint.sh)), and 
then to use the following command to do the fuzzing. 

```bash
$ ./scripts/fuzz-testbins.sh <testbin dir> <output dir> [optional flags for fuzzer]
```

### 3.3. Prometheus

To apply GFuzz, we need to change one testing setting of Prometheus. Since all detected bugs of Prometheus are from two versions, we create two repositories for the two versions ([Prometheus-1](https://github.com/gfuzz-asplos/prometheus-e0f1506254688cec85276cc939aeb536a4e029d1) and [Prometheus-2](https://github.com/gfuzz-asplos/prometheus-f08c89e569b2421bcc8ef7caf585fd8d3c2ccaba)) and conduct the required change. To reproduce the experiments on Prometheus, users can directly use the two created repositories. 

The first version is e0f1506254688cec85276cc939aeb536a4e029d1. Users can execute the following command to apply GFuzz to the version. 

``` bash
$ ./scripts/fuzz-git.sh https://github.com/gfuzz-asplos/prometheus-e0f1506254688cec85276cc939aeb536a4e029d1 ba019add3f94b5ef224fbf2e537afe4f3878ffbe $(pwd)/tmp/out
```

To only use one unit test in the first version, users can execute the following command.

``` bash
$ ./scripts/fuzz-git.sh https://github.com/gfuzz-asplos/prometheus-e0f1506254688cec85276cc939aeb536a4e029d1 ba019add3f94b5ef224fbf2e537afe4f3878ffbe $(pwd)/tmp/out --pkg <pkg_name> --func <pkg_name> 
```


The second version is f08c89e569b2421bcc8ef7caf585fd8d3c2ccaba. Users can execute the following command to apply GFuzz to that version. 

``` bash
$ ./scripts/fuzz-git.sh https://github.com/gfuzz-asplos/prometheus-f08c89e569b2421bcc8ef7caf585fd8d3c2ccaba 141e016e260734ade6b2ba2cbdc8435bfce70262 $(pwd)/tmp/out
```

To only use one unit test in the first version, users can execute the following command.

``` bash
$ ./scripts/fuzz-git.sh https://github.com/gfuzz-asplos/prometheus-f08c89e569b2421bcc8ef7caf585fd8d3c2ccaba 141e016e260734ade6b2ba2cbdc8435bfce70262 $(pwd)/tmp/out --pkg <pkg_name> --func <pkg_name> 
```

### 3.4. etcd


### 3.5. Go-Ethereum

Users can execute the following command to apply GFuzz to 
Go-Ethereum at a particular version:

``` bash
$ ./scripts/fuzz-git.sh https://github.com/ethereum/go-ethereum <commit_hash> $(pwd)/tmp/out

# e.g., 
# ./scripts/fuzz-git.sh https://github.com/ethereum/go-ethereum 56e9001a1a8ddecc478943170b00207ef46109b9 $(pwd)/tmp/out
```

Users can execute the following command
to force GFuzz to only use one unit test. 

``` bash
$ ./scripts/fuzz-git.sh https://github.com/ethereum/go-ethereum <commit_hash> $(pwd)/tmp/out --pkg <pkg_name> --func <func_name>

# e.g.,
# ./scripts/fuzz-git.sh https://github.com/ethereum/go-ethereum 56e9001a1a8ddecc478943170b00207ef46109b9 $(pwd)/tmp/out --pkg github.com/ethereum/go-ethereum/console --func TestInteractive
```


### 3.6. TiDB

The following command applies GFuzz to 
a particular version of TiDB:

``` bash
$ ./scripts/fuzz-git.sh https://github.com/pingcap/tidb <commit_hash> $(pwd)/tmp/out

```

Users can execute the following command
to only fuzz one unit test. 

``` bash
$ ./scripts/fuzz-git.sh https://github.com/pingcap/tidb <commit_hash> $(pwd)/tmp/out --pkg <pkg_name> --func <func_name>
```

### 3.7. gRPC-go

For gRPC and prometheus, some testing settings need to be changed firstly. 
Details can be found at [docs/fuzz-trick.md](docs/fuzz-trick.md). 
Then, you can use the following command to run the fuzzing. 

```bash
$ ./scripts/fuzz-mount.sh <repository dir> <output dir> [optional flags for fuzzer]
```

We compare GFuzz with GCatch in our evaluation. To check whether 
GCatch can detect a bug, please see instruction at section [Using GCatch to test GFuzz bugs](https://github.com/system-pclub/GFuzz#7-using-gcatch-to-test-gfuzz-bugs) below.





## 4. Tab Table-2-Overhead

This tab shows the overhead of GFuzz's sanitizer. 

Users can execute the following command to measure the overhead
on an application (e.g., grpc): 

``` bash

$ cd benchmark

# If you run before, skip it.
$ ./clone-repos.sh repos

# If you run before, skip it.
$ ./build.sh

# /builder is the mapped directory of host directory 'tmp/builder', which is output of ./build.sh
$ ./benchmark.sh benchmark --dir /builder/grpc/native --mode native
$ ./benchmark.sh benchmark --dir /builder/grpc/inst --mode inst
```


## 5. Tab Table-3

In Section 7.2 of the paper, we manually studied whether reordering messages can
help detect channel-related bugs in two public sets of Go concurrency bugs. 
This tab shows the detailed labeling. 



## 6. Figure 5 of the paper


We evaluate GFuzz on grpc in Figure 5. 

``` bash
# First of all, setup grpc
$ cd ~
$ git clone https://github.com/grpc/grpc-go.git
$ cd grpc-go

# Checkout the version we are evaluating on
$ git checkout 9280052d36656451dd7568a18a836c2a74edaf6c 
```

It is required to use a fresh grpc folder each time we begin the GFuzz fuzzing. 
To run GFuzz on the grpc library:
``` bash
# Copy a new fresh grpc folder to fuzz on. 
$ cp -r /path/to/grpc/ /path/to/grpc_0/
$ sudo ./script/fuzz-mount.sh /path/to/grpc_0/ /path/to/output/folder/GFuzz_out/
```
For fuzzing without feedback:

``` bash
$ cp -r /path/to/grpc/ /path/to/grpc_1/
$ sudo ./script/fuzz-mount.sh /path/to/grpc_1/ /path/to/output/folder/GFuzz_no_feedback/ --isIgnoreFeedback 1
```

For fuzzing without mutations:

``` bash
$ cp -r /path/to/grpc/ /path/to/grpc_2/
$ sudo ./script/fuzz-mount.sh /path/to/grpc_2/ /path/to/output/folder/GFuzz_no_mutation/ --isNoMutation 1
```

For fuzzing without oracle:

``` bash
$ ./benchmark/clone-repos.sh ./repos

# If you have ran this script before, skip it
$ ./benchmark/build.sh

# Build an uninstrumented grpc
# /builder is the mapped directory of host directory 'tmp/builder', which is output of ./build.sh
$ ./benchmark.sh benchmark --dir /builder/grpc/native --mode native

# Run GFuzz with the compiled grpc
$ sudo ./script/fuzz-testbins.sh ./tmp/builder/grpc/native/ /path/to/output/folder/GFuzz_no_oracle/
```

After fuzzing for 3 hours with each configs, we can plot Figure 5 using the following script: 

``` bash
# Install python3 dependent libraries
$ pip3 install matplotlib click datetime

$ python3 ./script/plot_Figure_5.py --with-feedback-path /path/to/output/folder/GFuzz_out/ --no-feedback-path /path/to/output/folder/GFuzz_no_feedback/ --no-mutation-path /path/to/output/folder/GFuzz_no_mutation/ --no-oracle-path /path/to/output/folder/GFuzz_no_oracle/
```

## 7. Using GCatch to test GFuzz bugs: 

Let's set up GCatch using a Docker environment

``` bash
$ cd ~
$ git clone https://github.com/system-pclub/GCatch.git
$ cd ./GCatch/GCatch
# This might take a while
$ sudo docker build -t gcatch_test .

# Upon finish the previous command
$ sudo docker run -it gcatch_test
```

Now we are in a Docker terminal, the following commands are all executed inside this Docker environment. 
Additionally, all bugs being detected by GFuzz can be found in [Tab Table-2-Bug](https://github.com/system-pclub/GFuzz#3-tab-table-2-bug). 

For testing grpc: 
All grpc packages start with *google.golang.org/grpc*. If the bug is located in grpc folder *internal/resolver*, then the module path would be *google.golang.org/grpc/internal/resolver*.

``` bash
$ cd /playground
$ git clone https://github.com/grpc/grpc-go.git
$ cd /playground/grpc-go

# checkout the specific buggy version

$ GO111MODULE=on GCatch -mod -mod-abs-path=/playground/grpc-go -mod-module-path=module_path -compile-error
```

For testing etcd: 
All etcd packages start with *go.etcd.io/etcd/.../v3*. For example, if the bug is located in etcd folder *tests/integration/snapshot*, then the module path would be *go.etcd.io/etcd/tests/v3/integration/snapshot*.

``` bash
$ cd /playground
$ git clone https://github.com/etcd-io/etcd.git
$ cd /playground/etcd

# checkout the specific buggy version

$ GO111MODULE=on GCatch -mod -mod-abs-path=/playground/etcd -mod-module-path=module_path -compile-error
```

For testing Kubernetes:
All Kubernetes packages start with *k8s.io/kubernetes*. For example, if the bug is located in Kubernets folder *pkg/kubelet/nodeshutdown/systemd*, then the module path is: *k8s.io/kubernetes/pkg/kubelet/nodeshutdown/systemd*

``` bash
$ cd /playground
$ git clone https://github.com/kubernetes/kubernetes.git
$ cd /playground/kubernetes

# checkout the specific buggy version

$ GO111MODULE=on GCatch -mod -mod-abs-path=/playground/kubernetes -mod-module-path=module_path  -compile-error
```

For testing Prometheus:
All Prometheus packages start with *github.com/prometheus/prometheus*. For example, if the bug is located in Prometheus folder *storage/remote*, then the module path is: *github.com/prometheus/prometheus/storage/remote*

``` bash
$ cd /playground
$ https://github.com/prometheus/prometheus.git
$ cd /playground/prometheus

# checkout the specific buggy version

$ GO111MODULE=on GCatch -mod -mod-abs-path=/playground/prometheus -mod-module-path=module_path  -compile-error
```

For testing go-Ethereum:
All go-Ethereum packages start with *github.com/ethereum/go-ethereum/*. For example, if the bug is located in Prometheus folder *core*, then the module path is: *github.com/ethereum/go-ethereum/core*

``` bash
$ cd /playground
$ git clone https://github.com/ethereum/go-ethereum.git
$ cd /playground/go-ethereum

# checkout the specific buggy version

$ GO111MODULE=on GCatch -mod -mod-abs-path=/playground/go-ethereum -mod-module-path=module_path  -compile-error
```

For testing tidb:
For tidb module names, most tidb packages start with *github.com/pingcap/tidb*. For example, if the bug is located in tidb folder *./ddl/*, then the module path is: *github.com/pingcap/tidb/ddl*. However, for bugs reported in *badger*, the module path is: *github.com/pingcap/badger*

``` bash
$ cd /playground
$ git clone https://github.com/pingcap/tidb.git
$ cd /playground/tidb

# checkout the specific buggy version

$ GO111MODULE=on GCatch -mod -mod-abs-path=/playground/tidb -mod-module-path=module_path  -compile-error
```

For testing Moby(Docker):
Docker fuzzing uses a slightly different routine. The Docker code must be stored in path */go/src/github.com/*. 

``` bash
$ cd /go/src/github.com
$ mkdir -p docker
$ cd docker
$ git clone https://github.com/moby/moby.git
$ mv moby docker
$ GO111MODULE=off GCatch -path=/go/src/github.com/docker/docker -include=github.com/docker/docker -r -compile-error
```


If any bug was found from any programs above, GCatch would output a bug report similar to the following format. 

``` bash
Successfully built whole program. Now running checkers
----------Bug[1]----------
        Type: BMOC      Reason: One or multiple channel operation is blocked.
-----Blocking at:
        File: /playground/prometheus/web/web.go:949
-----Blocking Path NO. 0
ChanMake :/playground/prometheus/web/web.go:947:12       '✓'
Chan_op :/playground/prometheus/web/web.go:948:13        '✓'
Chan_op :/playground/prometheus/web/web.go:949:12        Blocking
If :/playground/prometheus/web/web.go:949:22     '✗'
Return :/playground/prometheus/web/web.go:950:13         '✗'
```


## 8. Example Output of Fuzzing




Output Dir Structure

fuzzer.log: all the activities that produced by the fuzzer
tbin/*: the compiled test binary of target Golang repository by fuzzer. 
exec/*: full history run triggered by the fuzzer. Each folder represents one run


Here is a real world example for fuzzing gRPC.

Result Explain
First we can go through the fuzzer.log:

```
2021/11/29 23:22:36 [worker 3] received 25389-rand-google.golang.org-grpc-internal-transport.test-TestKeepaliveServerClosesUnresponsiveClient-1106
2021/11/29 23:22:40 [worker 3] found unique bug: /fuzz/target/internal/transport/keepalive_test.go:381
2021/11/29 23:22:40 [worker 3] found 1 unique bug(s)
```

If we saw log like ‘found xxx unique bug(s)’, this means that previous run with id ‘25389-rand-google.golang.org-grpc-internal-transport.test-TestKeepaliveServerClosesUnresponsiveClient-1106’ detected​​ a bug. If we look at exec folder in the output directory, you should see a folder with the exact name: {outputdir}/exec/25389-rand-google.golang.org-grpc-internal-transport.test-TestKeepaliveServerClosesUnresponsiveClient-1106

The output exists in exec/{run id}/stdout, usually the user should expect a pattern of ‘-----New Blocking Bug’, which is printed by oracle runtime.
Primitive location indicates where the primitive is been make ( something like position of ch := make(chan struct{}))
Primitive pointer indicates the memory address of this primitive at runtime (this is for oracle for eliminating false positive)
```
-----New Blocking Bug:
---Primitive location:
/fuzz/target/internal/transport/keepalive_test.go:381
---Primitive pointer:
0xc0002b0000
-----End Bug
```

At the bottom of stdout, we should expect an full goroutine stack trace, usually you will find a goroutine blocked with corresponding primitive (we can see from filename in the stack usually)

```
goroutine 50 [chan send]:
google.golang.org/grpc/internal/transport.TestKeepaliveServerClosesUnresponsiveClient.func2(0x941dc0, 0xc000286000, 0xc0003a68b0, 0xc0002aa000)
	/fuzz/target/internal/transport/keepalive_test.go:387 +0x134
created by google.golang.org/grpc/internal/transport.TestKeepaliveServerClosesUnresponsiveClient
	/fuzz/target/internal/transport/keepalive_test.go:382 +0x535
```

