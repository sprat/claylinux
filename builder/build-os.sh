#!/bin/sh
# shellcheck shell=dash
set -euo pipefail

BOOT_DIR=/boot
SYSTEM_DIR=/system
OUTPUT_DIR=/out
OUTPUT_NAME=claylinux
VOLUME_NAME=$(echo $OUTPUT_NAME | tr '[:lower:]' '[:upper:]')
SYSLINUX_DIR=/usr/share/syslinux

# Prepare the OS files
prepare_os() {
	mkdir -p $BOOT_DIR $OUTPUT_DIR

	echo "Copying the boot directory"
	cp -R $SYSTEM_DIR/boot/* $BOOT_DIR

	# TODO: generate a squashfs (probably)
	# echo "Generating the initramfs"
	# here=$(pwd)
	# cd $SYSTEM_DIR
	# find . -path './boot' -prune -o -print | cpio -oH newc | gzip -9 >$BOOT_DIR/initrd
	# cd "$here"
}

# Just copy the files to the output directory
create_raw() {
	echo "Copying the OS files to the output directory"

	mv $BOOT_DIR/* $OUTPUT_DIR
}

# Creates a disk image with a MBR and a single whole disk partition
create_img() {
	disk_file=$OUTPUT_DIR/$OUTPUT_NAME.img
	disk_label=dos  # for sfdisk
	syslinux_mbr=$SYSLINUX_DIR/mbr.bin
	syslinux_mbr_length=440
	part_file=/tmp/part.img
	boot_dir_size=$(du -sk $BOOT_DIR | cut -f1)  # in KB
	# add some headroom to account for the filesystem overhead & the missing syslinux bootloader files
	part_size=$(( boot_dir_size + 4 * 1024 ))

	echo "Generating the disk image"

	# create a FAT32 filesystem
	mkdir -p "$(dirname $part_file)"
	mkfs.vfat -n "$VOLUME_NAME" -F 32 -C $part_file $part_size
	# copy the system boot files (i.e. kernel, initramfs, etc.)
	mcopy -i $part_file $BOOT_DIR/* ::/
	# copy the syslinux bootloader
	syslinux --install $part_file

	# compute the disk size
	part_size=$(stat -c %s "$part_file")
	disk_size=$(( part_size + 2 * 1024 * 1024 ))

	# create a blank image file
	dd if=/dev/zero of=$disk_file bs=$disk_size count=1 conv=notrunc

	# add a single full disk partition formatted as FAT32 -- 0c = FAT32 with LBA
	sfdisk $disk_file <<-EOF
	label: $disk_label
	2MiB - 0c *
	EOF

	# overwrite the MBR
	dd if=$syslinux_mbr of=$disk_file bs=$syslinux_mbr_length count=1 conv=notrunc

	# copy in our partition data into the first partition of the disk image
	dd if=$part_file of=$disk_file bs=2M seek=1 conv=notrunc

	rm -rf $part_file
}

# Create an hybrid ISO image
create_iso() {
	iso_file=$OUTPUT_DIR/$OUTPUT_NAME.iso

	echo "Generating the iso image"

	mkdir -p $BOOT_DIR/isolinux
	cp $SYSLINUX_DIR/isolinux.bin $BOOT_DIR/isolinux
	cp $SYSLINUX_DIR/ldlinux.c32 $BOOT_DIR/isolinux

	# see https://wiki.syslinux.org/wiki/index.php?title=Isohybrid
	xorriso -as mkisofs \
	-output $iso_file \
	-isohybrid-mbr $SYSLINUX_DIR/isohdpfx.bin \
	-eltorito-catalog isolinux/boot.cat \
	-eltorito-boot isolinux/isolinux.bin \
	-no-emul-boot \
	-boot-load-size 4 \
	-boot-info-table \
	-joliet \
	-full-iso9660-filenames \
	-volid "$VOLUME_NAME" \
	$BOOT_DIR
	# -input-charset utf8 \
	# -joliet-long \
}

prepare_os
create_iso
