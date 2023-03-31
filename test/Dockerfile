# syntax=docker/dockerfile:1
# hadolint ignore=DL3006
FROM alpine AS system

# basic setup:
# - add the ucode packages
# - enable passwordless root login
RUN apk add --no-cache intel-ucode amd-ucode && passwd -d root

# =========================================================
# hadolint ignore=DL3006
FROM builder AS build
ARG FORMAT=raw
RUN --mount=from=system,target=/system build-image --format $FORMAT

# =========================================================
FROM scratch AS output
COPY --from=build /out /
