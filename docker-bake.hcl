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

variable "FORMAT" {
  # default output format for the test image
  default = "efi"
}

group "default" {
  targets = ["builder", "alpine-lts", "alpine-virt"]
}

group "test" {
  targets = ["test-img", "efi-firmware"]
}

target "builder" {
  inherits = ["_image"]
  context = "builder"
  tags = tag("builder")
}

# TODO: we should factor the alpine images
target "alpine-lts" {
  inherits = ["_image"]
  context = "alpine"
  tags = tag("alpine-lts")
  args = {
    FLAVOR = "lts"
  }
}

target "alpine-virt" {
  inherits = ["_image"]
  context = "alpine"
  tags = tag("alpine-virt")
  args = {
    FLAVOR = "virt"
  }
}

target "test-img" {
  context = "test-img"
  contexts = {
    "claylinux/alpine-virt" = "target:alpine-virt"
    "claylinux/builder" = "target:builder"
  }
  args = {
    FORMAT = "${FORMAT}"
  }
  output = ["type=local,dest=out"]
}

target "efi-firmware" {
  context = "efi-firmware"
  output = ["type=local,dest=out"]
}

target "_image" {
  pull = true
  platforms = split(",", "${PLATFORMS}")
}

function "tag" {
  params = [name]
  result = ["${REPOSITORY}/${name}:${TAG}"]
}
