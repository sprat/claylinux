variable "REPOSITORY" {
  default = "claylinux"
}

variable "TAG" {
  default = ""
}

variable "PLATFORMS" {
  # comma-separated list of platform, e.g. "linux/amd64,linux/arm64", leave empty to use the default platform
  default = ""
}

group "lint" {
  targets = ["yamllint", "hadolint", "shellcheck"]
}

target "shellcheck" {
  inherits = ["_lint"]
  target = "shellcheck"
}

target "hadolint" {
  inherits = ["_lint"]
  target = "hadolint"
}

target "yamllint" {
  inherits = ["_lint"]
  target = "yamllint"
}

target "_lint" {
  dockerfile = "lint.Dockerfile"
  output = ["type=cacheonly"]
}

group "default" {
  targets = ["builder", "alpine-lts", "alpine-virt"]
}

group "test" {
  targets = ["test-os-efi", "test-os-raw", "test-os-qcow2", "test-os-iso"]
}

target "builder" {
  inherits = ["_image"]
  context = "builder"
  tags = tags("builder")
}

# TODO: we should factor the alpine images
target "alpine-lts" {
  inherits = ["_image"]
  context = "alpine"
  tags = tags("alpine-lts")
  args = {
    FLAVOR = "lts"
  }
}

target "alpine-virt" {
  inherits = ["_image"]
  context = "alpine"
  tags = tags("alpine-virt")
  args = {
    FLAVOR = "virt"
  }
}

target "_image" {
  pull = true
  platforms = split(",", "${PLATFORMS}")
}

target "efi-firmware" {
  context = "efi-firmware"
  output = ["type=local,dest=out"]
}

target "test-os-efi" {
  inherits = ["_test-os"]
  args = {
    FORMAT = "efi"
  }
}

target "test-os-raw" {
  inherits = ["_test-os"]
  args = {
    FORMAT = "raw"
  }
}

target "test-os-qcow2" {
  inherits = ["_test-os"]
  args = {
    FORMAT = "qcow2"
  }
}

target "test-os-iso" {
  inherits = ["_test-os"]
  args = {
    FORMAT = "iso"
  }
}

target "_test-os" {
  context = "test-os"
  contexts = {
    "claylinux/alpine-virt:latest" = "target:alpine-virt"
    "claylinux/builder:latest" = "target:builder"
  }
  output = ["type=local,dest=out"]
}

function "tags" {
  params = [name]
  result = [
    "${REPOSITORY}/${name}:latest",
    notequal(TAG, "") ? "${REPOSITORY}/${name}:${TAG}" : ""
  ]
}
