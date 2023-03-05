ARG VERSION=latest

# =========================================================
FROM alpine:$VERSION AS overlay

# install the alpine-base package (for the user space)
# TODO: or a subset of the packages?
# Note that we need to delete special files because they are transformed to regular files, which causes various
# problems. See https://github.com/moby/buildkit/issues/2914 and https://github.com/docker/buildx/issues/1027
RUN set -euxo pipefail; \
mkdir -p /out/etc; \
cp -R /etc/apk/ /out/etc/apk/; \
apk add --no-cache --root /out --initdb alpine-base; \
rm -rf /out/dev/*

# install the kernel (without installing the dependencies)
ARG FLAVOR=lts
RUN set -euxo pipefail; \
apk fetch --no-cache --quiet --stdout linux-$FLAVOR | tar -xz --directory=/out --exclude=.*; \
mv /out/boot/vmlinuz-$FLAVOR /out/boot/kernel; \
mv /out/boot/config-$FLAVOR /out/boot/config; \
rm /out/boot/System.map-$FLAVOR; \
printf "DEFAULT Linux\nLABEL Linux\nKERNEL kernel\nINITRD initramfs\n" >/out/boot/syslinux.cfg; \
if [[ "$FLAVOR" = "virt" ]]; then echo "APPEND console=ttyS0" >>/out/boot/syslinux.cfg; sed -i '/^#ttyS0/s/^#//g' /out/etc/inittab; fi

# copy our custom init script
COPY ./init /out

# setup the system
ARG HOSTNAME="claylinux"
ARG SYSINIT_SERVICES="devfs dmesg mdev hwdrivers"
ARG BOOT_SERVICES="modules sysctl hostname bootmisc syslog loadkmap networking hwclock"
ARG DEFAULT_SERVICES="firstboot acpid ntpd"
ARG SHUTDOWN_SERVICES="mount-ro killprocs savecache"

RUN set -euxo pipefail; \
echo "$HOSTNAME" >/out/etc/hostname; \
printf "127.0.1.1\t$HOSTNAME\n" >>/out/etc/hosts; \
printf "auto lo\niface lo inet loopback\n\nauto eth0\niface eth0 inet dhcp\n" >/out/etc/network/interfaces; \
rc_add() { local lvl=$1; shift; for svc; do ln -s /etc/init.d/$svc /out/etc/runlevels/$lvl/$svc; done; }; \
rc_add sysinit $SYSINIT_SERVICES; \
rc_add boot $BOOT_SERVICES; \
rc_add default $DEFAULT_SERVICES; \
rc_add shutdown $SHUTDOWN_SERVICES;

# =========================================================
FROM scratch AS final
COPY --from=overlay /out /
ENTRYPOINT ["/bin/sh"]
