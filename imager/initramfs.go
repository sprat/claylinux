package imager

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"kraftkit.sh/cpio"
)

// See https://github.com/unikraft/kraftkit/blob/staging/initrd/directory.go
func addToCpioArchive(writer *cpio.Writer, path string, fileInfo fs.FileInfo, name string) error {
	var err error
	var data []byte

	fmt.Printf("Adding %s\n", name)

	mode := fileInfo.Mode()

	targetLink := ""
	if mode & os.ModeSymlink != 0 {
		// find link target
		targetLink, err = os.Readlink(path)
		data = []byte(targetLink)
	} else if mode.IsRegular() {
		// TODO: don't read the whole file in RAM?
		data, err = os.ReadFile(path)
	}
	if err != nil {
		return err
	}

	header, err := cpio.FileInfoHeader(fileInfo, targetLink)
	if err != nil {
		return err
	}

	header.Name = name

	err = writer.WriteHeader(header)
	if err != nil {
		return err
	}

	_, err = writer.Write(data)
	if err != nil {
		return err
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func buildInitramfs(rootfsDir string, buildDir string) (string, error) {
	initRamFsPath := filepath.Join(buildDir, "initramfs.img")

	// create the initramfs image file
	file, err := os.Create(initRamFsPath)
	if err != nil {
		return "", err
	}

    // remember to close the file
    defer file.Close()

	// create a new cpio archive
	writer := cpio.NewWriter(file)

	// add the rootfs files to the archive
	err = filepath.WalkDir(rootfsDir, func(path string, dirEntry os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == rootfsDir {  // ignore the root entry
			return nil
		}

		fileInfo, err := dirEntry.Info()
		if err != nil {
			return err
		}

		name := "." + filepath.ToSlash(strings.TrimPrefix(path, rootfsDir))

		switch name {  // handle some special cases
		case "./init":
			// ignore the existing file, because we'll add our own ./init binary later
			return nil
		case "./boot":
			// store the directory (so that it can be mounted in the target system), but ignore its contents
			// (kernel, cmdline, ...) since it will be used to build the UKI
			err = addToCpioArchive(writer, path, fileInfo, name)
			if err != nil {
				return err
			}
			return fs.SkipDir
		case "./etc/hosts.target", "./etc/resolv.conf.target":
			// rename the files by removing the ".target" suffix
			// we need alternate names for these files because the container runtime mount these files
			// in order to provide internet access to our build container. We cannot use and store
			// these files in the final image as the target system may not have the same network settings
			// (and it's not possible anyway).
			name = strings.TrimSuffix(name, ".target")
		}

		return addToCpioArchive(writer, path, fileInfo, name)
	})
	if err != nil {
		return "", err
	}

	// add our special ./init binary
	initPath := "/usr/share/claylinux/init"
	initFileInfo, err := os.Stat(initPath)
	if err != nil {
		return "", err
	}

	err = addToCpioArchive(writer, initPath, initFileInfo, "./init")
	if err != nil {
		return "", err
	}

	// check the errors on close
	// TODO: how to use a defer?
	if err := writer.Close(); err != nil {
		return "", err
	}

	return initRamFsPath, nil

	// TODO: compress
}
