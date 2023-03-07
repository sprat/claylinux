# syntax=docker/dockerfile:1
FROM base AS system

# basic setup:
# - add the ucode packages
# - enable passwordless root login
RUN set -euo pipefail; \
apk add --no-cache intel-ucode amd-ucode; \
passwd -d root

# =========================================================
FROM builder AS build
RUN --mount=from=system,target=/system build-img

# =========================================================
FROM scratch AS output
COPY --from=build /out /
