#!/bin/bash -e
OUTPUT_DIR=/builder

# build native part
cd /repos/grpc-go
exclude_paths="(abi)|(fuzzer)|(integration)" /gfuzz/scripts/gen-testbins.sh $OUTPUT_DIR/grpc/native


# instrument runtime, code and do instrumentation part
/gfuzz/scripts/patch.sh
cd /repos/grpc-go
/gfuzz/bin/inst --dir .
exclude_paths="(abi)|(fuzzer)|(integration)" /gfuzz/scripts/gen-testbins.sh $OUTPUT_DIR/grpc/inst