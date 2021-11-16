#!/bin/bash -e
cd "$(dirname "$0")"/..

TARGET_GOMOD_DIR="_examples/fuzz/hello"
OUT_DIR="tmp/hello"
INST_STATS=${OUT_DIR}/stats

mkdir -p ${OUT_DIR}

bin/inst --dir=${TARGET_GOMOD_DIR} --statsOut=${INST_STATS}
rm -rf $OUT_DIR
bin/fuzzer --goModDir=${TARGET_GOMOD_DIR} --outputDir=${OUT_DIR} --instStats=${INST_STATS}