# syntax=docker/dockerfile:1

# Setup our custom OS
# hadolint ignore=DL3006
FROM alpine AS system

# =========================================================
# Generate the OS image
# hadolint ignore=DL3006
FROM builder AS image
ARG FORMAT
RUN --mount=from=system,target=/system build-image --format $FORMAT

# =========================================================
# Save our image into the output directory
FROM scratch AS output
COPY --from=image /out /
