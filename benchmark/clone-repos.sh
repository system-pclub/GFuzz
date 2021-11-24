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
cd etcd && git checkout dae29bb719dd69dc119146fc297a0628fcc1ccf8 && cd ..

git clone https://github.com/grpc/grpc-go.git
cd grpc-go && git checkout 9280052d36656451dd7568a18a836c2a74edaf6c && cd ..

git clone https://github.com/ethereum/go-ethereum.git
cd go-ethereum && git checkout 123e934e72dbbc63281858ec20c29beb6f70d9ba && cd ..

git clone https://github.com/prometheus/prometheus.git
cd prometheus && git checkout e0f1506254688cec85276cc939aeb536a4e029d1 && cd ..

git clone https://github.com/pingcap/tidb.git 
cd tidb && git checkout 7e6690df8e8d5474b1872edbd279bb1b3c510ee5 && cd ..

git clone https://github.com/kubernetes/kubernetes.git 
cd kubernetes && git checkout 97d40890d00acf721ecabb8c9a6fec3b3234b74b && cd ..

git clone https://github.com/moby/moby.git
cd moby && git checkout 791640417b67036bbc7d13597cad55bb5fcead2b && cd ..