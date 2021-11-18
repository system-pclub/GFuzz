#!/bin/bash -e
cd "$(dirname "$0")"

TARGET_GO_MOD_DIR=$1
OUTPUT_DIR=$2

# Start fuzzing
/gfuzz/bin/fuzz -goModDir $TARGET_GO_MOD_DIR -outputDir $OUTPUT_DIR $@
