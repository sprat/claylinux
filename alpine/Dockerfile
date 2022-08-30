# syntax=docker/dockerfile:1.4
ARG VERSION=latest
ARG KERNEL_FLAVOR=lts

# =========================================================
FROM alpine:$VERSION AS base

# =========================================================
FROM base AS kernel
ARG KERNEL_FLAVOR
RUN apk add --no-cache linux-$KERNEL_FLAVOR intel-ucode amd-ucode
RUN mkdir -p /out/boot /out/lib \
&&  mv /boot/vmlinuz-$KERNEL_FLAVOR /out/boot/kernel \
&&  mv /boot/config-$KERNEL_FLAVOR /out/boot/config \
&&  mv /boot/intel-ucode.img /boot/amd-ucode.img /out/boot \
&&  mv /lib/modules /out/lib

# ========================================
FROM base AS system

# by default, an alpine image comes with the following packages (at least in the 3.16 branch):
# alpine-baselayout, alpine-baselayout-data, alpine-keys, apk-tools, busybox, ca-certificates-bundle, libc-utils,
# libcrypto1.1, libssl1.1, musl, musl-utils, scanelf, ssl_client, zlib
# ... and a default ISO install comes with alpine-base at least, which adds:
# alpine-conf, busybox-initscripts, busybox-suid, openrc
RUN apk add --no-cache alpine-base

# enable some services
# RUN rc_add() { for svc in $2; do rc-update add $svc $1; done; } \
#  && rc_add sysinit "devfs dmesg mdev hwdrivers" \
#  && rc_add boot "hwclock modules sysctl hostname bootmisc syslog loadkmap networking urandom" \
#  && rc_add default "acpid ntpd" \
#  && rc_add shutdown "mount-ro killprocs savecache"

# copy the system configuration files
# COPY system/ /

# copy the kernel & its modules
COPY --from=kernel /out /
