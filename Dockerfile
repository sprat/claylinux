# syntax=docker/dockerfile:1.4
ARG ALPINE_VERSION=latest
ARG KERNEL_FLAVOR=lts

# =========================================================
FROM alpine:$ALPINE_VERSION AS packager
ARG KERNEL_FLAVOR
RUN apk add --no-cache linux-$KERNEL_FLAVOR intel-ucode amd-ucode
RUN mkdir -p /out/boot /out/lib \
&&  mv /boot/vmlinuz-$KERNEL_FLAVOR /out/boot/kernel \
&&  mv /boot/config-$KERNEL_FLAVOR /out/boot/config \
&&  mv /boot/intel-ucode.img /boot/amd-ucode.img /out/boot \
&&  mv /lib/modules /out/lib

# =========================================================
FROM scratch
COPY --from=packager /out /
