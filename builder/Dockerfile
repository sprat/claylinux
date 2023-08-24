# syntax=docker/dockerfile:1
FROM alpine:latest AS builder
SHELL ["/bin/ash", "-euxo", "pipefail", "-c"]
RUN apk add --no-cache \
bash \
binutils \
coreutils \
cpio \
dosfstools \
findutils \
gummiboot-efistub \
mtools \
pigz \
qemu-img \
sfdisk \
xorriso \
zstd \
xz
RUN --mount=source=init.sh,target=/src/init \
mkdir -p /tmp/busybox /tmp/initfs /usr/local/share/claylinux/; \
cd /tmp/initfs; \
apk fetch --quiet --stdout --no-cache busybox-static | tar -zx -C /tmp bin/busybox.static; \
mv /tmp/bin/busybox.static busybox; \
cp /src/init init; \
find . -print | cpio --quiet -oH newc >/usr/local/share/claylinux/init.img; \
rm -rf /tmp/*
COPY build-image.sh /usr/local/bin/build-image
WORKDIR /out
ENTRYPOINT ["build-image"]
