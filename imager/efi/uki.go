package efi

import (
	"debug/pe"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
)

func BuildUKI(rootfsDir, buildDir string) (string, error) {
	stubPath := getStubPath()
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
	firstAddress := align(base + uint64(lastSection.VirtualAddress) + uint64(lastSection.VirtualSize), alignment)
	//fmt.Printf("firstAddress=0x%x\n", firstAddress)

	// prepare objcopy arguments
	kernelPath, err := findKernel(rootfsDir)
	if err != nil {
		return "", err
	}
	// initRamFs = efi.BuildInitramfs(rootfsDir, buildDir)

	sections := []Section{
		{".osrel", rootfsDir + "/etc/os-release"},
		// {".uname", kernel-release},
		// {".cmdline", cmdline},
		// {".initrd", initramfs},
		// {".ucode", ...},
		{".linux", kernelPath},  // should be the last section, because everything after can be overwritten by inplace decompression
	}

	args, err := getSectionsArgs(sections, firstAddress, alignment)
	if err != nil {
		return "", err
	}

	// run objcopy
	outputPath := path.Join(buildDir, "claylinux.efi")
	args = append(args, stubPath, outputPath)
	fmt.Printf("objcopy args: %+v", args)
	cmd := exec.Command("objcopy", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return "", err
	}
	return outputPath, nil
}

/*
objcopy \
--add-section .osrel=os-release --change-section-vma .osrel=0x20000 \
--add-section .cmdline=cmdline.txt --change-section-vma .cmdline=0x30000 \
--add-section .dtb=devicetree.dtb --change-section-vma .dtb=0x40000 \
--add-section .splash=splash.bmp --change-section-vma .splash=0x100000 \
--add-section .linux=vmlinux --change-section-vma .linux=0x2000000 \
--add-section .initrd=initrd.cpio --change-section-vma .initrd=0x3000000 \
/usr/lib/systemd/boot/efi/linuxx64.efi.stub \
foo-unsigned.efi

sbsign \
--key mykey.pem \
--cert mykey.crt \
--output foo.efi \
foo-unsigned.efi
*/


/*
func buildEFI() {
	os.Chdir(buildDir)
	fmt.Println("Building the EFI executable")
	buildInitramfs()
	size := getFileSizeMib("initramfs")
	fmt.Printf("The size of the initramfs is: %d MiB\n", size)
	kernel := runOutput("find", "/system/boot", "-name", "vmlinu*", "-print")
	run("sh", "-c", "tr '\\n' ' ' </system/boot/cmdline >cmdline")
	run("sh", "-c", "basename /system/lib/modules/* >kernel-release")

	// Compose the UKI section file
	ukiStanza := fmt.Sprintf(
		".osrel /system/etc/os-release\n.uname kernel-release\n.cmdline cmdline\n.initrd initramfs\n.linux %s\n", kernel)
	cmd := exec.Command("build_uki")
	cmd.Stdin = strings.NewReader(ukiStanza)
	if err := cmd.Run(); err != nil {
		die("build_uki failed: " + err.Error())
	}

	run("find", ".", "!", "-name", "*.efi", "-delete")
}
*/
