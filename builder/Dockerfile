# syntax=docker/dockerfile:1
FROM alpine:3 AS builder
SHELL ["/bin/ash", "-euxo", "pipefail", "-c"]

# add the necessary tools to generate the OS
RUN apk add --no-cache bash dosfstools limine-cd limine-deploy mtools pigz sfdisk xorriso zstd xz

# copy the build script
COPY build-image /usr/local/bin/

WORKDIR /out
ENTRYPOINT ["build-image"]
