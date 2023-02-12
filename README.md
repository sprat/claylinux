# claylinux

Claylinux is a toolbox to build live OS images using BuildKit & some Dockerfiles.

To build the images, use:
```bash
docker buildx bake
```

## Development notes

The boot process of a Linux system is the following:
1. The BIOS starts and tries to find a bootable device, respecting the boot order configuration in the BIOS. It looks
in the MBR part of each device (first 512 bytes) in order to find a bootloader (i.e. a particular signature in the MBR).
2. The first-stage bootloader (from e.g. syslinux, grub) tries to find a partition with the `boot` flag. If found, it
executes the second-stage bootloader.
3. The second-stage bootloader read its configuration (e.g. `syslinux.cfg`) in order to determine the kernel,
initrd/initramfs and kernel command-line to use.
4. The kernel is started, then it runs the `/init` script found in the initramfs
5. Then the Alpine's init script tries to find & prepare the final system and `switch_root` on it.

The Alpine's init script supports multiple boot modes: `diskless`, `data` and `sys`:
- `sys` correspond to a classical installation on a hard disk, the system configuration persists after a reboot
- `diskless` is used to run a system fully in RAM from a CD-ROM or USB key. Every change is lost after a reboot. In
this mode, you can load a custom configuration of the system from a `apkovl` file.
- `data` is similar to `diskless`, but the `/var` directory is persisted on a separate disk.

During the boot, the Alpine's init script tries to find the boot media using the `nlplug-findfs` utility
(found in the `mkinitfs` package). This utility use the `root` parameter (if provided) to determine the device
containing the final rootfs. If no `root` parameter is specified (i.e. diskless, the `nlplug-findfs` tries to find a

The official Alpine ISO image does not use a `root=`  kernel parameter: it seems that the CDROM boot media is found
because it contains a `.boot_repository` file, and the boot process continues thanks to that.

Some interesting links:
- https://wiki.alpinelinux.org/wiki/Create_a_Bootable_Device
- https://wiki.alpinelinux.org/wiki/PXE_boot
- https://wiki.alpinelinux.org/wiki/Alpine_Package_Keeper#Upgrading_.22diskless.22_and_.22data.22_disk_mode_installs
- https://wiki.alpinelinux.org/wiki/Alpine_Source_Map_by_boot_sequence
