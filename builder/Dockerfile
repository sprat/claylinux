# syntax=docker/dockerfile:1.11.1
FROM golang:1.23.3-alpine AS init
WORKDIR /go/src
RUN \
--mount=source=.,target=. \
--mount=type=cache,target=/root/.cache/go-build \
CGO_ENABLED=0 go build -o /go/bin/init -v --ldflags '-s -w -extldflags=-static'

# =========================================================
FROM alpine:3.20.3 AS image
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
COPY --from=init /go/bin/init /usr/share/claylinux/init
COPY build-image.sh /usr/bin/build-image
WORKDIR /out
ENTRYPOINT ["build-image"]
