#!/bin/bash -e
cd "$(dirname "$0")"

OUTPUT_DIR=$1

# Start fuzzing
/gfuzz/bin/fuzz -goModDir /fuzz/target -outputDir $OUTPUT_DIR $@
