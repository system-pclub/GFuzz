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

git clone https://github.com/etcd-io/etcd.git 
cd etcd && git checkout bbe1e78e6242a57d54c4b96d8c49ea1e094c3cbb && cd ..

git clone https://github.com/grpc/grpc-go.git
cd grpc-go && git checkout 574137db7de3c10e010d5023626169f13540cef1 && cd ..

git clone https://github.com/ethereum/go-ethereum.git
cd go-ethereum && git checkout d3e3a460ec947c9e1e963d1a35f887d95f23f99d && cd ..

git clone https://github.com/prometheus/prometheus.git
cd prometheus && git checkout ee7e0071d10e795cb78a8e50764e45cbcaf8c29a && cd ..

git clone https://github.com/pingcap/tidb.git 
cd tidb && git checkout b8107d7e51d26a400d8ff0d6a64cf307e4485f19 && cd ..

git clone https://github.com/kubernetes/kubernetes.git 
cd kubernetes && git checkout ebc87c39d3f453cc3d93c968e0fff2228c37b653 && cd ..

# We cannot simply clone moby since it does not support GO111MODULE
# It implies that dependencies of this project would require special cares
# git clone https://github.com/moby/moby.git
# cd moby && git checkout 91dc595e96483184e090d653870eb95d95f96904 && cd ..