# syntax=docker/dockerfile:1.12.1
ARG VERSION=3.21.0

# =========================================================
FROM alpine:$VERSION AS overlay
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
FROM scratch AS image
COPY --from=overlay /out /
ENTRYPOINT ["/bin/sh"]
