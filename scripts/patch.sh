#!/bin/bash -e
cd "$(dirname "$0")"/.. 

GOROOT=$(go env GOROOT)
GOROOT_SRC=$GOROOT/src
RUNTIME=$GOROOT_SRC/runtime

cp runtime/my* $RUNTIME
cp runtime/select.go $RUNTIME/select.go
cp runtime/chan.go $RUNTIME/chan.go
cp runtime/runtime2.go $RUNTIME/runtime2.go
cp runtime/proc.go $RUNTIME/proc.go
cp -r goFuzz/gooracle $GOROOT_SRC
cp -r time $GOROOT_SRC
cp -r sync $GOROOT_SRC
cp reflect/value.go $GOROOT_SRC/reflect/value.go