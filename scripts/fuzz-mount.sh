#!/bin/bash -e
cd "$(dirname "$0")"/..

GOMOD_DIR=$(realpath $1)
OUT_DIR=$(realpath $2)
shift 2

docker build -f docker/fuzzer/Dockerfile -t gfuzz:latest .

docker run --rm -it \
-v $GOMOD_DIR:/fuzz/target \
-v $OUT_DIR:/fuzz/output \
-v $(pwd)/tmp/pkgmod:/go/pkg/mod \
gfuzz:latest true /fuzz/target /fuzz/output $@
