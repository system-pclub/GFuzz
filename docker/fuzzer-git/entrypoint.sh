#!/bin/bash -e
cd "$(dirname "$0")"

OUT_DIR=$1
shift 1

# Start fuzzing
/gfuzz/bin/fuzzer --gomod /fuzz/target --out $OUT_DIR $@
