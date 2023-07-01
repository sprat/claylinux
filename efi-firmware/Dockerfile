# syntax=docker/dockerfile:1

# Download an EFI firmware for the qemu VM
FROM alpine:latest AS efi-firmware
SHELL ["/bin/ash", "-euxo", "pipefail", "-c"]
RUN \
mkdir -p /out; \
apk add --no-cache ovmf; \
cp /usr/share/OVMF/OVMF.fd /out

# =========================================================
# Save the files from the output directory in a new image that will be exported to disk
FROM scratch AS output
COPY --from=efi-firmware /out /
