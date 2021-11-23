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

COPY scripts /gfuzz/scripts

WORKDIR /repos

RUN /gfuzz/scripts/clone-repos.sh .

# override leakcheck.go to prevent sideeffect of leakcheck to the benchmark
COPY docker/builder/leakcheck.go /repos/grpc-go/internal/leakcheck/leakcheck.go

# avoid gRPC testsuite
RUN cd /grpc-go && grep -rl 'func (s) Test' ./ | xargs sed -i 's/func (s)/func/g'

WORKDIR /gfuzz

# copy source files to docker
COPY patch ./patch

COPY pkg ./pkg
COPY cmd ./cmd
COPY docker/builder/entrypoint.sh docker/builder/entrypoint.sh
COPY go.mod go.sum Makefile VERSION ./
RUN make tidy
RUN BUILD=docker make
RUN chmod +x docker/builder/entrypoint.sh scripts/gen-testbins.sh
ENTRYPOINT [ "/gfuzz/docker/builder/entrypoint.sh" ]