# syntax=docker/dockerfile:1.12.1
# Generate the OS image
# hadolint ignore=DL3006
FROM claylinux/builder:latest AS build
ARG FORMAT
RUN --mount=from=claylinux/alpine-virt:latest,target=/system build-image --format "$FORMAT"

# =========================================================
# Generate a qemu image running our custom OS image
FROM alpine:3.21.0 AS emulator
RUN apk add --no-cache bash qemu-system-x86_64 ovmf
COPY entrypoint.sh /entrypoint
ENTRYPOINT ["/entrypoint"]
COPY --from=build /out /images
ARG FORMAT
ENV FORMAT="$FORMAT"
