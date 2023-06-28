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
    tags = ["${NAMESPACE}alpine-lts:${TAG}"]
    args = {
        FLAVOR = "lts"
    }
}

target "alpine-virt" {
    context = "alpine"
    tags = ["${NAMESPACE}alpine-virt:${TAG}"]
    args = {
        FLAVOR = "virt"
    }
}

target "test" {
    context = "test"
    contexts = {
        "claylinux/alpine-virt" = "target:alpine-virt"
        "claylinux/builder" = "target:builder"
    }
    output = ["type=local,dest=out"]
}
