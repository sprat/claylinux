package imager

import (
	"fmt"
	"path/filepath"
	"os"

	"github.com/sprat/claylinux/imager/efi"
)

func (i Image) Build() error {
	// show the detected EFI architecture
	fmt.Printf("EFI architecture: %s\n", efi.Arch)

	// ensure that the specified RootFsDir exists
	_, err := os.Stat(i.RootFsDir)
	if err != nil {
		return fmt.Errorf("The rootfs directory %s does not exist", i.RootFsDir)
	}

	// make sure the output directory exists
	os.Mkdir(filepath.Dir(i.Output), 0750)

	// create a temporary build directory
	buildDir, err := os.MkdirTemp("", "claylinux-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(buildDir)

	// build the Unified Kernel Image
	fmt.Println("Building the Unified Kernel Image...")
	efiFile, err := BuildUKI(i.RootFsDir, buildDir)
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
