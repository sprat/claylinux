# syntax=docker/dockerfile:1
FROM koalaman/shellcheck-alpine:stable AS shellcheck
SHELL ["/bin/ash", "-euxo", "pipefail", "-c"]
WORKDIR /src
RUN --mount=type=bind,target=. \
shellcheck --version && grep -rlE '^#!/bin/.*sh' ./* | xargs shellcheck

# =========================================================
FROM hadolint/hadolint:latest-alpine AS hadolint
WORKDIR /src
RUN --mount=type=bind,target=. \
hadolint --version && find . -name Dockerfile -exec hadolint {} +

# =========================================================
FROM cytopia/yamllint:alpine AS yamllint
RUN --mount=type=bind,target=. \
yamllint -v && yamllint -s -f colored .
