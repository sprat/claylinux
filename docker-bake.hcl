variable "NAMESPACE" {
  default = "claylinux/"
}

variable "TAG" {
  default = "latest"
}

group "default" {
    targets = ["builder", "alpine-lts", "alpine-virt"]
}

target "builder" {
    context = "builder"
    tags = ["${NAMESPACE}builder:${TAG}"]
}

# TODO: we should factor the alpine images
target "alpine-lts" {
    context = "alpine"
    tags = ["${NAMESPACE}alpine:${TAG}-lts"]
    args = {
        FLAVOR = "lts"
    }
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
