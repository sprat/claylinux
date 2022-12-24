# syntax=docker/dockerfile:1.4
ARG VERSION=latest

FROM alpine:$VERSION AS builder

# add the necessary tools to generate the OS
RUN apk add --no-cache sfdisk dosfstools mtools syslinux xorriso

# copy the build script
COPY build-os.sh /usr/local/bin/build-os

ENTRYPOINT ["build-os"]
