#!/bin/bash
set -euo pipefail

# exit with an error message
die() {
	echo "Error: $*" >&2
	exit 1
}

# get the size of the file in bytes
get_size() {
	stat -c %s "$1"
}

# convert a number of bytes into MiB (i.e. 1024 * 1024 bytes), rounded to the next value
in_mib() {
	echo $(( ($1 + (1<<20) - 1) >> 20 ))
}

# get the size of the file in mibytes
get_size_mib() {
	in_mib "$(get_size "$1")"
}

# round the value given in $1 to the next multiple of the value given in $2
# e.g. align 9 4 -> 12, align 8 4 -> 8
align() {
	echo "$(( ($1 + $2 - 1) / $2 * $2 ))"
}

# build the EFI executable
build_efi() {
	local size kernel

	pushd "$build_dir" >/dev/null

	echo "Building the EFI executable"
	build_initramfs
	size=$(get_size_mib initramfs)
	echo "The size of the initramfs is: $size MiB"
	kernel=$(find /system/boot -name 'vmlinu*' -print)
	space_separated </system/boot/cmdline >cmdline
	basename /system/lib/modules/* >kernel-release

	# build the EFI UKI file
	# TODO: add .dtb section on ARM?
	build_uki <<-EOF
	.osrel /system/etc/os-release
	.uname kernel-release
	.cmdline cmdline
	.initrd initramfs
	.linux $kernel
	EOF

	# delete all the temporary files
	find . ! -name '*.efi' -delete

	popd >/dev/null
}

# build the initramfs
build_initramfs() {
	mkdir initramfs_files

	# copy the init script
	cp /usr/share/claylinux/init initramfs_files

	# copy /etc/hosts.target as /etc/hosts
	if [[ -f /system/etc/hosts.target ]]; then
		mkdir -p initramfs_files/etc
		cp /system/etc/hosts.target initramfs_files/etc/hosts
	fi

	# copy etc/resolv.conf.target as etc/resolv.conf
	if [[ -f /system/etc/resolv.conf.target ]]; then
		mkdir -p initramfs_files/etc
		cp /system/etc/resolv.conf.target initramfs_files/etc/resolv.conf
	fi

	# create an initramfs with these files
	find initramfs_files -mindepth 1 -printf '%P\0' \
	| cpio --quiet -o0H newc -D initramfs_files -F initramfs.img

	# append the system files into the initramfs image, except /boot, /etc/hosts.target and /etc/resolv.conf.target
	find /system \
	-path /system/boot -prune -o \
	! -path /system/init \
	! -path /system/etc/hosts.target \
	! -path /system/etc/resolv.conf.target \
	-mindepth 1 -printf '%P\0' \
	| cpio --quiet -o0AH newc -D /system -F initramfs.img

	# compress the initramfs
	compress initramfs.img

	# build the final initramfs by concatenating the ucode images & our compressed initramfs image
	# see https://docs.kernel.org/arch/x86/microcode.html
	echo "$(find /system/boot/ -name '*-ucode.img') initramfs.img" | xargs cat >initramfs

	# remove the temporary files
	find . ! -name initramfs -delete
}

# create a Unified Kernel Image from the sections passed in the standard input
build_uki() {
	# the sections addresses should be aligned to PAGE_ALIGN(), i.e. 2<<12 == 4096 bytes
	local args=() alignment=4096 size offset

	# compute the start offset of the new sections
	offset="$(objdump -h -w "$efi_stub" | awk 'END { offset=("0x"$4)+0; size=("0x"$3)+0; print offset + size }')"
	offset=$(align "$offset" $alignment)

	# compute the objcopy arguments
	while read -r section file
	do
		# add the section to the parameters
		args+=(
			--add-section
			"$section=$file"
			--change-section-vma
			"$section=$offset"
		)

		# compute the offset for the next section
		size="$(get_size "$file")"
		size=$(align "$size" $alignment)
		offset=$(( offset + size ))
	done

	objcopy "${args[@]}" "$efi_stub" "$efi_file"
}

# detect the current EFI architecture
get_efi_arch() {
	local machine_arch
	machine_arch="$(uname -m)"
	case "$machine_arch" in
		aarch64)
			echo "aa64"
			;;
		arm*)
			echo "arm"
			;;
		i686)
			echo "ia32"
			;;
		x86_64)
			echo "x64"
			;;
		*)
			die "unsupported architecture: $machine_arch"
			;;
	esac
}

# convert a multi-line input into a space separated list
space_separated() {
	paste -d' ' -s
}

# compress the initramfs with the specified scheme
compress() {
	case "$compression" in
		none)
			;;
		gz)
			pigz -9 "$1"
			mv "$1".gz "$1"
			;;
		xz)
			xz -C crc32 -9 -T0 "$1"
			mv "$1".xz "$1"
			;;
		zstd)
			zstd -19 -T0 --rm "$1"
			mv "$1".zstd "$1"
			;;
		*)
			die "invalid compression scheme '$compression'"
			;;
	esac
}

# just copy the build files to the output directory
generate_efi() {
	echo "Copying the OS files to the output directory"
	mv "$efi_file" "$output".efi
}

# generate the EFI system partition
generate_esp() {
	local size

	echo "Generating the EFI System Partition (ESP)"

	# compute the ESP size:
	# - measure the apparent size of the files to copy, in bytes
	# - add some headroom to account for the filesystem overhead (less than 2%)
	# - convert the result into MiB
	size=$(get_size "$efi_file")
	size=$(in_mib $(( size * 102 / 100 )))
	echo "The size of the ESP is: $size MiB"

	# create the FAT32 filesystem
	mkfs.vfat -n "$volume" -F 32 -C "$esp_file" "$(( size << 10 ))" >/dev/null

	# copy the EFI executable into the filesystem
	mmd -i "$esp_file" ::/EFI
	mmd -i "$esp_file" ::/EFI/boot
	mcopy -i "$esp_file" "$efi_file" "::/EFI/boot/boot${efi_arch}.efi"

	rm "$efi_file"
}

# generate an ISO image
generate_iso() {
	generate_esp

	echo "Generating the ISO image"

	xorrisofs \
	-e "$(basename "$esp_file")" \
	-no-emul-boot \
	-joliet \
	-full-iso9660-filenames \
	-rational-rock \
	-sysid LINUX \
	-volid "$volume" \
	-output "$output".iso \
	"$esp_file" 2>/dev/null

	rm "$esp_file"
}

# generate a raw disk image with a single whole disk FAT32 EFI partition on GPT
generate_raw() {
	local disk_file esp_size disk_size

	generate_esp

	echo "Generating the disk image"
	disk_file="$output".img

	# compute the disk size in bytes: get the ESP size and add 1MB at both ends for the partition tables
	esp_size=$(get_size_mib "$esp_file")
	disk_size=$(( esp_size + 2 ))
	echo "The size of the disk image is: $disk_size MiB"

	# create a blank image file
	truncate -s ${disk_size}M "$disk_file"

	# add a single full disk partition formatted as FAT32 with LBA
	sfdisk --quiet "$disk_file" <<-EOF
	label: gpt
	first-lba: 34
	start=1MiB size=${esp_size}MiB name="EFI system partition" type=uefi
	EOF

 	# copy in our partition data into the first partition of the disk image
	dd if="$esp_file" of="$disk_file" bs=1M seek=1 conv=notrunc status=none

	rm "$esp_file"
}

# convert the raw disk image to another format
convert_image() {
	local format="$1"

	generate_raw

	echo "Converting the disk image to $format format"
	qemu-img convert -f raw -O "$format" "$output".img "$output"."$format"

	rm "$output".img
}

# generate a disk image in qcow2 format
generate_qcow2() {
	convert_image qcow2
}

# generate a disk image in vmdk format
generate_vmdk() {
	convert_image vmdk
}

# generate a disk image in vhdx format
generate_vhdx() {
	convert_image vhdx
}

# generate a disk image in vdi format
generate_vdi() {
	convert_image vdi
}

# validate and set the output format
set_format() {
	case "$1" in
		efi|iso|raw|qcow2|vmdk|vhdx|vdi)
			format="$1"
			;;
		*)
			die "invalid format '$1'"
			;;
	esac
}

# defaults
output=/out/claylinux
format=raw
volume=CLAYLINUX
compression=gz
efi_arch=$(get_efi_arch)
efi_stub="/usr/lib/systemd/boot/efi/linux${efi_arch}.efi.stub"

usage=$(cat <<-EOF
	Usage: $(basename "$0") [OPTIONS ...]

	Build an OS image from the root filesystem found in /system

	Options:
	  -f, --format FORMAT        Output format (default: $format)
	  -o, --output OUTPUT        Output image path/name, without any extension (default: $output)
	  -v, --volume VOLUME        Volume name/label of the boot partition (default: $volume)
	  -c, --compression COMP     Compression format for the initramfs: none | gz | xz | zstd (default: $compression)

	Output formats:
	  - efi: EFI executable (saved as OUTPUT.efi), for use with a custom bootloader or with PXE boot
	  - iso: ISO9660 CD-ROM image (saved as OUTPUT.iso)
	  - raw: raw disk image with a single FAT32 boot partition (saved as OUTPUT.img)
	  - qcow2: disk image in QCOW2 format (saved as OUTPUT.qcow2)
	  - vmdk: disk image in VMDK format (saved as OUTPUT.vmdk)
	  - vhdx: disk image in VHDX format (saved as OUTPUT.vhdx)
	  - vdi: disk image in VDI format (saved as OUTPUT.vdi)
	EOF
)

# parse the command-line arguments
while [[ "$#" -gt 0 ]]; do
	case "$1" in
		-h|--help)
			echo "$usage"
			exit 0
			;;
		-f|--format)
			set_format "$2"
			shift 2
			;;
		-o|--output)
			output="$2"
			shift 2
			;;
		-v|--volume)
			volume="$2"
			shift 2
			;;
		-c|--compression)
			compression="$2"
			shift 2
			;;
		-*)
			die "invalid option '$1'"
			;;
		*)
			die "invalid parameter '$1'"
			;;
	esac
done

[[ -d /system ]] || die "the /system directory does not exist, please copy/mount your root filesystem here"

build_dir=$(mktemp -d)
efi_file="$build_dir"/claylinux.efi
esp_file="$build_dir"/claylinux.esp
build_efi
mkdir -p "$(dirname "$output")"
generate_"$format"
rmdir "$build_dir"
