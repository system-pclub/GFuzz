inst:
	go build -o bin/inst gfuzz/cmd/inst

test:
	go test -v gfuzz/pkg/...
