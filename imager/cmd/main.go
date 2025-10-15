package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sprat/claylinux/imager"
)

func main() {
	image := imager.Image{}

	// parse command-line arguments
	flag.StringVar(&image.RootFsDir, "rootfs", "/system", "Root filesystem directory")
	flag.StringVar(&image.Output, "output", "/out/claylinux", "Output image name with path, without any extension")
	flag.StringVar(&image.Format, "format", "efi", "Output format (efi, iso, raw, qcow2, vmdk, vhdx, vdi)")
	// flag.StringVar(&volume, "volume", "CLAYLINUX", "Volume label for the boot partition")
	// flag.StringVar(&compression, "compression", "gz", "Compression format for initramfs: none|gz|xz|zstd")
	flag.Parse()

	// build the image
	err := image.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}
