# claylinux

Claylinux is a toolbox to build live OS images using BuildKit & some Dockerfiles.

To build the images, use:
```bash
docker buildx bake
```

## Development notes

Alpine Linux has multiple boot modes: diskless, data and sys.

During the boot, the init script found in the initramfs tries to find the boot media using the `nlplug-findfs` (found
in the `mkinitfs` package).

The official Alpine ISO image does not use a `root=`  kernel parameter: it seems that the CDROM boot media is found
because it contains a `.boot_repository` file, and the boot process continues thanks to that.

Some interesting links:
- https://wiki.alpinelinux.org/wiki/Create_a_Bootable_Device
- https://wiki.alpinelinux.org/wiki/PXE_boot
- https://wiki.alpinelinux.org/wiki/Alpine_Package_Keeper#Upgrading_.22diskless.22_and_.22data.22_disk_mode_installs
