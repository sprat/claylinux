# syntax=docker/dockerfile:1

# Setup our custom OS
# hadolint ignore=DL3006
FROM claylinux/alpine-virt AS system

# =========================================================
# Generate the OS image
# hadolint ignore=DL3006
FROM claylinux/builder AS build
ARG FORMAT
RUN --mount=from=system,target=/system build-image --format $FORMAT

# =========================================================
# Download an EFI firmware for the qemu VM
# hadolint ignore=DL3006
FROM alpine AS ovmf
RUN apk add --no-cache ovmf

# =========================================================
# Save the files from the output directory in a new image that will be exported to disk
FROM scratch AS output
COPY --from=ovmf /usr/share/OVMF/OVMF.fd /
COPY --from=build /out /
