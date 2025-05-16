package main

import (
	"log"
	"os"

	"github.com/sprat/claylinux/imager/efi"
)

func main() {
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

	log.Printf("EFI architecture: %s", efi.Arch)

	tmpDir := "/tmp"
	os.Mkdir(tmpDir, 0700)

	buildDir, err := os.MkdirTemp(tmpDir, "claylinux-imager-")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Build directory: %s", buildDir)
	defer os.RemoveAll(buildDir)

	rootfsDir := "/system"
	if _, err := os.Stat(rootfsDir); err != nil {
		log.Fatalf("Error: the %s directory does not exist, please copy/mount your root filesystem here", rootfsDir)
	}

	efiFile, err := efi.BuildUKI(rootfsDir, buildDir)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Rename(efiFile, "/out/claylinux.efi")
	if err != nil {
		log.Fatal(err)
	}

	/*
	espFile = filepath.Join(buildDir, "claylinux.esp")

	os.Chdir(buildDir)
	buildEFI()
	os.Chdir("/")

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
}
