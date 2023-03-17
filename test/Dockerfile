# syntax=docker/dockerfile:1
FROM base AS system

# basic setup:
# - add the ucode packages
# - enable passwordless root login
RUN set -euxo pipefail; \
apk add --no-cache intel-ucode amd-ucode; \
passwd -d root

# =========================================================
FROM builder AS build
ARG FORMAT
RUN --mount=from=system,target=/system build-image

# =========================================================
FROM scratch AS output
COPY --from=build /out /
