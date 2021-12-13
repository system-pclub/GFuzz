#!/bin/bash -e


COMMIT=$1

DEST_DIR=./repos


if [ -z "$COMMIT" ]
then
    echo "COMMIT_HASH is missing"
    echo "Usage: ./fuzz-etcd.sh <COMMIT_HASH>"
    exit 1
fi

mkdir -p $DEST_DIR
cd $DEST_DIR

git clone https://github.com/etcd-io/etcd.git
cd etcd && git checkout $COMMIT && cd ..

#cd "$(dirname "$0")"/..

cd ..

# This script will automatically generate the corresponding test binary files(with and without instrumented go runtime) to the ./tmp/builder

docker build -f docker/builder/Dockerfile -t gfuzzbuilder:latest .


docker run --rm -it \
-v $(pwd)/tmp/builder:/builder \
-v $(pwd)/tmp/pkgmod:/go/pkg/mod \
gfuzzbuilder:latest

