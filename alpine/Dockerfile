# syntax=docker/dockerfile:1.5
ARG VERSION=latest

# =========================================================
FROM alpine:$VERSION AS kernel
ARG FLAVOR=lts
RUN apk add --no-cache --no-scripts linux-$FLAVOR linux-firmware-none \
&& for f in /boot/*-$FLAVOR; do mv $f ${f%-$FLAVOR}; done \
&& mkdir -p /out/boot /out/lib \
&& mv /boot /out \
&& mv /lib/modules /out/lib

# =========================================================
FROM alpine:$VERSION AS base
COPY --from=kernel /out /
RUN apk add --no-cache alpine-base

# Generate the default syslinux configuration
ARG CONSOLE=ttyS0
RUN CMDLINE="console=$CONSOLE" \
&& printf "DEFAULT Linux\nLABEL Linux\n  KERNEL vmlinuz\n  INITRD initramfs\n  APPEND $CMDLINE\n" >/boot/syslinux.cfg

# Setup the system
ARG HOSTNAME="claylinux"
# modloop in sysinit?
ARG SYSINIT_SERVICES="devfs dmesg mdev hwdrivers"
ARG BOOT_SERVICES="modules sysctl hostname bootmisc syslog loadkmap networking hwclock"
ARG DEFAULT_SERVICES="firstboot acpid ntpd"
ARG SHUTDOWN_SERVICES="mount-ro killprocs savecache"

RUN echo "$HOSTNAME" >/etc/hostname \
&& printf "auto lo\niface lo inet loopback\n\nauto eth0\niface eth0 inet dhcp\n" >/etc/network/interfaces \
&& passwd -d root \
&& rc_add() { local runlevel=$1; shift; for service; do rc-update add $service $runlevel; done; } \
&& rc_add sysinit $SYSINIT_SERVICES \
&& rc_add boot $BOOT_SERVICES \
&& rc_add default $DEFAULT_SERVICES \
&& rc_add shutdown $SHUTDOWN_SERVICES

# uncomment ttyS0 line
RUN sed -i '/^#ttyS0/s/^#//g' /etc/inittab
# RUN sed -i '/^tty/s/^/#/g' /etc/inittab && sed -i '/^#'${CONSOLE}'/s/^#//g' /etc/inittab

COPY ./init /
