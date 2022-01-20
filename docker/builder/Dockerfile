FROM golang:1.16.4

RUN apt update && apt install -y python3
RUN apt-get update && apt-get install -y --no-install-recommends \
		build-essential \
		curl \
		cmake \
		gcc \
		git \
		libapparmor-dev \
		libbtrfs-dev \
		libdevmapper-dev \
		libseccomp-dev \
		ca-certificates \
		e2fsprogs \
		iptables \
		pkg-config \
		pigz \
		procps \
		xfsprogs \
		xz-utils \
		\
		aufs-tools \
		vim-common \
	&& rm -rf /var/lib/apt/lists/*

WORKDIR /repos

COPY benchmark/repos/etcd ./etcd
COPY benchmark/repos/go-ethereum ./go-ethereum
COPY benchmark/repos/grpc-go ./grpc-go
COPY benchmark/repos/grpc-go ./grpc-go-goleak
COPY benchmark/repos/kubernetes ./kubernetes
COPY benchmark/repos/prometheus ./prometheus
COPY benchmark/repos/tidb ./tidb

COPY benchmark/repos/moby /go/src/github.com/docker/docker


# override leakcheck.go to prevent sideeffect of leakcheck to the benchmark
COPY docker/builder/leakcheck.go /repos/grpc-go/internal/leakcheck/leakcheck.go
COPY docker/builder/leakcheck-allow-oracle.go /repos/grpc-go-goleak/internal/leakcheck/leakcheck.go

# avoid gRPC testsuite
RUN cd grpc-go && grep -rl 'func (s) Test' ./ | xargs sed -i 's/func (s)/func/g' && cd ..
# Don't replace since goleak only will be called from testsuite way
# RUN cd grpc-go-goleak && grep -rl 'func (s) Test' ./ | xargs sed -i 's/func (s)/func/g' && cd ..


WORKDIR /gfuzz

# copy source files to docker
COPY scripts ./scripts
COPY patch ./patch
COPY pkg ./pkg
COPY cmd ./cmd
COPY docker/builder/entrypoint.sh docker/builder/entrypoint.sh
COPY go.mod go.sum Makefile VERSION ./
RUN make tidy
RUN BUILD=docker make
RUN chmod +x docker/builder/entrypoint.sh scripts/gen-testbins.sh
ENTRYPOINT [ "/gfuzz/docker/builder/entrypoint.sh" ]