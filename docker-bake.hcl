variable "REPOSITORY" {
  default = "claylinux"
}

variable "TAG" {
  default = "dev"
}

variable "PLATFORMS" {
  # comma-separated list of platform, e.g. "linux/amd64,linux/arm64"
  # leave empty to use the default build platform
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

group "test-os-all" {
  targets = ["test-os-efi", "test-os-raw", "test-os-qcow2", "test-os-iso"]
}

group "all" {
  targets = ["default", "test-os-all", "efi-firmware"]
}

target "builder" {
  inherits = ["_oci-image"]
  context = "builder"
  tags = tag("builder")
}

# TODO: we should factor the alpine images
target "alpine-lts" {
  inherits = ["_oci-image"]
  context = "alpine"
  tags = tag("alpine-lts")
  args = {
    FLAVOR = "lts"
  }
}

target "alpine-virt" {
  inherits = ["_oci-image"]
  context = "alpine"
  tags = tag("alpine-virt")
  args = {
    FLAVOR = "virt"
  }
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
    "claylinux/alpine-virt" = "target:alpine-virt"
    "claylinux/builder" = "target:builder"
  }
  output = ["type=local,dest=out"]
}

target "_oci-image" {
  pull = true
  platforms = split(",", "${PLATFORMS}")
}

function "tag" {
  params = [name]
  result = ["${REPOSITORY}/${name}:${TAG}"]
}
