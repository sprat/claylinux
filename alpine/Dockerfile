# syntax=docker/dockerfile:1.5
ARG VERSION=latest

FROM alpine:$VERSION AS base

# TODO: maybe we shouldn't install the packages to avoid adding more files than necessary
ARG FLAVOR=lts
RUN apk add --no-cache linux-$FLAVOR linux-firmware-none \
&& rm /boot/initramfs-$FLAVOR \
&& rm /boot/boot

# TODO: ucode should be optional
RUN apk add --no-cache intel-ucode amd-ucode
RUN apk add --no-cache alpine-base

# Generate the default syslinux configuration
ARG CONSOLE=ttyS0
RUN KERNEL="vmlinuz-$FLAVOR" \
&& INITRD="intel-ucode.img,amd-ucode.img,initramfs" \
&& CMDLINE="console=$CONSOLE" \
&& printf "DEFAULT Linux\nLABEL Linux\n  KERNEL $KERNEL\n  INITRD $INITRD\n  APPEND $CMDLINE\n" >/boot/syslinux.cfg

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
