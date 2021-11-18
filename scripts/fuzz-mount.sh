#!/bin/bash -e
cd "$(dirname "$0")"/..

TARGET_GO_MOD_DIR=$1
OUTPUT_DIR=$2
shift 2

docker build -f docker/fuzzer/Dockerfile -t gfuzz:latest .

docker run --rm -it \
-v $OUTPUT_DIR:/fuzz/output \
-v $TARGET_GO_MOD_DIR:/fuzz/target \
gfuzz:latest /fuzz/target /fuzz/output $@
