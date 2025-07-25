#!/bin/bash
set -euo pipefail

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
imager "$build_dir"
mkdir -p "$(dirname "$output")"
generate_"$format"
rmdir "$build_dir"
