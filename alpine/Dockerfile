# syntax=docker/dockerfile:1.5
ARG VERSION=latest

# =========================================================
FROM alpine:$VERSION AS kernel
ARG FLAVOR=lts
RUN apk add --no-cache --no-scripts linux-$FLAVOR linux-firmware-none \
&& printf "DEFAULT Linux\nLABEL Linux\nKERNEL vmlinuz-$FLAVOR\nINITRD initramfs\n" >/boot/syslinux.cfg \
&& mkdir -p /out/boot /out/lib \
&& mv /boot /out \
&& mv /lib/modules /out/lib

# =========================================================
FROM alpine:$VERSION AS base
COPY --from=kernel /out /
RUN apk add --no-cache alpine-base
COPY ./init /

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
