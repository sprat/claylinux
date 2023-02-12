# syntax=docker/dockerfile:1.5
ARG VERSION=latest

FROM alpine:$VERSION AS base

# This docker image will be used to generate a bootable OS image, so we must add a kernel, an initramfs and an init
# system into it, because a docker image does not contain that normally.
# The official Alpine Linux "profiles" used to build the official images are described in the mkimg.*.sh scripts
# here: https://gitlab.alpinelinux.org/alpine/aports/-/tree/master/scripts
RUN apk add --no-cache alpine-base
ARG FLAVOR=lts
RUN apk add --no-cache linux-$FLAVOR linux-firmware-none
# TODO: ucode should be optional
RUN apk add --no-cache intel-ucode amd-ucode

# Generate the default syslinux configuration
# add root=/dev/ram0 on the cmdline?
ARG CONSOLE=ttyS0
RUN KERNEL="vmlinuz-$FLAVOR" \
 && INITRD="intel-ucode.img,amd-ucode.img,initramfs-$FLAVOR" \
 && CMDLINE="console=$CONSOLE" \
 && printf "DEFAULT Linux\nLABEL Linux\n  KERNEL $KERNEL\n  INITRD $INITRD\n  APPEND $CMDLINE\n" >/boot/syslinux.cfg

# Setup the system
RUN echo "claylinux" >/etc/hostname

# enable some services
# RUN rc_add() { for svc in $2; do rc-update add $svc $1; done; } \
#  && rc_add sysinit "devfs dmesg mdev hwdrivers" \
#  && rc_add boot "hwclock modules sysctl hostname bootmisc syslog loadkmap networking urandom" \
#  && rc_add default "acpid ntpd" \
#  && rc_add shutdown "mount-ro killprocs savecache"

# RUN sed -i '/tty/d' /etc/inittab
COPY ./init /
# COPY ./init /usr/share/mkinitfs/initramfs-init
# RUN KERNEL_VERSION=$(basename $(ls -d /lib/modules/*)) \
# && mkinitfs $KERNEL_VERSION

# echo "Generating the initramfs"
RUN find / -path /boot -prune -o -xdev -print | cpio -oH newc | gzip -9 >/boot/initramfs-$FLAVOR

# ENTRYPOINT ["/sbin/init"]

# ================================================================
FROM base AS system

