package imager

import (
	"bytes"

	//"github.com/cavaliergopher/cpio"
	//"github.com/u-root/u-root/pkg/cpio"
	"kraftkit.sh/cpio"
)

func BuildInitramfs(rootfsDir string, buildDir string) error {
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new cpio archive.
	w := cpio.NewWriter(buf)

	// Add some files to the archive.
	var files = []struct {
		Name, Body string
	}{
		{"readme.txt", "This archive contains some text files."},
		{"gopher.txt", "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
		{"todo.txt", "Get animal handling license."},
	}
	for _, file := range files {
		hdr := &cpio.Header{
			Name: file.Name,
			Mode: 0600,
			Size: int64(len(file.Body)),
		}
		if err := w.WriteHeader(hdr); err != nil {
			return err
		}
		if _, err := w.Write([]byte(file.Body)); err != nil {
			return err
		}
	}

	// Make sure to check the error on Close.
	if err := w.Close(); err != nil {
		return err
	}

	return nil

	/*
	os.Mkdir("initramfs_files", 0755)
	run("cp", "/usr/share/claylinux/init", "initramfs_files")

	if fileExists("/system/etc/hosts.target") {
		os.MkdirAll("initramfs_files/etc", 0755)
		run("cp", "/system/etc/hosts.target", "initramfs_files/etc/hosts")
	}
	if fileExists("/system/etc/resolv.conf.target") {
		os.MkdirAll("initramfs_files/etc", 0755)
		run("cp", "/system/etc/resolv.conf.target", "initramfs_files/etc/resolv.conf")
	}
	// cpio for initramfs_files
	run("sh", "-c", `find initramfs_files -mindepth 1 -printf '%P\0' | cpio --quiet -o0H newc -D initramfs_files -F initramfs.img`)
	// Add system files except boot, hosts.target, resolv.conf.target
	run("sh", "-c", `find /system -path /system/boot -prune -o ! -path /system/init ! -path /system/etc/hosts.target ! -path /system/etc/resolv.conf.target -mindepth 1 -printf '%P\0' | cpio --quiet -o0AH newc -D /system -F initramfs.img`)
	compress("initramfs.img")

	ucode := runOutput("find", "/system/boot/", "-name", "*-ucode.img")
	imgs := "initramfs.img"
	if ucode != "" {
		imgs = ucode + " " + imgs
	}
	run("sh", "-c", fmt.Sprintf("cat %s >initramfs", imgs))
	run("find", ".", "!", "-name", "initramfs", "-delete")
	*/
}

/*
func compress(filename string) {
	"none", "gz", "xz", "zstd"
	"invalid compression scheme: " + compression
}
*/
