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

# install the alpine-base package
# TODO: or a subset of the packages?
# Note that we need to delete character devices because they are transformed to regular files
# is exported, which causes various problems later. See https://github.com/moby/buildkit/issues/2914 and
# https://github.com/docker/buildx/issues/1027
RUN mkdir -p /out/etc \
&& cp -R /etc/apk/ /out/etc/apk/ \
&& apk add --no-cache --root /out --initdb alpine-base \
&& rm -rf /out/dev/*

# setup the system
ARG HOSTNAME="claylinux"
# TODO: modloop in sysinit?
ARG SYSINIT_SERVICES="devfs dmesg mdev hwdrivers"
ARG BOOT_SERVICES="modules sysctl hostname bootmisc syslog loadkmap networking hwclock"
ARG DEFAULT_SERVICES="firstboot acpid ntpd"
ARG SHUTDOWN_SERVICES="mount-ro killprocs savecache"
# TODO: it would be good to add the hostname in /etc/hosts, but it's complicated because /etc/hosts is overlayed
RUN echo "$HOSTNAME" >/out/etc/hostname \
&& printf "auto lo\niface lo inet loopback\n\nauto eth0\niface eth0 inet dhcp\n" >/out/etc/network/interfaces \
&& sed -i 's/root:[*!]:/root::/g' /out/etc/shadow \
&& rc_add() { set -eu; local lvl=$1; shift; for svc; do ln -s /etc/init.d/$svc /out/etc/runlevels/$lvl/$svc; done; } \
&& rc_add sysinit $SYSINIT_SERVICES \
&& rc_add boot $BOOT_SERVICES \
&& rc_add default $DEFAULT_SERVICES \
&& rc_add shutdown $SHUTDOWN_SERVICES

# copy our custom init script
COPY ./init /out

# =========================================================
FROM scratch AS final
COPY --from=overlay /out /
ENTRYPOINT ["/bin/sh"]
