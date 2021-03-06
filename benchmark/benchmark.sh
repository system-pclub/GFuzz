#!/bin/bash -e
cd "$(dirname "$0")"/.. 

docker build -f docker/benchmark/Dockerfile -t gfuzzbenchmark:latest .
docker run -it --rm \
-v $(pwd)/tmp/builder:/builder \
-v $(pwd)/tmp/pkgmod:/go/pkg/mod \
gfuzzbenchmark:latest $@