# syntax=docker/dockerfile:1.5
ARG VERSION=latest

FROM alpine:$VERSION AS base

# TODO: maybe we shouldn't install the packages to avoid adding more files than necessary
ARG FLAVOR=lts
RUN apk add --no-cache linux-$FLAVOR linux-firmware-none \
&& rm /boot/initramfs-$FLAVOR
# TODO: ucode should be optional
RUN apk add --no-cache intel-ucode amd-ucode

# Generate the default syslinux configuration
# add root=/dev/ram0 on the cmdline?
ARG CONSOLE=ttyS0
RUN KERNEL="vmlinuz-$FLAVOR" \
&& INITRD="intel-ucode.img,amd-ucode.img,initramfs" \
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
