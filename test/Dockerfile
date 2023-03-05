# syntax=docker/dockerfile:1
FROM base AS system

# TODO: make sure the ucodes are used in the syslinux config and/or in the initramfs
# RUN apk add --no-cache intel-ucode amd-ucode
# enable passwordless root login
RUN sed -i 's/root:[*!]:/root::/g' /etc/shadow

# =========================================================
FROM builder AS build
RUN --mount=from=system,target=/system build-img

# =========================================================
FROM scratch AS output
COPY --from=build /out /
