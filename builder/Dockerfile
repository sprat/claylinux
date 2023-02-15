# syntax=docker/dockerfile:1.5
ARG VERSION=latest

FROM alpine:$VERSION AS builder

# add the necessary tools to generate the OS
RUN apk add --no-cache bash sfdisk dosfstools mtools syslinux xorriso xz zstd pigz

# copy the build script
COPY build-os.sh /usr/local/bin/build-os

ENTRYPOINT ["build-os"]
