# syntax=docker/dockerfile:1

# Setup our custom OS
FROM claylinux/alpine-virt:latest AS system

# =========================================================
# Generate the OS image
FROM claylinux/builder:latest AS build
ARG FORMAT
RUN --mount=from=system,target=/system build-image --format "$FORMAT"

# =========================================================
# Save the files from the output directory in a new image that will be exported to disk
FROM scratch AS output
COPY --from=build /out /
