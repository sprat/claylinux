# syntax=docker/dockerfile:1
FROM koalaman/shellcheck-alpine:stable AS shellcheck
WORKDIR /src
RUN --mount=type=bind,target=. \
shellcheck --version && find . -name '*.sh' -exec shellcheck {} +

# =========================================================
FROM hadolint/hadolint:latest-alpine AS hadolint
WORKDIR /src
RUN --mount=type=bind,target=. \
hadolint --version && find . -name Dockerfile -exec hadolint {} +

# =========================================================
FROM cytopia/yamllint:alpine AS yamllint
RUN --mount=type=bind,target=. \
yamllint -v && yamllint -s -f colored .
