# syntax=docker/dockerfile:1
FROM alpine:latest AS builder

# add the necessary tools to generate the OS
RUN set -euo pipefail; \
PACKAGES="bash dosfstools file mtools pigz sfdisk xorriso zstd xz"; \
[[ "$(uname -m)" == x86* ]] && PACKAGES="$PACKAGES syslinux"; \
apk add --no-cache $PACKAGES

# copy the build script
COPY build-image /usr/local/bin/

WORKDIR /out
ENTRYPOINT ["build-image"]
