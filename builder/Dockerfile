# syntax=docker/dockerfile:1
# build an "init" image with our init script & busybox that will allow the kernel to perform an init of the target system
FROM alpine:latest AS init-image
SHELL ["/bin/ash", "-euxo", "pipefail", "-c"]
WORKDIR /src
RUN apk fetch --quiet --stdout --no-cache busybox-static | tar -zx --exclude '.*' && mv bin/busybox.static ./busybox
COPY init.sh ./init
RUN mkdir /out && printf "init\nbusybox" | cpio -oH newc >/out/init.img

# =========================================================
# create the builder image whose purpose is to generate OS images
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
COPY --from=init-image /out/ /usr/share/claylinux/
COPY build-image.sh /usr/local/bin/build-image
WORKDIR /out
ENTRYPOINT ["build-image"]
