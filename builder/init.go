package main

import (
	"log"
	"os"

	"github.com/otiai10/copy"
	"golang.org/x/sys/unix"
)

// See also:
// - https://stackoverflow.com/questions/51779243/copy-a-folder-in-go
// - https://github.com/moby/moby/blob/master/daemon/graphdriver/copy/copy.go

func relocateRootFS() error {
	const newRoot = "/newroot"

	// create a directory to host the new root mountpoint
	if err := os.Mkdir(newRoot, 0777); err != nil {
		return err
	}

	// mount a tmpfs as new root
	if err := unix.Mount("rootfs", newRoot, "tmpfs", 0, ""); err != nil {
		return err
	}

	// copy all the files to the new root
	options := copy.Options{
		Skip: func(srcinfo os.FileInfo, src, dest string) (bool, error) { // skip the new root directory
			return src == newRoot, nil
		},
	}
	// TODO: this does not copy the special files (device files for sure, and fifo probably)
	if err := copy.Copy("/", newRoot, options); err != nil {
		return err
	}

	// delete everything in / except the new root
	files, err := os.ReadDir("/")
	if err != nil {
		return err
	}
	for _, file := range files {
		fullName := "/" + file.Name()
		if fullName != newRoot {
			if err := os.RemoveAll(fullName); err != nil {
				return err
			}
		}
	}

	// chdir to the new root directory
	if err := os.Chdir(newRoot); err != nil {
		return err
	}

	// mount --move cwd (i.e. /newroot) to /
	if err := unix.Mount(".", "/", "", unix.MS_MOVE, ""); err != nil {
		return err
	}

	// chroot to .
	if err := unix.Chroot("."); err != nil {
		return err
	}

	// chdir to "/" to fix up . and ..
	return os.Chdir("/")
}

func main() {
	const ramfsMagic = 0x858458f6
	const tmpfsMagic = 0x01021994
	var sfs unix.Statfs_t

	if err := unix.Statfs("/", &sfs); err != nil {
		log.Fatalf("Cannot statfs /: %v", err)
	}

	// Some programs (e.g. runc) refuse to work if the rootfs is a tmpfs or ramfs.
	// So, we need to copy all the files into a new tmpfs and make it the new rootfs.
	if sfs.Type == ramfsMagic || sfs.Type == tmpfsMagic {
		if err := relocateRootFS(); err != nil {
			log.Fatalf("Cannot relocate rootfs: %v", err)
		}
	}

	// Run /sbin/init
	if err := unix.Exec("/sbin/init", []string{"init"}, os.Environ()); err != nil {
		log.Fatalf("Cannot exec /sbin/init")
	}
}
