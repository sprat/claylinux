package imager

import (
	"debug/pe"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

type PEFile struct {
	Path string
	alignment uint64
	nextAddress uint64
	args []string
}

func NewPEFile(path string) (*PEFile, error) {
	peFile, err := pe.Open(path)
	if err != nil {
		return nil, err
	}

	defer peFile.Close()

	var base uint64  // should be uint32 on 32-bits platforms, but it has no consequence here
	var alignment uint64

	switch optionalHeader := peFile.OptionalHeader.(type) {
	case *pe.OptionalHeader32:
		base = uint64(optionalHeader.ImageBase)
		alignment = uint64(optionalHeader.SectionAlignment)
	case *pe.OptionalHeader64:
		base = uint64(optionalHeader.ImageBase)
		alignment = uint64(optionalHeader.SectionAlignment)
	default:
		return nil, errors.New("optional header should be present in EFI files")
	}

	// sections are sorted by increasing virtual address
	// so we just have to take the last one to find the next available address
	lastSection := peFile.Sections[len(peFile.Sections) - 1]
	nextAddress := alignAddress(base + uint64(lastSection.VirtualAddress) + uint64(lastSection.VirtualSize), alignment)

	return &PEFile{
		Path: path,
		alignment: alignment,
		nextAddress: nextAddress,
		args: []string{},
	}, nil
}

func (p *PEFile) AddSection(name string, path string) error {
	size, err := fileSize(path)
	if err != nil {
		return err
	}

	p.args = append(
		p.args,
		"--add-section",
		fmt.Sprintf("%s=%s", name, path),
		"--change-section-vma",
		fmt.Sprintf("%s=0x%x", name, p.nextAddress),
	)
	p.nextAddress = alignAddress(p.nextAddress + uint64(size), p.alignment)
	return nil
}

func (p *PEFile) Finalize(path string) error {
	p.args = append(p.args, p.Path, path)
	cmd := exec.Command("objcopy", p.args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func fileSize(path string) (int64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}
