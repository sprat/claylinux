# syntax=docker/dockerfile:1
FROM alpine:latest AS builder

# add the necessary tools to generate the OS
RUN apk add --no-cache bash dosfstools mtools pigz sfdisk syslinux xorriso zstd xz

# copy the build script
COPY build-image /usr/local/bin/

WORKDIR /out
ENTRYPOINT ["build-image"]
