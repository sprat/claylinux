# syntax=docker/dockerfile:1.6.0
# Setup a custom OS
# hadolint ignore=DL3006
FROM base AS system
ARG UCODE=none
RUN if [ "$UCODE" != "none" ]; then apk add --no-cache "${UCODE}-ucode"; fi

# =========================================================
# Generate the OS image
# hadolint ignore=DL3006
FROM builder AS test
ARG FORMAT=efi
RUN --mount=from=system,target=/system build-image --format "$FORMAT"
