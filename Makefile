VERSION=`cat VERSION`

ifndef BUILD
	BUILD=`git log -1 --format=%h`-`date '+%s'`
endif

BIN_FUZZER_LD_FLAGS="-X main.Version=${VERSION} -X main.Build=${BUILD}"
BIN_INST_LD_FLAGS="-X main.Version=${VERSION} -X main.Build=${BUILD}"

.PHONY: all test clean tidy
all: bin/inst bin/fuzzer

tidy:
	go mod tidy

bin/inst: 
	go build -o bin/inst -ldflags $(BIN_INST_LD_FLAGS) gfuzz/cmd/inst

bin/fuzzer: 
	go build -o bin/fuzzer -ldflags $(BIN_FUZZER_LD_FLAGS) gfuzz/cmd/fuzzer

test:
	go test -v gfuzz/pkg/...

clean:
	rm -rf bin
