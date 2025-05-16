package imager

import (
	"os"
	"path/filepath"

	"github.com/sprat/claylinux/imager/efi"
)

// Build the Unified Kernel Image
func (i Image) buildUKI() (string, error) {
	stubPath := efi.GetStubPath()
	peFile, err := NewPEFile(stubPath)
	if err != nil {
		return "", err
	}

	// cmdline
	cmdline := i.getCmdline()
	err = peFile.AddSection(".cmdline", cmdline)
	if err != nil {
		return "", err
	}

	// initramfs
	initRamFsPath, err := i.buildInitRamFs()
	if err != nil {
		return "", err
	}
	err = peFile.AddSection(".initrd", initRamFsPath)
	if err != nil {
		return "", err
	}

	// kernel release
	kernelRelease, err := i.getKernelRelease()
	if err != nil {
		return "", err
	}
	kernelReleasePath := filepath.Join(i.BuildDir, "kernel-release")
	os.WriteFile(kernelReleasePath, []byte(kernelRelease), 0644)
	err = peFile.AddSection(".uname", kernelReleasePath)
	if err != nil {
		return "", err
	}

	// os release
	osReleasePath := filepath.Join(i.RootFsDir, "etc", "os-release")
	err = peFile.AddSection(".osrel", osReleasePath)
	if err != nil {
		return "", err
	}

	// TODO: include ucode
	// ucodePath, err := findSingleFile(filepath.Join(i.RootFsDir, "boot", "*-ucode.img"))
	// peFile.AddSection("".ucode", ucodePath)

	// peFile.AddSection("".dtb", ...)
	// peFile.AddSection("".splash", ...)

	// kernel
	// Should be the last section because everything after can be overwritten by in-place kernel decompression
	kernelPath, err := i.findKernel()
	if err != nil {
		return "", err
	}
	err = peFile.AddSection(".linux", kernelPath)
	if err != nil {
		return "", err
	}

	// run objcopy
	outputPath := filepath.Join(i.BuildDir, "claylinux.efi")
	err = peFile.Finalize(outputPath)
	if err != nil {
		return "", err
	}
	return outputPath, nil
}
