# syntax=docker/dockerfile:1.6.0
ARG FORMAT=efi

# Setup our custom OS
# hadolint ignore=DL3006
FROM alpine-virt AS system
RUN apk add --no-cache intel-ucode

# =========================================================
# Generate the OS image
# hadolint ignore=DL3006
FROM builder AS test
ARG FORMAT
RUN --mount=from=system,target=/system build-image --format "$FORMAT"

# =========================================================
# Generate a qemu image running our custom OS image
FROM alpine:3.19.0 AS vm
RUN apk add --no-cache bash qemu-system-x86_64 ovmf
COPY run-vm.sh /usr/local/bin/run-vm
ENTRYPOINT ["/usr/local/bin/run-vm"]
COPY --from=test /out /out
ARG FORMAT
ENV FORMAT="$FORMAT"
