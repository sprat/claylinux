# syntax=docker/dockerfile:1.7.0
FROM alpine:3.19.1 AS image
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
RUN \
mkdir -p /usr/share/claylinux; \
apk fetch --quiet --stdout --no-cache busybox-static | tar -zxO bin/busybox.static >/usr/share/claylinux/busybox.static; \
chmod +x /usr/share/claylinux/busybox.static
COPY init.sh /usr/share/claylinux/init
COPY build-image.sh /usr/bin/build-image
WORKDIR /out
ENTRYPOINT ["build-image"]
