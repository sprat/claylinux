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

variable "FORMAT" {
  # default format for the emulator output
  default = "efi"
}

group "lint" {
  targets = ["yamllint", "hadolint", "shellcheck"]
}

group "default" {
  targets = ["lint", "imager", "alpine-lts", "alpine-edge", "alpine-virt"]
}

group "all" {
  targets = ["default", "test", "emulator"]
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
  dockerfile = "lint/Dockerfile"
  output = ["type=cacheonly"]
}

target "imager" {
  inherits = ["_multiplatform_image"]
  context = "imager"
  tags = tags("imager")
}

# TODO: we should factor the alpine images
target "alpine-lts" {
  inherits = ["_multiplatform_image"]
  context = "alpine"
  tags = tags("alpine-lts")
  args = {
    FLAVOR = "lts"
  }
}

target "alpine-edge" {
  inherits = ["_multiplatform_image"]
  context = "alpine"
  tags = tags("alpine-edge")
  args = {
    FLAVOR = "edge"
  }
}

target "alpine-virt" {
  inherits = ["_multiplatform_image"]
  context = "alpine"
  tags = tags("alpine-virt")
  args = {
    FLAVOR = "virt"
    CMDLINE = "console=tty0 console=ttyS0"
    INITTAB_TTYS = <<-EOF
      tty1::respawn:/sbin/getty 38400 tty1
      ttyS0::respawn:/sbin/getty -L 0 ttyS0 vt100
    EOF
  }
}

target "_multiplatform_image" {
  target = "image"
  pull = true
  platforms = split(",", "${PLATFORMS}")
}

target "test" {
  name = "test-${item.target}-${item.format}-${item.ucode}"
  matrix = {
    item = [
      {format = "efi", ucode="intel", target="alpine-lts"},
      {format = "raw", ucode="intel", target="alpine-edge"},
      {format = "iso", ucode="amd", target="alpine-lts"},
      {format = "qcow2", ucode="none", target="alpine-virt"},
      {format = "vmdk", ucode="none", target="alpine-virt"},
      {format = "vhdx", ucode="none", target="alpine-virt"},
      {format = "vdi", ucode="none", target="alpine-virt"}
    ]
  }
  context = "test"
  output = ["type=cacheonly"]
  contexts = {
    "base" = "target:${item.target}"
    "imager" = "target:imager"
  }
  args = {
    FORMAT = "${item.format}"
    UCODE = "${item.ucode}"
  }
}

target "emulator" {
  context = "emulator"
  target = "image"
  output = ["type=image"]
  tags = tags("emulator")
  args = {
    FORMAT = "${FORMAT}"
  }
  contexts = {
    "imager" = "target:imager"
    "alpine-virt" = "target:alpine-virt"
  }
}

function "tags" {
  params = [name]
  result = [
    "${REPOSITORY}/${name}:latest",
    notequal(TAG, "") ? "${REPOSITORY}/${name}:${TAG}" : ""
  ]
}
