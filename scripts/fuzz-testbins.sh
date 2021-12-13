#!/bin/bash -e
cd "$(dirname "$0")"/..


TESTBINS_DIR=$1
OUT_DIR=$2
shift 2


docker build -f docker/fuzzer-bin/Dockerfile -t gfuzzbin:latest .

docker run -it \
-v $TESTBINS_DIR:/fuzz/target \
-v $OUT_DIR:/fuzz/output \
-v $(pwd)/tmp/pkgmod:/go/pkg/mod \
gfuzzbin:latest /fuzz/output --bin "/fuzz/target/*" $@
