# syntax=docker/dockerfile:1.6.0
# Setup a custom OS
# hadolint ignore=DL3006
FROM alpine-virt AS system
RUN apk add --no-cache intel-ucode

# =========================================================
# Generate the OS image
# hadolint ignore=DL3006
FROM builder AS test
ARG FORMAT=efi
RUN --mount=from=system,target=/system build-image --format "$FORMAT"
