variable "NAMESPACE" {
  default = "claylinux/"
}

variable "TAG" {
  default = "latest"
}

group "default" {
    targets = ["builder", "alpine", "alpine-virt"]
}

target "builder" {
    context = "builder"
    tags = ["${NAMESPACE}builder:${TAG}"]
}

target "alpine" {
    context = "alpine"
    tags = ["${NAMESPACE}alpine:${TAG}"]
}

target "alpine-virt" {
    context = "alpine"
    tags = ["${NAMESPACE}alpine:${TAG}-virt"]
    args = {
        FLAVOR = "virt"
    }
}

target "test" {
    context = "test"
    contexts = {
        base = "target:alpine-virt"
        builder = "target:builder"
    }
    output = ["type=local,dest=out"]
}
