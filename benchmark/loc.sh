#!/bin/bash -e

DEST_DIR=$1
if [ -z "$DEST_DIR" ]
then
    echo "DEST_DIR is missing"
    echo "Usage: ./clone-repos.sh <DEST_DIR>"
    exit 1
fi

mkdir -p $DEST_DIR
cd $DEST_DIR

cloc ./etcd
cloc ./grpc-go
cloc ./go-ethereum
cloc ./prometheus
cloc ./tidb
cloc ./kubernetes
cloc ./moby