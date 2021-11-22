FROM golang:1.16.4

ARG GIT_URL
ARG GIT_COMMIT

RUN apt update && apt install -y python3
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
&& git checkout GIT_COMMIT \
&& cd ..

RUN /gfuzz/bin/inst --dir=/fuzz/target

# patch golang runtime in the container
RUN chmod +x scripts/patch.sh \
&& ./scripts/patch.sh

RUN chmod +x docker/fuzzer-git/entrypoint.sh

ENTRYPOINT [ "/gfuzz/docker/fuzzer-git/entrypoint.sh" ]