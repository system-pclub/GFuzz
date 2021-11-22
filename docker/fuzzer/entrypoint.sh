#!/bin/bash -e
cd "$(dirname "$0")"

NEED_INST=$1
GOMOD_DIR=$2
OUT_DIR=$3
shift 3

if [ "$NEED_INST" = true ]; then
    /gfuzz/bin/inst --dir $GOMOD_DIR
fi

# Start fuzzing
/gfuzz/bin/fuzzer \
--gomod $GOMOD_DIR \
--out $OUT_DIR \
$@
