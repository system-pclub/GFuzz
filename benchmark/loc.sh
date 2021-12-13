#!/bin/bash -e

DEST_DIR=$1
if [ -z "$DEST_DIR" ]
then
    echo "DEST_DIR is missing"
    echo "Usage: ./clone-repos.sh <DEST_DIR>"
    exit 1
fi

mkdir -p $DEST_DIR
#cd $DEST_DIR

./cloc ./${DEST_DIR}/etcd
./cloc ./${DEST_DIR}/grpc-go
./cloc ./${DEST_DIR}/go-ethereum
./cloc ./${DEST_DIR}/prometheus
./cloc ./${DEST_DIR}/tidb
./cloc ./${DEST_DIR}/kubernetes
./cloc ./${DEST_DIR}/moby
