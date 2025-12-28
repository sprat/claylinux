package imager

import (
	"fmt"
	"path/filepath"
	"os"

	"github.com/sprat/claylinux/imager/efi"
)

func (i Image) Build() error {
	var err error

	// show the detected EFI architecture
	fmt.Printf("EFI architecture: %s\n", efi.Arch)

	// ensure that the specified RootFsDir exists
	if _, err := os.Stat(i.RootFsDir); err != nil {
		return fmt.Errorf("the rootfs directory %s does not exist", i.RootFsDir)
	}

	// make sure the output directory exists
	if err := os.MkdirAll(filepath.Dir(i.Output), 0750); err != nil {
		return err
	}

	// create a temporary build directory
	i.BuildDir, err = os.MkdirTemp("", "claylinux-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(i.BuildDir)

	// build the Unified Kernel Image
	fmt.Println("Building the Unified Kernel Image...")
	efiFile, err := i.buildUKI()
	if err != nil {
		return err
	}

	fmt.Println("Writing the output")
	err = os.Rename(efiFile, i.Output + ".efi")
	if err != nil {
		return err
	}

	/*
	espFile = spec.output + ".esp"

	os.MkdirAll(filepath.Dir(output), 0755)
	switch format {
	case "efi":
		generateEFI()
	case "iso":
		generateISO()
	case "raw":
		generateRaw()
	case "qcow2", "vmdk", "vhdx", "vdi":
		convertImage(format)
	default:
		die("invalid format: " + format)
	}
	*/

	return nil
}
