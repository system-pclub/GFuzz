FROM golang:1.16.4
RUN apt update && apt install -y python3

WORKDIR /gfuzz
# copy source files to docker
COPY patch ./patch
COPY pkg ./pkg
COPY scripts ./scripts

# Patch golang runtime in the container
RUN chmod +x scripts/patch.sh \
&& ./scripts/patch.sh


WORKDIR /gfuzz