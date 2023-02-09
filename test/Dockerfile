# syntax=docker/dockerfile:1.5
FROM builder AS build
RUN --mount=type=bind,from=system,target=/system build-os

# =========================================================
FROM scratch AS output
COPY --from=build /out /
