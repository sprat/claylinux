# syntax=docker/dockerfile:1
# build an "init" image with our init script & busybox that will allow the kernel to start/init the target system
FROM alpine:3 AS init-image
SHELL ["/bin/ash", "-euxo", "pipefail", "-c"]
WORKDIR /src
RUN apk fetch --quiet --stdout --no-cache busybox-static | tar -zx --exclude '.*' && mv bin/busybox.static bin/busybox
COPY init .
RUN mkdir /out && find . -print0 | cpio -0oH newc >/out/init.img

# =========================================================
FROM alpine:3 AS builder
SHELL ["/bin/ash", "-euxo", "pipefail", "-c"]

# add the necessary tools to generate the OS
RUN apk add --no-cache bash binutils dosfstools gummiboot-efistub mtools ovmf pigz qemu-img sfdisk xorriso zstd xz

# copy the init image
COPY --from=init-image /out/ /usr/share/claylinux/

# copy our build script
COPY build-image /usr/local/bin/

# set the working directory
WORKDIR /out

ENTRYPOINT ["build-image"]
