# syntax=docker/dockerfile:1.8.1
FROM koalaman/shellcheck-alpine:v0.10.0 AS shellcheck
WORKDIR /src
RUN --mount=type=bind,target=. \
shellcheck --version && find . -name '*.sh' -exec shellcheck {} +

# =========================================================
FROM hadolint/hadolint:2.12.0-alpine AS hadolint
WORKDIR /src
RUN --mount=type=bind,target=. \
hadolint --version && find . -name Dockerfile -exec hadolint {} +

# =========================================================
FROM toolhippie/yamllint:1.35.1 AS yamllint
WORKDIR /src
RUN --mount=type=bind,target=. \
yamllint -v && yamllint -s -f colored .
