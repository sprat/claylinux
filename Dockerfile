# syntax=docker/dockerfile:1.19.0
FROM koalaman/shellcheck-alpine:v0.11.0 AS shellcheck
WORKDIR /src
RUN --mount=type=bind,target=. \
shellcheck --version && find . -name '*.sh' -exec shellcheck {} +

# =========================================================
FROM hadolint/hadolint:v2.14.0-alpine AS hadolint
WORKDIR /src
RUN --mount=type=bind,target=. \
hadolint --version && find . -name Dockerfile -exec hadolint {} +

# =========================================================
FROM toolhippie/yamllint:1.37.1 AS yamllint
WORKDIR /src
RUN --mount=type=bind,target=. \
yamllint -v && yamllint -s -f colored .

# =========================================================
FROM golang:1.25.2-alpine AS init
WORKDIR /go/src
RUN \
--mount=source=init,target=. \
--mount=type=cache,target=/root/.cache/go-build \
CGO_ENABLED=0 go build -o /go/bin/init -v --ldflags '-s -w -extldflags=-static'

# =========================================================
FROM alpine:3.22.2 AS alpine-base

# =========================================================
FROM alpine-base AS imager
SHELL ["/bin/ash", "-euxo", "pipefail", "-c"]
RUN \
echo "@edge-community https://dl-cdn.alpinelinux.org/alpine/edge/community" >>/etc/apk/repositories && \
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
systemd-efistub@edge-community
COPY --from=init /go/bin/init /usr/share/claylinux/init
COPY build-image.sh /usr/bin/build-image
WORKDIR /out
ENTRYPOINT ["build-image"]

# =========================================================
FROM alpine-base AS bootable-alpine-rootfs
SHELL ["/bin/ash", "-euxo", "pipefail", "-c"]

# install the alpine-base package (for the user space)
# TODO: or a subset of the packages?
RUN \
mkdir -p /out/etc; \
cp -R /etc/apk/ /out/etc/apk/; \
apk add --no-cache --root /out --initdb alpine-base

# install the kernel (without installing the dependencies)
ARG FLAVOR=lts
RUN apk fetch --no-cache --quiet --stdout linux-$FLAVOR | tar -xz --directory=/out --exclude=.*

# setup the system
# by default, we remove all the TTYs in inittab, so they should be added back with the INITTAB_TTYS argument
# e.g. you can use "ttyS0::respawn:/sbin/getty -L 0 ttyS0 vt100" to enable login on the serial console
ARG \
HOSTNAME="claylinux" \
CMDLINE="" \
INITTAB_TTYS="tty1::respawn:/sbin/getty 38400 tty1" \
NETWORK_INTERFACES="auto lo\niface lo inet loopback\n\nauto eth0\niface eth0 inet dhcp" \
SYSINIT_SERVICES="devfs dmesg mdev hwdrivers" \
BOOT_SERVICES="modules sysctl hostname bootmisc syslog networking hwclock" \
DEFAULT_SERVICES="acpid ntpd" \
SHUTDOWN_SERVICES="mount-ro killprocs savecache"

RUN \
echo "$HOSTNAME" >/out/etc/hostname; \
mv /out/etc/hosts /out/etc/hosts.target; \
printf "127.0.1.1\t%s\n" "$HOSTNAME" >>/out/etc/hosts.target; \
printf "%b\n" "$CMDLINE" >/out/boot/cmdline; \
sed -i -E '/tty/d;/^#/d;/^$/d' /out/etc/inittab; \
printf "%b\n" "$INITTAB_TTYS" >>/out/etc/inittab; \
printf "%b\n" "$NETWORK_INTERFACES" >/out/etc/network/interfaces; \
rc_add() { for svc in $2; do ln -s "/etc/init.d/$svc" "/out/etc/runlevels/$1/$svc"; done; }; \
rc_add sysinit "$SYSINIT_SERVICES"; \
rc_add boot "$BOOT_SERVICES"; \
rc_add default "$DEFAULT_SERVICES"; \
rc_add shutdown "$SHUTDOWN_SERVICES"

# =========================================================
FROM scratch AS bootable-alpine
COPY --from=bootable-alpine-rootfs /out /
ENTRYPOINT ["/bin/sh"]

# =========================================================
# Prepare a custom test rootfs
FROM bootable-alpine AS test-rootfs
ARG UCODE=none
RUN if [ "$UCODE" != "none" ]; then apk add --no-cache "${UCODE}-ucode"; fi

# =========================================================
# Generate a test OS image from the rootfs we prepared
# hadolint ignore=DL3006
FROM imager AS test
ARG FORMAT=efi
RUN --mount=from=test-rootfs,target=/system build-image --format "$FORMAT"

# =========================================================
# Generate a qemu image running our custom OS image
FROM alpine-base AS emulator
RUN apk add --no-cache bash qemu-system-x86_64 ovmf
COPY emulator.sh /entrypoint
ENTRYPOINT ["/entrypoint"]
COPY --from=test /out /images
ARG FORMAT
ENV FORMAT="$FORMAT"
