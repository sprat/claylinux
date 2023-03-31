# syntax=docker/dockerfile:1
FROM alpine:3 AS builder
SHELL ["/bin/ash", "-euxo", "pipefail", "-c"]

# add the necessary tools to generate the OS
# syslinux is not available on non-PC platforms
RUN \
apk add --no-cache bash dosfstools file mtools pigz sfdisk xorriso zstd xz; \
apk add --no-cache syslinux || true

# copy the build script
COPY build-image /usr/local/bin/

WORKDIR /out
ENTRYPOINT ["build-image"]
