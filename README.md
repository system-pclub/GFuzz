# The code, analysis scripts and results for ASPLOS 2022 Artifact Evaluation

Version: 1.2\
Update:  Jan 20, 2022\
Paper:   Who Goes First? Detecting Go Concurrency Bugs via Message Reordering

This document is to help users reproduce the results we reported in our submission. 
It contains the following descriptions:

## 0. Artifact Expectation

The code and the scripts of our built tool and the paper's final version are 
released in this repository. The detailed information of our experiments is 
released in an excel file. All the experiments can be executed using Docker 20.10.8. 
We expect users to use a docker of this version or higher to reproduce our experiments. 

## 1. Artifact Overview

Our paper presents GFuzz, a dynamic detector for channel-related concurrency
bugs in Go programs. For artifact evaluation, we release 
- (1) the tool we built,
- (2) the paper's final version,
- (3) information of evaluated benchmarks, 
- (4) information of detected bugs, 
- (5) execution overhead of GFuzz's sanitizer, 
- and (6) scripts to compare the effectiveness of GFuzz's different features.  


Items (1), (2) and (6) can be checked out by executing the following commands

``` bash
$ git clone https://github.com/system-pclub/GFuzz.git
$ cd GFuzz
$ git checkout asplos-artifact
```


Items (3), (4) and (5) are released using a Google Sheet file [asplos-710-artifact-2](https://docs.google.com/spreadsheets/d/11lUctFnLQAdGj_GK3wVfqMq_dHotW100rI96w-ez1tU/edit?usp=sharing). 
Particularly, (3), (4) and (5) are related to Table 2 in the paper. 
All columns and tabs discussed later are in the Google Sheet file, 
unless otherwise specified. 

Item (6) is to reproduce Figure 7 in the paper. 


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

# run script to count specific app. Same way for the others.
# Make sure that --dir is /builder/xxx/inst.
$ ./benchmark.sh count-tests --dir /builder/etcd/inst
```

## 3. Tab Table-2-Bug 

This tab shows the detailed information of each detected bugs, including at which application
version we find the bug (Column B), where we report the bug (Column H), what is the current 
status of the filed bug report (Columns J-L), which category the bug belongs to (Columns Q-Y), whether 
GCatch can detect the bug (Column AA), what is the reason if GCatch fails (Columns AB-AF), 
and the unit test we used to find the bug (columns AG-AH). 

### 3.1. Kubernetes

Users can execute the following command to apply GFuzz to 
fuzz Kubernetes at a particular version 
(e.g., 97d40890d00acf721ecabb8c9a6fec3b3234b74b):

``` bash
$ ./scripts/fuzz-git.sh https://github.com/kubernetes/kubernetes 97d40890d00acf721ecabb8c9a6fec3b3234b74b $(PWD)/tmp/out

```

By default, GFuzz will use all unit tests of a given application. To ease
the reproduction of our results, we enhance GFuzz to only use one unit 
test (columns AG-AH). For example, users can execute the following command
to inspect whether GFuzz can still detect the bug at row 4. 

``` bash
$ ./scripts/fuzz-git.sh https://github.com/kubernetes/kubernetes 97d40890d00acf721ecabb8c9a6fec3b3234b74b $(PWD)/tmp/out --pkg k8s.io/kubernetes/pkg/kubelet/cm/devicemanager --func TestAllocate
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

- The first version is e0f1506254688cec85276cc939aeb536a4e029d1. Users can execute the following commands to apply GFuzz to the version or one unit test in the version. 

``` bash
# fuzz the whole version
$ ./scripts/fuzz-git.sh https://github.com/gfuzz-asplos/prometheus-e0f1506254688cec85276cc939aeb536a4e029d1 ba019add3f94b5ef224fbf2e537afe4f3878ffbe $(PWD)/tmp/out
# fuzz one unit test
$ ./scripts/fuzz-git.sh https://github.com/gfuzz-asplos/prometheus-e0f1506254688cec85276cc939aeb536a4e029d1 ba019add3f94b5ef224fbf2e537afe4f3878ffbe $(PWD)/tmp/out --pkg <pkg_name> --func <func_name>
```

- The second version is f08c89e569b2421bcc8ef7caf585fd8d3c2ccaba. Users can execute the following commands to apply GFuzz to the version or one unit test in the version. 

``` bash
# fuzz the whole version
$ ./scripts/fuzz-git.sh https://github.com/gfuzz-asplos/prometheus-f08c89e569b2421bcc8ef7caf585fd8d3c2ccaba 141e016e260734ade6b2ba2cbdc8435bfce70262 $(PWD)/tmp/out
# fuzz one unit test
$ ./scripts/fuzz-git.sh https://github.com/gfuzz-asplos/prometheus-f08c89e569b2421bcc8ef7caf585fd8d3c2ccaba 141e016e260734ade6b2ba2cbdc8435bfce70262 $(PWD)/tmp/out --pkg <pkg_name> --func <func_name> 
```


### 3.4. etcd

For etcd, the recommended way is to use the following command to build its tests firstly. 

``` bash
# build all etcd tests at a particular version
$ ./scripts/prepare-etcd.sh <COMMIT_HASH>
# e.g.,
# ./scripts/prepare-etcd.sh 6bb26ef008f5465bd11b078f0a2e3ae95fdc6d4a
```

Then, users can use the following commands to fuzz the whole built etcd or one unit test in etcd. 

``` bash
# fuzz the whole etcd
$ ./scripts/fuzz-testbins.sh $(PWD)/tmp/builder/etcd/inst $(PWD)/tmp/out
# fuzz one unit test
$ ./scripts/fuzz-testbins.sh $(PWD)/tmp/builder/etcd/inst $(PWD)/tmp/out --func <func_name>
```


### 3.5. Go-Ethereum

Users can execute the following command to apply GFuzz to 
Go-Ethereum at a particular version:

``` bash
$ ./scripts/fuzz-git.sh https://github.com/ethereum/go-ethereum <commit_hash> $(PWD)/tmp/out

# e.g., 
# ./scripts/fuzz-git.sh https://github.com/ethereum/go-ethereum 56e9001a1a8ddecc478943170b00207ef46109b9 $(PWD)/tmp/out
```

Users can execute the following command
to force GFuzz to only use one unit test. 

``` bash
$ ./scripts/fuzz-git.sh https://github.com/ethereum/go-ethereum <commit_hash> $(PWD)/tmp/out --pkg <pkg_name> --func <func_name>

# e.g.,
# ./scripts/fuzz-git.sh https://github.com/ethereum/go-ethereum 56e9001a1a8ddecc478943170b00207ef46109b9 $(PWD)/tmp/out --pkg github.com/ethereum/go-ethereum/console --func TestInteractive
```


### 3.6. TiDB

The following command applies GFuzz to 
a particular version of TiDB:

``` bash
$ ./scripts/fuzz-git.sh https://github.com/pingcap/tidb <commit_hash> $(PWD)/tmp/out

```

Users can execute the following command
to only fuzz one unit test. 

``` bash
$ ./scripts/fuzz-git.sh https://github.com/pingcap/tidb <commit_hash> $(PWD)/tmp/out --pkg <pkg_name> --func <func_name>
```

### 3.7. gRPC-go

Similar to Prometheus, we need to change testing settings to apply GFuzz to gRPC-go. 
All detected bugs of gRPC-go come from three versions. Thus, we create three repositories 
respectively ([grpc-go-1](https://github.com/gfuzz-asplos/grpc-go-0bc741730b8171fc51cdaf826caea5119c411009), 
[grpc-go-2](https://github.com/gfuzz-asplos/grpc-go-83f9def5feb388c4fd7e6586bd55cf6bf6d46a01), 
and [grpc-go-3](https://github.com/gfuzz-asplos/grpc-go-9280052d36656451dd7568a18a836c2a74edaf6c)) 
and conduct the required changes. 

- The first version is 0bc741730b8171fc51cdaf826caea5119c411009. Users can execute the following command to apply GFuzz to the version or one unit test in the version. 

``` bash
# fuzz the whole application
$ ./scripts/fuzz-git.sh https://github.com/gfuzz-asplos/grpc-go-0bc741730b8171fc51cdaf826caea5119c411009 fe42d65231bf2c83c940db3b46849e250c3bdf2b $(PWD)/tmp/out
# fuzz one unit test
$  ./scripts/fuzz-git.sh https://github.com/gfuzz-asplos/grpc-go-0bc741730b8171fc51cdaf826caea5119c411009 fe42d65231bf2c83c940db3b46849e250c3bdf2b $(PWD)/tmp/out --pkg <pkg_name> --func <func_name> 
```

- The second version is 83f9def5feb388c4fd7e6586bd55cf6bf6d46a01. Users can execute the following command to apply GFuzz to the version or one unit test in the version. 

``` bash
# fuzz the whole application
$ ./scripts/fuzz-git.sh https://github.com/gfuzz-asplos/grpc-go-83f9def5feb388c4fd7e6586bd55cf6bf6d46a01 b95c0c0923d938b8acb7c841f0a04ade8f7d5fbf $(PWD)/tmp/out
# fuzz one unit test
$ ./scripts/fuzz-git.sh https://github.com/gfuzz-asplos/grpc-go-83f9def5feb388c4fd7e6586bd55cf6bf6d46a01 b95c0c0923d938b8acb7c841f0a04ade8f7d5fbf $(PWD)/tmp/out --pkg <pkg_name> --func <func_name> 
```


- The third version is 9280052d36656451dd7568a18a836c2a74edaf6c. Users can execute the following command to apply GFuzz to the version or one unit test in the version. 

``` bash
# fuzz the whole application
$ ./scripts/fuzz-git.sh https://github.com/gfuzz-asplos/grpc-go-9280052d36656451dd7568a18a836c2a74edaf6c 93d5a0f32dadc51585082f0d7786605b65fa6160 $(PWD)/tmp/out
# fuzz one unit test
$ ./scripts/fuzz-git.sh https://github.com/gfuzz-asplos/grpc-go-9280052d36656451dd7568a18a836c2a74edaf6c 93d5a0f32dadc51585082f0d7786605b65fa6160 $(PWD)/tmp/out --pkg <pkg_name> --func <func_name> 
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
$ ./benchmark.sh benchmark --dir /builder/grpc-go/native --mode native --out /builder/out/grpc-go-native.out
$ ./benchmark.sh benchmark --dir /builder/grpc-go/inst --mode inst --out /builder/out/grpc-go-inst.out

# After you have both results, compare common parts of them
$ ./filter.py ../tmp/builder/out/grpc-go-native.out ../tmp/builder/out/grpc-go-inst.out

# you should see following output:
common tests: 832
first average 0.1982 # second/test, first means first arg, which is ../tmp/builder/out/grpc-go-native.out
seond average 0.2032 # second/test, second means second arg, which is ../tmp/builder/out/grpc-go-inst.outß
```


## 5. Figure 5 of the paper


We evaluate GFuzz on grpc in Figure 5. 




``` bash
# If you have ran this script before, skip it
$ ./benchmark/clone-repos.sh ./repos

# If you have ran this script before, skip it
$ ./benchmark/build.sh
```

```bash
# Run gfuzz without feedback
./scripts/fuzz-testbins.sh $(pwd)/tmp/builder/grpc-go/inst ~/gfuzz/output/grpc-nfb-fixed --ignorefeedback --fixedsetimeout
# Run gfuzz with feedback
./scripts/fuzz-testbins.sh $(pwd)/tmp/builder/grpc-go/inst ~/gfuzz/output/grpc-se --scoreenergy --allowdupcfg --fixedsetimeout
# Run gfuzz without select enforcement
./scripts/fuzz-testbins.sh $(pwd)/tmp/builder/grpc-go/inst ~/gfuzz/output/grpc-nose --scoreenergy --allowdupcfg --nose
# Run gfuzz without oracle(sanitizer)
./scripts/fuzz-testbins.sh $(pwd)/tmp/builder/grpc-go/inst ~/gfuzz/output/grpc-nooracle --scoreenergy --allowdupcfg --fixedsetimeout --nooracle
```

After fuzzing for 12 hours with each configs, we can check when and what bugs they have found by:

```bash
$ ./scripts/analyze.py --bug-analyze --gfuzz-out-dir <gfuzz output>

# you should be able to see something like following:

bug statistics:
used (hours), buggy primitive location, gfuzz exec
0.05 h,/repos/grpc-go/xds/internal/testutils/fakeserver/server.go:160,1053-rand-google.golang.org-grpc-xds-internal-client-v2.test-TestLDSHandleResponseWithoutWatch-60
0.09 h,/repos/grpc-go/balancer/grpclb/grpclb_test.go:831,1168-rand-google.golang.org-grpc-balancer-grpclb.test-TestDropRequest-132

```

## 6. Using GCatch to test GFuzz bugs: 

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


## 7. Example Output of Fuzzing




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

