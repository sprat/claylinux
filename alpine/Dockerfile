# syntax=docker/dockerfile:1.5
ARG VERSION=latest

# =========================================================
FROM alpine:$VERSION AS overlay

# install the kernel
ARG FLAVOR=lts
RUN mkdir -p /out \
&& apk fetch --no-cache --quiet --stdout linux-$FLAVOR | tar -xz --directory=/out --exclude=.PKGINFO --exclude=.SIGN*

# create the according syslinux configuration
RUN printf "DEFAULT Linux\nLABEL Linux\nKERNEL vmlinuz-$FLAVOR\nINITRD initramfs\n" >/out/boot/syslinux.cfg

# copy our custom init script
COPY ./init /out

# =========================================================
FROM alpine:$VERSION AS base
COPY --from=overlay /out /
RUN apk add --no-cache alpine-base

# Setup the system
ARG HOSTNAME="claylinux"
# TODO: modloop in sysinit?
ARG SYSINIT_SERVICES="devfs dmesg mdev hwdrivers"
ARG BOOT_SERVICES="modules sysctl hostname bootmisc syslog loadkmap networking hwclock"
ARG DEFAULT_SERVICES="firstboot acpid ntpd"
ARG SHUTDOWN_SERVICES="mount-ro killprocs savecache"

# TODO: it would be good to add the hostname in /etc/hosts, but it's complicated because /etc/hosts is overlayed
RUN echo "$HOSTNAME" >/etc/hostname \
&& printf "auto lo\niface lo inet loopback\n\nauto eth0\niface eth0 inet dhcp\n" >/etc/network/interfaces \
&& sed -i 's/root:[*!]:/root::/g' /etc/shadow \
&& rc_add() { local runlevel=$1; shift; for service; do rc-update add $service $runlevel; done; } \
&& rc_add sysinit $SYSINIT_SERVICES \
&& rc_add boot $BOOT_SERVICES \
&& rc_add default $DEFAULT_SERVICES \
&& rc_add shutdown $SHUTDOWN_SERVICES
