variable "REPOSITORY" {
  default = "claylinux"
}

variable "TAG" {
  default = "latest"
}

group "default" {
    targets = ["builder", "alpine"]
}

target "builder" {
    context = "builder"
    tags = ["${REPOSITORY}/builder:${TAG}"]
}

target "alpine" {
    context = "alpine"
    tags = ["${REPOSITORY}/alpine:${TAG}"]
}
