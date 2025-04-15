# syntax=docker/dockerfile:1.14.0
FROM golang:1.24.2-alpine AS init
WORKDIR /go/src
RUN \
--mount=source=init,target=. \
--mount=type=cache,target=/root/.cache/go-build \
CGO_ENABLED=0 go build -o /go/bin/init -v --ldflags '-s -w -extldflags=-static'

# =========================================================
FROM alpine:3.21.3 AS image
SHELL ["/bin/ash", "-euxo", "pipefail", "-c"]
RUN \
echo "@testing https://dl-cdn.alpinelinux.org/alpine/edge/testing" >>/etc/apk/repositories && \
apk add --no-cache \
bash \
binutils \
coreutils \
cpio \
dosfstools \
findutils \
mtools \
pigz \
qemu-img \
sfdisk \
xorriso \
zstd \
xz \
systemd-efistub@testing
COPY --from=init /go/bin/init /usr/share/claylinux/init
COPY build-image.sh /usr/bin/build-image
WORKDIR /out
ENTRYPOINT ["build-image"]
