#!/bin/bash -xe
cd "$(dirname "$0")"

OUT_DIR=$1
shift 1

/gfuzz/bin/fuzzer \
--out $OUT_DIR \
"$@"


