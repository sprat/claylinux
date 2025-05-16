package imager

import (
	"debug/pe"
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sprat/claylinux/imager/efi"
)

// Build the Unified Kernel Image
func (i Image) buildUKI() (string, error) {
	stubPath := efi.GetStubPath()
	peFile, err := pe.Open(stubPath)
	if err != nil {
		return "", err
	}
	defer peFile.Close()

	//fmt.Printf("SizeOfOptionalHeader=0x%x, NumberOfSections=%v\n", peFile.SizeOfOptionalHeader, peFile.NumberOfSections)

	var base uint64  // should be uint32 on 32-bits platforms, but it has no consequence here
	var alignment uint64
	if optionalHeader, ok := peFile.OptionalHeader.(*pe.OptionalHeader64); ok {
		base = uint64(optionalHeader.ImageBase)
		alignment = uint64(optionalHeader.SectionAlignment)
	} else if optionalHeader, ok := peFile.OptionalHeader.(*pe.OptionalHeader32); ok {
		base = uint64(optionalHeader.ImageBase)
		alignment = uint64(optionalHeader.SectionAlignment)
	} else {
		return "", errors.New("optional header should be present in EFI files")
	}

	//fmt.Printf("base=0x%x, alignment=0x%x\n", base, alignment)

	// sections are sorted by increasing virtual address
	// so we just have to take the last one to find the next available address
	lastSection := peFile.Sections[len(peFile.Sections) - 1]
	firstAddress := alignAddress(base + uint64(lastSection.VirtualAddress) + uint64(lastSection.VirtualSize), alignment)
	//fmt.Printf("firstAddress=0x%x\n", firstAddress)

	// prepare objcopy arguments
	kernelPath, err := i.findKernel()
	if err != nil {
		return "", err
	}

	cmdlinePath := i.getCmdline()

	initRamFsPath, err := i.buildInitRamFs()
	if err != nil {
		return "", err
	}

	kernelRelease, err := i.getKernelRelease()
	if err != nil {
		return "", err
	}

	kernelReleasePath := filepath.Join(i.BuildDir, "kernel-release")
	os.WriteFile(kernelReleasePath, []byte(kernelRelease), 0644)

	osReleasePath := filepath.Join(i.RootFsDir, "etc", "os-release")

	//ucodePath, err := findSingleFile(filepath.Join(rootfsDir, "boot", "*-ucode.img"))

	// Linux should be the last section because everything after can be overwritten by in-place
	// kernel decompression
	sections := []Section{
		{".cmdline", cmdlinePath},
		{".initrd", initRamFsPath},
		{".osrel", osReleasePath},
		{".uname", kernelReleasePath},
		// {".ucode", ucodePath},
		// {".dtb", ...},
		// {".splash", ...},
		{".linux", kernelPath},
	}

	args, err := getSectionsArgs(sections, firstAddress, alignment)
	if err != nil {
		return "", err
	}

	// run objcopy
	outputPath := filepath.Join(i.BuildDir, "claylinux.efi")
	args = append(args, stubPath, outputPath)
	cmd := exec.Command("objcopy", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return "", err
	}
	return outputPath, nil
}
