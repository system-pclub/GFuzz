#!/bin/bash -e
cd "$(dirname "$0")"/.. 

GOROOT=$(go env GOROOT)
GOROOT_SRC=$GOROOT/src
RUNTIME=$GOROOT_SRC/runtime

cp patch/runtime/my* $RUNTIME
cp patch/runtime/select.go $RUNTIME/select.go
cp patch/runtime/chan.go $RUNTIME/chan.go
cp patch/runtime/runtime2.go $RUNTIME/runtime2.go
cp patch/runtime/proc.go $RUNTIME/proc.go

mkdir -p $GOROOT_SRC/gfuzz/pkg
cp -r pkg/oraclert $GOROOT_SRC/gfuzz/pkg
cp -r pkg/selefcm $GOROOT_SRC/gfuzz/pkg
cp -r patch/time $GOROOT_SRC
cp -r patch/sync $GOROOT_SRC
cp patch/reflect/value.go $GOROOT_SRC/reflect/value.go