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

# gRPC code
RUN git clone https://github.com/grpc/grpc-go.git
RUN cd grpc-go && git checkout 9280052d36656451dd7568a18a836c2a74edaf6c

# override leakcheck.go to prevent sideeffect of leakcheck to the benchmark
COPY docker/builder/leakcheck.go internal/leakcheck/leakcheck.go

# avoid testsuite
RUN grep -rl 'func (s) Test' ./ | xargs sed -i 's/func (s)/func/g'

WORKDIR /gfuzz

# copy source files to docker
COPY patch ./patch
COPY scripts ./scripts
COPY pkg ./pkg
COPY cmd ./cmd
COPY docker/builder/entrypoint.sh docker/builder/entrypoint.sh
COPY go.mod go.sum Makefile VERSION ./
RUN make tidy
RUN BUILD=docker make
RUN chmod +x docker/builder/entrypoint.sh scripts/gen-testbins.sh
ENTRYPOINT [ "/gfuzz/docker/builder/entrypoint.sh" ]