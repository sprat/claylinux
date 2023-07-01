# syntax=docker/dockerfile:1

# Setup our custom OS
# hadolint ignore=DL3006
FROM claylinux/alpine-virt AS system

# =========================================================
# Generate the OS image
# hadolint ignore=DL3006
FROM claylinux/builder AS build
ARG FORMAT
RUN --mount=from=system,target=/system build-image --format "$FORMAT"

# =========================================================
# Save the files from the output directory in a new image that will be exported to disk
FROM scratch AS output
COPY --from=build /out /
