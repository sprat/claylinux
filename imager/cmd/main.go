package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"os"

	"github.com/sprat/claylinux/imager"
	"github.com/sprat/claylinux/imager/efi"
)

func run() error {
	/*
	output = "/out/claylinux"
	format = "raw"
	volume = "CLAYLINUX"
	compression = "gz"

	flag.StringVar(&format, "format", format, "Output format (efi, iso, raw, qcow2, vmdk, vhdx, vdi)")
	flag.StringVar(&output, "output", output, "Output image path/name, without extension")
	flag.StringVar(&volume, "volume", volume, "Volume label for the boot partition")
	flag.StringVar(&compression, "compression", compression, "Compression format for initramfs: none|gz|xz|zstd")
	flag.Parse()
	*/

	fmt.Printf("EFI architecture: %s\n", efi.Arch)

	tmpDir := "/tmp"
	outDir := "/out"
	rootfsDir := "/system"

	// check if the rootfsDir is populated
	if _, err := os.Stat(rootfsDir); err != nil {
		return errors.New(fmt.Sprintf("the %s directory does not exist, please copy or mount your root filesystem here", rootfsDir))
	}

	// make sure these directories exist
	os.Mkdir(tmpDir, 0700)
	os.Mkdir(outDir, 0700)

	// create a temporary build directory
	buildDir, err := os.MkdirTemp(tmpDir, "claylinux-")
	if err != nil {
		return err
	}
	fmt.Printf("Temporary build directory: %s\n", buildDir)
	defer os.RemoveAll(buildDir)

	// build the Unified Kernel Image
	efiFile, err := imager.BuildUKI(rootfsDir, buildDir)
	if err != nil {
		return err
	}

	err = os.Rename(efiFile, filepath.Join(outDir, "claylinux.efi"))
	if err != nil {
		return err
	}

	/*
	espFile = filepath.Join(buildDir, "claylinux.esp")

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

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
