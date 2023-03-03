variable "NAMESPACE" {
  default = "claylinux/"
}

variable "TAG" {
  default = "latest"
}

group "default" {
    targets = ["builder", "alpine"]
}

target "builder" {
    context = "builder"
    tags = ["${NAMESPACE}builder:${TAG}"]
}

target "alpine" {
    context = "alpine"
    tags = ["${NAMESPACE}alpine:${TAG}"]
}

target "test" {
    context = "test"
    contexts = {
        base = "target:alpine"
        builder = "target:builder"
    }
    output = ["type=local,dest=out"]
}
