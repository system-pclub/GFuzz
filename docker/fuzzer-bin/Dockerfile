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
WORKDIR /gfuzz

# copy source files to docker
COPY patch ./patch
COPY scripts ./scripts
COPY pkg ./pkg
COPY cmd ./cmd
COPY docker/fuzzer-bin/entrypoint.sh docker/fuzzer-bin/entrypoint.sh
COPY go.mod go.sum Makefile VERSION ./
# build inst and fuzzer
RUN make tidy
RUN BUILD=docker make

# patch golang runtime in the container
RUN chmod +x scripts/patch.sh \
&& ./scripts/patch.sh

RUN chmod +x docker/fuzzer-bin/entrypoint.sh

# RUN groupadd gfgroup
# RUN useradd -r -u 1001 -g gfgroup gfuser
# RUN chown gfuser:gfgroup ./scripts/fuzz.sh && chmod +x ./scripts/fuzz.sh
# USER gfuser
# RUN chmod +x ./scripts/fuzz.sh
# ENTRYPOINT [ "scripts/fuzz.sh" ] 

ENTRYPOINT [ "/gfuzz/docker/fuzzer-bin/entrypoint.sh" ]