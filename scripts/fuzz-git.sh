#!/bin/bash -e
cd "$(dirname "$0")"

TARGET_GIT=$1
TARGET_GIT_COMMIT=$2
OUTPUT_DIR=$3
shift 3


docker build \
&& --build-arg GIT_URL=$TARGET_GIT
&& --build-arg GIT_COMMIT=$TARGET_GIT_COMMIT
&& -t gfuzzgit:latest .

docker run --rm -it \
-v "${CWD}/tmp/pkgmod":/go/pkg/mod \
-v $OUTPUT_DIR:/fuzz/output \
gfuzzgit:latest /fuzz/target /fuzz/output $@