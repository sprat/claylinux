# syntax=docker/dockerfile:1
ARG ALPINE_VERSION=latest
ARG KERNEL_FLAVOR=lts

FROM alpine:$ALPINE_VERSION as builder
ARG KERNEL_FLAVOR
RUN apk add --no-cache linux-$KERNEL_FLAVOR \
&& mkdir -p /out/boot /out/lib \
&& mv /boot/vmlinuz-$KERNEL_FLAVOR /out/boot/vmlinuz \
&& mv /boot/config-$KERNEL_FLAVOR /out/boot/config \
&& mv /lib/modules /out/lib

FROM scratch
COPY --from=builder /out /
