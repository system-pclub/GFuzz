#!/bin/bash -e
cd "$(dirname "$0")"/.. 

docker build -f dev.Dockerfile -t gfuzz:dev .

docker run -it --rm -v $(pwd):/gfuzz gfuzz:dev bash