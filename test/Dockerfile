# syntax=docker/dockerfile:1
FROM base AS system

# TODO: make sure the ucodes are used in the syslinux config and/or in the initramfs
# RUN apk add --no-cache intel-ucode amd-ucode
# enable passwordless root login
RUN sed -i 's/root:[*!]:/root::/g' /etc/shadow
# enable the serial console (kernel commandline & inittab)
RUN echo "APPEND console=ttyS0" >>/boot/syslinux.cfg && sed -i '/^#ttyS0/s/^#//g' /etc/inittab;

# =========================================================
FROM builder AS build
RUN --mount=type=bind,from=system,target=/system build-os

# =========================================================
FROM scratch AS output
COPY --from=build /out /
