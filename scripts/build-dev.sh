#!/bin/bash -e
cd "$(dirname "$0")"/.. 

docker build -f docker/dev/Dockerfile -t gfuzz:dev .