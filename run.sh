#!/bin/sh
docker buildx bake test || exit 1
grep -v '^#' qemu.options | xargs qemu-system-x86_64
