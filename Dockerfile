# syntax=docker/dockerfile:1.4
ARG ALPINE_VERSION=latest
ARG KERNEL_FLAVOR=lts

FROM alpine:$ALPINE_VERSION AS packager
ARG KERNEL_FLAVOR
RUN apk add --no-cache linux-$KERNEL_FLAVOR \
&& mkdir -p /out/boot /out/lib \
&& mv /boot/vmlinuz-$KERNEL_FLAVOR /out/boot/vmlinuz \
&& mv /boot/config-$KERNEL_FLAVOR /out/boot/config \
&& mv /lib/modules /out/lib

FROM scratch
COPY --from=packager /out /
