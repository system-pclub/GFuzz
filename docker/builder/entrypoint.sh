#!/bin/bash -x
OUTPUT_DIR=/builder

# build native part
# echo "Building etcd"
# cd /repos/etcd
# /gfuzz/scripts/gen-testbins.sh $OUTPUT_DIR/etcd/native

# cd /repos/etcd/api
# /gfuzz/scripts/gen-testbins.sh $OUTPUT_DIR/etcd/native

# cd /repos/etcd/client/v2
# /gfuzz/scripts/gen-testbins.sh $OUTPUT_DIR/etcd/native

# cd /repos/etcd/client/v3
# /gfuzz/scripts/gen-testbins.sh $OUTPUT_DIR/etcd/native

# cd /repos/etcd/etcdctl
# /gfuzz/scripts/gen-testbins.sh $OUTPUT_DIR/etcd/native

# cd /repos/etcd/pkg
# /gfuzz/scripts/gen-testbins.sh $OUTPUT_DIR/etcd/native

# cd /repos/etcd/server
# /gfuzz/scripts/gen-testbins.sh $OUTPUT_DIR/etcd/native

# cd /repos/etcd/tests
# /gfuzz/scripts/gen-testbins.sh $OUTPUT_DIR/etcd/native

# echo "Building go-ethereum"
# cd /repos/go-ethereum
# /gfuzz/scripts/gen-testbins.sh $OUTPUT_DIR/go-ethereum/native

# echo "Building grpc-go"
# cd /repos/grpc-go
# /gfuzz/scripts/gen-testbins.sh $OUTPUT_DIR/grpc-go/native

# echo "Building kubernetes"
# cd /repos/kubernetes
# /gfuzz/scripts/gen-testbins.sh $OUTPUT_DIR/kubernetes/native

# echo "Building prometheus"
# cd /repos/prometheus
# /gfuzz/scripts/gen-testbins.sh $OUTPUT_DIR/prometheus/native

# echo "Building tidb"
# cd /repos/tidb
# /gfuzz/scripts/gen-testbins.sh $OUTPUT_DIR/tidb/native

# echo "Building moby"
# export GO111MODULE=off
# pkg_list=$(go list github.com/docker/docker/... | grep -vE "(integration)")

# for pkg in $pkg_list
# do
#     echo "generating test bin for $pkg"
#     name=$(echo "$pkg" | sed "s/\//-/g")
#     go test -c -o $OUTPUT_DIR/moby/native/$name.test $pkg
# done
# export GO111MODULE=on

# instrument runtime, code and do instrumentation part
/gfuzz/scripts/patch.sh


# echo "Building etcd"
cd /repos/etcd
/gfuzz/bin/inst --dir . --check-syntax-err --recover-syntax-err

/gfuzz/scripts/gen-testbins.sh $OUTPUT_DIR/etcd/inst

cd /repos/etcd/api
/gfuzz/scripts/gen-testbins.sh $OUTPUT_DIR/etcd/inst

cd /repos/etcd/client/v2
/gfuzz/scripts/gen-testbins.sh $OUTPUT_DIR/etcd/inst

cd /repos/etcd/client/v3
/gfuzz/scripts/gen-testbins.sh $OUTPUT_DIR/etcd/inst

cd /repos/etcd/etcdctl
/gfuzz/scripts/gen-testbins.sh $OUTPUT_DIR/etcd/inst

cd /repos/etcd/pkg
/gfuzz/scripts/gen-testbins.sh $OUTPUT_DIR/etcd/inst

cd /repos/etcd/server
/gfuzz/scripts/gen-testbins.sh $OUTPUT_DIR/etcd/inst

cd /repos/etcd/tests
/gfuzz/scripts/gen-testbins.sh $OUTPUT_DIR/etcd/inst

echo "Building go-ethereum"
cd /repos/go-ethereum
/gfuzz/bin/inst --dir . --check-syntax-err --recover-syntax-err
/gfuzz/scripts/gen-testbins.sh $OUTPUT_DIR/go-ethereum/inst

echo "Building grpc-go"
cd /repos/grpc-go
/gfuzz/bin/inst --dir . --check-syntax-err --recover-syntax-err
/gfuzz/scripts/gen-testbins.sh $OUTPUT_DIR/grpc-go/inst

echo "Building kubernetes"
cd /repos/kubernetes
/gfuzz/bin/inst --dir . --check-syntax-err --recover-syntax-err
/gfuzz/scripts/gen-testbins.sh $OUTPUT_DIR/kubernetes/inst

echo "Building prometheus"
cd /repos/prometheus
/gfuzz/bin/inst --dir . --check-syntax-err --recover-syntax-err
/gfuzz/scripts/gen-testbins.sh $OUTPUT_DIR/prometheus/inst

echo "Building tidb"
cd /repos/tidb
/gfuzz/bin/inst --dir . --check-syntax-err --recover-syntax-err
/gfuzz/scripts/gen-testbins.sh $OUTPUT_DIR/tidb/inst

echo "Building moby"

/gfuzz/bin/inst --dir /go/src/github.com/docker/docker --check-syntax-err --recover-syntax-err

export GO111MODULE=off
pkg_list=$(go list github.com/docker/docker/... | grep -vE "(integration)")

for pkg in $pkg_list
do
    echo "generating test bin for $pkg"
    name=$(echo "$pkg" | sed "s/\//-/g")
    go test -c -o $OUTPUT_DIR/moby/inst/$name.test $pkg
done
export GO111MODULE=on