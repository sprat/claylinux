#!/busybox.static sh
# shellcheck shell=dash
set -eu

# redirect all outputs to the console
exec >/dev/console 2>/dev/console

# Some programs (e.g. containerd) refuse to work if / is a "rootfs" mount, as in an initramfs.
# So, we need to copy all the files into a new tmpfs and "switch_root" to it
b=/busybox.static
$b mkdir -p /newroot
$b mount -t tmpfs root /newroot
$b find / -path /newroot -prune -o ! -path $b ! -path /init -mindepth 1 -maxdepth 1 -exec $b mv -t /newroot {} +
exec $b switch_root /newroot /sbin/init "$@"
