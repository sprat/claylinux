# syntax=docker/dockerfile:1

# Setup our custom OS
# hadolint ignore=DL3006
FROM alpine AS system

# =========================================================
# Generate the OS image
# hadolint ignore=DL3006
FROM builder AS build
ARG FORMAT
RUN --mount=from=system,target=/system build-image --format $FORMAT

# =========================================================
# Save the files from the output directory in a new image
FROM scratch AS output
COPY --from=build /usr/share/OVMF/OVMF.fd /
COPY --from=build /out /
