ARG VERSION=latest

FROM alpine:$VERSION AS builder

# add the necessary tools to generate the OS
RUN apk add --no-cache bash sfdisk dosfstools mtools rsync syslinux xorriso xz zstd pigz

# copy the build script
COPY build-img /usr/local/bin/

ENTRYPOINT ["build-img"]
