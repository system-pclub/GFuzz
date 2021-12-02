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

```
git clone https://github.com/system-pclub/GFuzz.git

cd GFuzz

git checkout asplos-artifact
```


Items (2), (3), (4) and (5) are released using a Google Sheet file "asplos-710-artifact" 
(https://docs.google.com/spreadsheets/d/1tLcgsfYlll0g20KMYgDKkAtwZtk426dMSUZ6SvXk04s/edit#gid=0). 
All columns and tabs discussed later are in the Google Sheet file, unless otherwise specified. 










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

