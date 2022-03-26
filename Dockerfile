# syntax=docker/dockerfile:1.3
ARG ALPINE_VERSION=latest
ARG KERNEL_FLAVOR=lts

FROM alpine:$ALPINE_VERSION as builder
ARG KERNEL_FLAVOR
RUN apk add --no-cache linux-$KERNEL_FLAVOR

FROM scratch
ARG KERNEL_FLAVOR
COPY --from=builder /boot/vmlinuz-$KERNEL_FLAVOR /boot/vmlinuz
COPY --from=builder /boot/config-$KERNEL_FLAVOR /boot/config
COPY --from=builder /lib/modules /lib/modules
