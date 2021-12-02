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
(1) the tool we built, 
(2) information of evaluated benchmarks, 
(3) information of detected bugs, 
(4) execution overhead of GFuzz's sanitizer, 
and (5) study results of whether 
GFuzz can help detect bugs in two public concurrency bug sets. 

Item (1) can be checked out by executing the following commands

``` bash
git clone https://github.com/system-pclub/GFuzz.git

cd GFuzz

git checkout asplos-artifact
```


Items (2), (3), (4) and (5) are released using a Google Sheet file "asplos-710-artifact" 
(https://docs.google.com/spreadsheets/d/1tLcgsfYlll0g20KMYgDKkAtwZtk426dMSUZ6SvXk04s/edit#gid=0). 
All columns and tabs discussed later are in the Google Sheet file, unless otherwise specified. 



## 2. Tab Table-2-Benchmark

This tab shows the information of our evaluated benchmarks. 
Column F shows the versions we use to count the line numbers (Column D)
and the unit-test numbers (Column E). 

The line number of an application (e.g., Kubernetes) at a particular version 
(e.g., 97d40890d00acf721ecabb8c9a6fec3b3234b74b)
can be counted by executing the following command:

``` bash
XXX
```

The unit tests of an application (e.g., Kubernetes) at a particular version 
(e.g., 97d40890d00acf721ecabb8c9a6fec3b3234b74b)
can be counted by executing the following command:

``` bash
XXX
```

## 3. Tab Table-2-Bug 

This tab shows the detailed information of the detected bugs, including which application
version we found a bug (Column B), where we report a bug (Column E),  
what is the current status of a filed bug report (Columns G--J), bug categories (Columns L--V), 
whether GCatch can detect a bug (Column X), 
the reasons why GCatch fails (Columns Y--AC), and the unit test we used to find a bug (columns AE--AF). 

Users can execute the following command to apply GFuzz to 
fuzz an application (e.g., Kubernetes) of a particular version 
(e.g., 97d40890d00acf721ecabb8c9a6fec3b3234b74b):

``` bash
XXX
```

By default, GFuzz will use all unit tests of a given application. To ease
the reproduction of our results, we enhance GFuzz to only use one unit 
test (columns AE--AF). For example, users can execute the following command
to inspect whether GFuzz can still detect the bug at row XXX. 

``` bash
XXX
```

We compare GFuzz with GCatch in our evaluation. To check whether 

GCatch can detect a bug (e.g., XXX), users can execute the following command:


``` bash
XXX
```




## 4. Tab Table-2-Overhead

This tab shows the overhead of GFuzz’s sanitizer. 

Users can execute the following command to measure the overhead
on an application (e.g., XXX): 

``` bash
XXX
```


## 5. Tab Table-3

In Section 7.2 of the paper, we manually studied whether reordering messages can
help detect channel-related bugs in two public sets of Go concurrency bugs. 
This tab shows the detailed labeling. 



## 6. Figure 5 of the paper


We evaluate GFuzz on grpc in Figure 5. 

``` bash
# First of all, setup grpc
cd ~
git clone https://github.com/grpc/grpc-go.git
cd grpc-go

# Checkout the version we are evaluating on
git checkout 9280052d36656451dd7568a18a836c2a74edaf6c 
```

It is required to use a fresh grpc folder each time we begin the GFuzz fuzzing. 
To run GFuzz on the grpc library:
``` bash
# Copy a new fresh grpc folder to fuzz on. 
cp -r /path/to/grpc/ /path/to/grpc_0/
sudo ./script/fuzz-mount.sh /path/to/grpc_0/ /path/to/output/folder/GFuzz_out/
```
For fuzzing without feedback:

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

After fuzzing for 3 hours with each configs, we can plot Figure 5 using the following script: 

``` bash
# Install python3 dependent libraries
pip3 install matplotlib click datetime

python3 ./script/plot_Figure_5.py --with-feedback-path /path/to/output/folder/GFuzz_out/ --no-feedback-path /path/to/output/folder/GFuzz_no_feedback/ --no-mutation-path grpc_no_feedback_all_stage_0 --no-oracle-path /path/to/output/folder/GFuzz_no_oracle/
```

## 7. Using GCatch to test GFuzz bugs: 

Let's set up GCatch using a Docker environment

``` bash
cd ~
git clone https://github.com/system-pclub/GCatch.git
cd ./GCatch/GCatch
# This might take a while
sudo docker build -t gcatch_test .

# Upon finish the previous command
sudo docker run -it gcatch_test
```

Now we are in a Docker terminal, the following commands are all executed inside this Docker environment. 
Additionally, all bugs being detected by GFuzz can be found in Google Sheets: **XXX**, table *Table-2-Bug*. 

For testing grpc: 
All grpc packages start with *google.golang.org/grpc**. If the bug is located in grpc folder *internal/resolver*, then the module path would be *google.golang.org/grpc/internal/resolver*.

``` bash
cd /playground
git clone https://github.com/grpc/grpc-go.git
cd /playground/grpc-go.git

# checkout the specific buggy version

GO111MODULE=on GCatch -mod -mod-abs-path=/playground/grpc-go -mod-module-path=module_path -compile-error
```


For testing etcd:
All etcd package start with *go.etcd.io/etcd/.../v3*. For example, If the bug is located in etcd folder *tests/integration/snapshot*, then the module path would be *go.etcd.io/etcd/tests/v3/integration/snapshot*.

``` bash
cd /playground
git clone https://github.com/etcd-io/etcd.git
cd /playground/etcd

# checkout the specific buggy version

GO111MODULE=on GCatch -mod -mod-abs-path=/playground/etcd -mod-module-path=module_path -compile-error
```

For testing etcd: 
All etcd packages start with *go.etcd.io/etcd/.../v3*. For example, if the bug is located in etcd folder *tests/integration/snapshot*, then the module path would be *go.etcd.io/etcd/tests/v3/integration/snapshot*.

``` bash
cd /playground
git clone https://github.com/etcd-io/etcd.git
cd /playground/etcd

# checkout the specific buggy version

GO111MODULE=on GCatch -mod -mod-abs-path=/playground/etcd -mod-module-path=module_path -compile-error
```

For testing Kubernetes:
All Kubernetes packages start with *k8s.io/kubernetes*. For example, if the bug is located in Kubernets folder *pkg/kubelet/nodeshutdown/systemd*, then the module path is: *k8s.io/kubernetes/pkg/kubelet/nodeshutdown/systemd*

``` bash
cd /playground
git clone https://github.com/kubernetes/kubernetes.git
cd /playground/kubernetes

# checkout the specific buggy version

GO111MODULE=on GCatch -mod -mod-abs-path=/playground/kubernetes -mod-module-path=module_path  -compile-error
```

For testing Prometheus:
All Prometheus packages start with *github.com/prometheus/prometheus*. For example, if the bug is located in Prometheus folder *storage/remote*, then the module path is: *github.com/prometheus/prometheus/storage/remote*

``` bash
cd /playground
https://github.com/prometheus/prometheus.git
cd /playground/prometheus

# checkout the specific buggy version

GO111MODULE=on GCatch -mod -mod-abs-path=/playground/prometheus -mod-module-path=module_path  -compile-error
```

For testing go-Ethereum:
All go-Ethereum packages start with *github.com/ethereum/go-ethereum/*. For example, if the bug is located in Prometheus folder *core*, then the module path is: *github.com/ethereum/go-ethereum/core*

``` bash
cd /playground
git clone https://github.com/ethereum/go-ethereum.git
cd /playground/go-ethereum

# checkout the specific buggy version

GO111MODULE=on GCatch -mod -mod-abs-path=/playground/go-ethereum -mod-module-path=module_path  -compile-error
```

For testing tidb:
For tidb module names, most tidb packages start with *github.com/pingcap/tidb*. For example, if the bug is located in tidb folder *./ddl/*, then the module path is: *github.com/pingcap/tidb/ddl*. However, for bugs reported in *badger*, the module path is: *github.com/pingcap/badger*

``` bash
cd /playground
git clone https://github.com/ethereum/go-ethereum.git
cd /playground/go-ethereum

# checkout the specific buggy version

GO111MODULE=on GCatch -mod -mod-abs-path=/playground/tidb -mod-module-path=module_path  -compile-error
```

For testing Moby(Docker):
Docker fuzzing uses a slightly different routine. The Docker code must be stored in path */go/src/github.com/*. 

``` bash
cd /go/src/github.com
mkdir -p docker
cd docker
git clone https://github.com/moby/moby.git
mv moby docker
GO111MODULE=off GCatch -path=/go/src/github.com/docker/docker -include=github.com/docker/docker -r -compile-error
```


If any bugs found from any programs above, GCatch would output a bug report similar to the following format. 

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




