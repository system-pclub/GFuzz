#!/bin/bash -e

# Examples:
# $ ./benchmark.sh custom --dir /builder/grpc/native --mode native
# $ ./benchmark.sh custom --dir /builder/grpc/inst --mode inst

docker build -f docker/benchmark/Dockerfile -t gfuzzbenchmark:latest .
docker run -it --rm \
-v $(pwd)/tmp/builder:/builder \
-v $(pwd)/tmp/pkgmod:/go/pkg/mod \
gfuzzbenchmark:latest $@