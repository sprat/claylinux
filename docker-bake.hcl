variable "REPOSITORY" {
  default = "claylinux/"
}

variable "TAG" {
  default = "latest"
}

group "default" {
    targets = ["builder", "alpine"]
}

group "all" {
    targets = ["builder", "alpine", "test"]
}

target "builder" {
    context = "builder"
    tags = ["${REPOSITORY}builder:${TAG}"]
}

target "alpine" {
    context = "alpine"
    tags = ["${REPOSITORY}alpine:${TAG}"]
}

target "test" {
    context = "test"
    contexts = {
        "base" = "target:alpine"
        "builder" = "target:builder"
    }
    output = ["type=local,dest=out"]
}
