VERSION=`cat VERSION`
BUILD=`git log -1 --format=%h`-`date '+%s'`

FUZZER_LD_FLAGS="-X main.Version=${VERSION} -X main.Build=${BUILD}"

inst:
	go mod tidy
	go build -o bin/inst gfuzz/cmd/inst

fuzzer:
	go build -o bin/fuzzer -ldflags $(FUZZER_LD_FLAGS) gfuzz/cmd/fuzzer
test:
	go test -v gfuzz/pkg/...
