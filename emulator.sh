#!/bin/bash
set -euo pipefail

: "${FORMAT:=efi}"

name=images/claylinux
case "$FORMAT" in
      efi)
            boot_opts=(-kernel "${name}.efi")
            ;;
      raw)
            boot_opts=(-drive "file=${name}.img,format=raw")
            ;;
      qcow2)
            boot_opts=(-drive "file=${name}.qcow2,format=qcow2")
            ;;
      vmdk)
            boot_opts=(-drive "file=${name}.vmdk,format=vmdk")
            ;;
      vhdx)
            boot_opts=(-drive "file=${name}.vhdx,format=vhdx")
            ;;
      vdi)
            boot_opts=(-drive "file=${name}.vdi,format=vdi")
            ;;
      iso)
            boot_opts=(-drive "file=${name}.iso,media=cdrom")
            ;;
      *)
            echo "invalid format '$FORMAT'" >&2
            exit 1
            ;;
esac

# Use the proper accelerator
if [ -e /dev/kvm ]; then
   boot_opts+=(-accel kvm -cpu host)
fi

# machine: q35 correspond to a modern PC
exec qemu-system-x86_64 \
-machine q35 \
-bios /usr/share/OVMF/OVMF.fd \
-m 2G \
-nodefaults \
-nographic \
-serial mon:stdio \
-netdev user,id=net \
-device virtio-net,netdev=net \
-device virtio-rng \
"${boot_opts[@]}" \
"$@"
# see https://gitlab.com/qemu-project/qemu/-/issues/513
# smm=on
# -global driver=cfi.pflash01,property=secure,value=on \
# -drive if=pflash,format=raw,unit=0,file=out/OVMF_CODE.fd,readonly=on \
# -drive if=pflash,format=raw,unit=1,file=out/OVMF_VARS.fd \
