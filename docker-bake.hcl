variable "NAMESPACE" {
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
    tags = ["${NAMESPACE}/builder:${TAG}"]
}

target "alpine" {
    context = "alpine"
    tags = ["${NAMESPACE}/alpine:${TAG}"]
}
