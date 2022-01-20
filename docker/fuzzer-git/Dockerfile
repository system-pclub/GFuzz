FROM golang:1.16.4

ARG GIT_URL
ARG GIT_COMMIT

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
WORKDIR /gfuzz

# copy source files to docker
COPY patch ./patch
COPY scripts ./scripts
COPY pkg ./pkg
COPY cmd ./cmd
COPY docker/fuzzer-git/entrypoint.sh docker/fuzzer-git/entrypoint.sh
COPY go.mod go.sum Makefile VERSION ./
RUN make tidy
RUN BUILD=docker make

RUN git clone ${GIT_URL} /fuzz/target
RUN cd /fuzz/target \
&& git checkout ${GIT_COMMIT} \
&& cd ..

RUN /gfuzz/bin/inst --check-syntax-err --recover-syntax-err --dir=/fuzz/target

# patch golang runtime in the container
RUN chmod +x scripts/patch.sh \
&& ./scripts/patch.sh

RUN chmod +x docker/fuzzer-git/entrypoint.sh

ENTRYPOINT [ "/gfuzz/docker/fuzzer-git/entrypoint.sh" ]