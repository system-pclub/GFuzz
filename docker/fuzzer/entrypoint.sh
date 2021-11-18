#!/bin/bash -e
cd "$(dirname "$0")"

TARGET_GO_MOD_DIR=$1
OUTPUT_DIR=$2

# Start fuzzing
/gfuzz/bin/fuzzer --gomod $TARGET_GO_MOD_DIR --out $OUTPUT_DIR $@
