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
  # default format for the VM output
  default = "efi"
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
  dockerfile = "lint/Dockerfile"
  output = ["type=cacheonly"]
}

group "default" {
  targets = ["builder", "alpine-lts", "alpine-virt"]
}

group "test" {
  targets = ["test-efi", "test-iso", "test-raw", "test-qcow2", "test-vmdk", "test-vhdx", "test-vdi"]
}

group "all" {
  targets = ["lint", "default", "test", "vm"]
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
    CMDLINE = "console=tty0 console=ttyS0"
    INITTAB_TTYS = <<-EOF
      tty1::respawn:/sbin/getty 38400 tty1
      ttyS0::respawn:/sbin/getty -L 0 ttyS0 vt100
    EOF
  }
}

target "_image" {
  target = "image"
  pull = true
  platforms = split(",", "${PLATFORMS}")
}

target "test-efi" {
  inherits = ["_test"]
  args = {
    FORMAT = "efi"
  }
}

target "test-iso" {
  inherits = ["_test"]
  args = {
    FORMAT = "iso"
  }
}

target "test-raw" {
  inherits = ["_test"]
  args = {
    FORMAT = "raw"
  }
}

target "test-qcow2" {
  inherits = ["_test"]
  args = {
    FORMAT = "qcow2"
  }
}

target "test-vmdk" {
  inherits = ["_test"]
  args = {
    FORMAT = "vmdk"
  }
}

target "test-vhdx" {
  inherits = ["_test"]
  args = {
    FORMAT = "vhdx"
  }
}

target "test-vdi" {
  inherits = ["_test"]
  args = {
    FORMAT = "vdi"
  }
}

target "vm" {
  inherits = ["_test"]
  target = "vm"
  output = ["type=image"]
  args = {
    FORMAT = "${FORMAT}"
  }
  tags = tags("vm")
}

target "_test" {
  context = "test"
  contexts = {
    "alpine-virt" = "target:alpine-virt"
    "builder" = "target:builder"
  }
  output = ["type=cacheonly"]
}

function "tags" {
  params = [name]
  result = [
    "${REPOSITORY}/${name}:latest",
    notequal(TAG, "") ? "${REPOSITORY}/${name}:${TAG}" : ""
  ]
}
