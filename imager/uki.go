package imager

import (
	"os"
	"path/filepath"

	"github.com/sprat/claylinux/imager/efi"
	"github.com/sprat/claylinux/imager/uki"
)

// Build the Unified Kernel Image
func (i Image) buildUKI() (string, error) {
	stubPath := efi.GetStubPath()
	stub, err := os.ReadFile(stubPath)
	if err != nil {
		return "", err
	}
	builder := uki.NewBuilder(stub)

	// cmdline
	cmdline, err := i.getCmdline()
	if err != nil {
		return "", err
	}
	builder.AddSection(".cmdline", []byte(cmdline))

	// initramfs
	// TODO: build in memory?
	initRamFsPath, err := i.buildInitRamFs()
	if err != nil {
		return "", err
	}
	initRamFs, err := os.ReadFile(initRamFsPath)
	if err != nil {
		return "", err
	}
	builder.AddSection(".initrd", initRamFs)

	// kernel release
	kernelRelease, err := i.getKernelRelease()
	if err != nil {
		return "", err
	}
	builder.AddSection(".uname", []byte(kernelRelease))

	// os release
	osRelease, err := i.getOSRelease()
	if err != nil {
		return "", err
	}
	builder.AddSection(".osrel", []byte(osRelease))

	// TODO: include ucode
	// ucodePath, err := findSingleFile(filepath.Join(i.RootFsDir, "boot", "*-ucode.img"))
	// builder.AddSection(".ucode", ucodePath)
	// builder.AddSection(".dtb", ...)
	// builder.AddSection(".splash", ...)

	// kernel
	// Should be the last section because everything after can be overwritten by in-place kernel decompression
	kernelPath, err := i.findKernelPath()
	if err != nil {
		return "", err
	}
	kernel, err := os.ReadFile(kernelPath)
	if err != nil {
		return "", err
	}
	builder.AddSection(".linux", kernel)

	outputPath := filepath.Join(i.BuildDir, "claylinux.efi")
	err = builder.Write(outputPath)
	if err != nil {
		return "", err
	}

	return outputPath, nil
}
