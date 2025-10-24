# claylinux

Claylinux is a tool to build and install [Bootable Containers](https://containers.github.io/bootable/) images. It
allows you to build Linux-based operating systems (OS) using [Docker](https://www.docker.com/) or
[Podman](https://podman.io/) tooling. Then, you can ship your OS images using [OCI](https://opencontainers.org/) image
registries such as [Docker Hub](https://hub.docker.com/), [GitHub Container Registry](https://ghcr.io),
[Quay.io](https://quay.io/)... This tool can also produce raw disk images (`.img`), CDROM/ISO images (`.iso`) and
pure EFI binaries also called
[Unified Kernel Images](https://uapi-group.org/specifications/specs/unified_kernel_image/) (`.efi`)...

The difference between bootable container images and traditional container images is that bootable images also
contain a Linux kernel and an init system in addition to the OS user-space. For that purpose, we provide some
[base images](https://hub.docker.com/u/claylinux) containing all you need to start building your own OS image.

This project is currently a **WORK IN PROGRESS**, some breaking changes may happen at any time.


## Getting Started

The process to follow to build an OS is the following:
1. start from a base docker image which contains both the Linux OS userspace filesystem and a Linux kernel. For example,
you can use the `claylinux/alpine-lts` image which correspond to an [Alpine Linux](https://www.alpinelinux.org/) OS with
a [linux-lts](https://pkgs.alpinelinux.org/package/edge/main/x86_64/linux-lts) kernel.
2. add software & configuration files to the image using a [Dockerfile](https://docs.docker.com/reference/dockerfile/).

Here is an example of Dockerfile which builds a custom Alpine Linux image with `nginx` installed:
```dockerfile
# Initialize the OS
FROM claylinux/alpine-lts:latest AS system
RUN apk add --no-cache nginx && rc-update add nginx default

# =========================================================
# Generate the OS image
FROM claylinux/imager:latest AS build
RUN --mount=from=system,target=/system build-image --format raw

# =========================================================
# Extract the files from the /out directory into a new image
FROM scratch AS output
COPY --from=build /out /
```

Build the OS image with this command:
```bash
docker buildx build --output=out .
```

A `claylinux.img` file is generated into the `out/`  directory, which is a bootable disk image. You can burn this image
on a USB key or a hard drive using for example [dd](https://www.man7.org/linux/man-pages/man1/dd.1.html) on Linux, or
[Rufus](https://rufus.ie/fr/) on Microsoft Windows. Then you can boot a machine on this media and use the Linux OS
you built.


## Tips

Since the OS image is built using docker, the `/etc/hosts` file is already used by docker to provide the networking to
the build container, so it can't be used as the `/etc/hosts` file of the target OS. And it's the same for the
`/etc/resolv.conf` file. So, if you want to change the `/etc/hosts` or `/etc/resolv.conf` file in your target OS,
please edit the `/etc/hosts.target` or `/etc/resolv.conf.target` file, which will be used as `/etc/hosts` or
`/etc/resolv.conf` in the target OS.


## Understanding how the OS works

There are a couple of things to know when using the claylinux tool (and some traps that you can fall into) which are
detailed here:
1. The whole operating system will be loaded into RAM and will run from there. So you should keep your image size
minimal and avoid cluttering your OS filesystem, as it will eat up RAM.
2. Due to point 1, no firmwares are added by default into the base images since they take much space. So, depending on
the target machine you'll boot your OS on, you'll have to add some firmwares manually into your image. For Alpine Linux,
the firmwares are found into the `linux-firmware-*` packages. You can easily identify which firmwares are missing by
looking for the firmware load errors in the `dmesg` output. Note that no additional firmware is needed for virtual
machine OSes using the `virt` kernels.
3. If you need to add big files into your system, you'd better put them on a persistent filesystem that will be
mounted in `/etc/fstab` instead of putting them into the OS image.
4. Also, due to point 1, every change you make on the root filesystem will be lost after a reboot. If you need to
persist some changes, add some mount points to persistent disks in `/etc/fstab` and use them.
5. The kernel will start your operating system using a minimal `init` script. We don't use the Linux distribution's
initramfs mechanisms (such as mkinitfs, mkinitramfs, dracut, ...) as they generally require the root filesystem to
be located on a separated partition/disk. Instead, we bundle everything in a single `.efi` file.
6. But thanks to 4 and 1, there's only one single `.efi` file to sign for Secure Boot, which contains both the kernel,
the initramfs and the OS userland filesystem, and which can't be falsified. So it's pretty secure.


## Development

To build the claylinux images, use:
```bash
docker buildx bake
```

And run some sanity tests with:
```bash
docker buildx bake test
```

Finally, you can build a test OS image & launch it in a qemu emulator using the following command:
```bash
docker buildx bake emulator
docker compose run --rm emulator
```
