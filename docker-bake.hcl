variable "REPOSITORY" {
  default = "claylinux"
}

variable "TAG" {
  default = "latest"
}

group "default" {
    targets = ["builder", "alpine-lts", "alpine-virt"]
}

target "builder" {
    context = "builder"
    tags = tag("builder")
}

# TODO: we should factor the alpine images
target "alpine-lts" {
    context = "alpine"
    tags = tag("alpine-lts")
    args = {
        FLAVOR = "lts"
    }
}

target "alpine-virt" {
    context = "alpine"
    tags = tag("alpine-virt")
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

function "tag" {
  params = [name]
  result = ["${REPOSITORY}/${name}:${TAG}"]
}
