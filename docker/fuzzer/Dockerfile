FROM golang:1.16.4
RUN apt update && apt install -y python3
WORKDIR /gfuzz

# copy source files to docker
COPY patch ./patch
COPY scripts ./scripts
COPY pkg ./pkg
COPY docker/fuzzer/entrypoint.sh docker/fuzzer/entrypoint.sh

# build inst and fuzzer
RUN make all

# patch golang runtime in the container
RUN chmod +x scripts/patch.sh \
&& ./scripts/patch.sh

RUN chmod +x docker/fuzzer/entrypoint.sh

ENTRYPOINT [ "/gfuzz/docker/fuzzer/entrypoint.sh" ]