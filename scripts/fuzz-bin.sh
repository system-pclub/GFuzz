#!/bin/bash -ex
cd "$(dirname "$0")"/..

# ./scripts/fuzz-bin.sh $(pwd)/tmp/builder/grpc-go/inst/google.golang.org-grpc-internal-transport.test ~/gfuzz/tmp/o1 --func TestHandlerTransport_HandleStreams_MultiWriteStatus

BIN=$1
OUT_DIR=$2
shift 2

docker build -f docker/fuzzer-bin/Dockerfile -t gfuzzbin:latest .

docker run -it --rm \
-v $BIN:/fuzz/target \
-v $OUT_DIR:/fuzz/output \
-v $(pwd)/tmp/pkgmod:/go/pkg/mod \
gfuzzbin:latest /fuzz/output --bin "/fuzz/target" $@

