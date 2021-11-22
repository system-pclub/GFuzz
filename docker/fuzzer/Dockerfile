FROM golang:1.16.4
RUN apt update && apt install -y python3
WORKDIR /gfuzz

# copy source files to docker
COPY patch ./patch
COPY scripts ./scripts
COPY pkg ./pkg
COPY cmd ./cmd
COPY docker/fuzzer/entrypoint.sh docker/fuzzer/entrypoint.sh
COPY go.mod go.sum Makefile VERSION ./
# build inst and fuzzer
RUN make tidy
RUN BUILD=docker make

# patch golang runtime in the container
RUN chmod +x scripts/patch.sh \
&& ./scripts/patch.sh

RUN chmod +x docker/fuzzer/entrypoint.sh

# RUN groupadd gfgroup
# RUN useradd -r -u 1001 -g gfgroup gfuser
# RUN chown gfuser:gfgroup ./scripts/fuzz.sh && chmod +x ./scripts/fuzz.sh
# USER gfuser
# RUN chmod +x ./scripts/fuzz.sh
# ENTRYPOINT [ "scripts/fuzz.sh" ] 

ENTRYPOINT [ "/gfuzz/docker/fuzzer/entrypoint.sh" ]