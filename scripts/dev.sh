#!/bin/bash -e
cd "$(dirname "$0")"/.. 

docker build -f docker/dev/Dockerfile -t gfuzz:dev .

docker run -it --rm \
-v $(pwd):/gfuzz \
-v $(pwd)/tmp/pkgmod:/go/pkg/mod \
gfuzz:dev bash