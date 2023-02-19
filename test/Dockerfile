# syntax=docker/dockerfile:1.5
FROM base AS system
# TODO: make sure the ucodes are used in the syslinux config and/or in the initramfs
# RUN apk add --no-cache intel-ucode amd-ucode
# TODO: set the console here

# =========================================================
FROM builder AS build
RUN --mount=type=bind,from=system,target=/system build-os

# =========================================================
FROM scratch AS output
COPY --from=build /out /
