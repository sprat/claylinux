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
initfs=/tmp/initfs; \
data=/usr/local/share/claylinux; \
mkdir -p "$initfs" "$data"; \
apk fetch --quiet --stdout --no-cache busybox-static | tar -zx bin/busybox.static >"$initfs"/busybox; \
cp /src/init "$initfs"/init; \
find "$initfs" -printf '%P\n' | grep -v '^$' | cpio --quiet -D "$initfs" -oH newc >"$data"/init.img; \
rm -rf "$initfs"
COPY build-image.sh /usr/local/bin/build-image
WORKDIR /out
ENTRYPOINT ["build-image"]
