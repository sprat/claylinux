package uki

import (
	"fmt"
	"os"

	"github.com/soyum2222/editPE"
)

// see:
// - https://secanablog.wordpress.com/2020/06/09/how-to-add-a-section-to-a-pe-file/
// - https://github.com/secana/PeNet/blob/master/src/PeNet/Editor/Section.cs
// - https://www.debugxp.com/posts/ExtendingPESections/
// - https://github.com/willscott/pefile-go
// - https://github.com/soyum2222/editPE

type Builder struct {
	pe editPE.PE
}

func NewBuilder(stub []byte) *Builder {
	pe := editPE.PE{}
	pe.Parse(stub)

	for index, header := range pe.ImageSectionHeaders {
		fmt.Printf("Section %x: %s (size %x / %x at addr %x / offset %x)\n", index, header.Name, header.SizeOfRawData, header.PhysicalAddressOrVirtualSize, header.VirtualAddress, header.PointerToRawData)
	}

	return &Builder{
		pe: pe,
	}
}

func (b *Builder) AddSection(name string, data []byte) {
	fmt.Printf("Adding section %s of size %x\n", name, len(data))
	b.pe.AddSection(name, uint32(len(data)))

	// TODO: there's something wrong here
	index := len(b.pe.ImageSectionHeaders) - 1
	header := b.pe.ImageSectionHeaders[index]
	offset := header.PointerToRawData
	fmt.Printf("Section %x: %s (size %x / %x at addr %x / offset %x)\n", index, header.Name, header.SizeOfRawData, header.PhysicalAddressOrVirtualSize, header.VirtualAddress, header.PointerToRawData)
	copy(data, b.pe.Raw[offset:])
}

func (b *Builder) Write(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR | os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(b.pe.Raw)
	if err != nil {
		return err
	}

	return nil
	/*
	stub, err := pe.NewFile(b.stubReader)
	if err != nil {
		return err
	}
	*/

	/*
	fmt.Printf("File: %+v\n", stub)
	numberOfSections := stub.NumberOfSections
	//sizeOfOptionalHeader := stub.SizeOfOptionalHeader
	//var base uint64  // should be uint32 on 32-bits platforms, but it has no consequence here
	var sectionAlignment uint64
	var fileAlignment uint64
	var sizeOfImage uint32
	var sizeOfHeaders uint32
	if optionalHeader, ok := stub.OptionalHeader.(*pe.OptionalHeader64); ok {
		fmt.Printf("Optional header: %+v\n", optionalHeader)
		//base = uint64(optionalHeader.ImageBase)
		sizeOfImage = optionalHeader.SizeOfImage
		sizeOfHeaders = optionalHeader.SizeOfHeaders
		sectionAlignment = uint64(optionalHeader.SectionAlignment)
		fileAlignment = uint64(optionalHeader.FileAlignment)
	} else if optionalHeader, ok := stub.OptionalHeader.(*pe.OptionalHeader32); ok {
		fmt.Printf("Optional header: %+v\n", optionalHeader)
		//base = uint64(optionalHeader.ImageBase)
		sizeOfImage = optionalHeader.SizeOfImage
		sizeOfHeaders = optionalHeader.SizeOfHeaders
		sectionAlignment = uint64(optionalHeader.SectionAlignment)
		fileAlignment = uint64(optionalHeader.FileAlignment)
	} else {
		return errors.New("optional header should be present in EFI files")
	}

	// sections are sorted by increasing virtual address
	// so we just have to take the last one to find the next available address
	lastSection := stub.Sections[len(stub.Sections) - 1]
	fmt.Printf("Last Section: %+v\n", lastSection)

	// append bytes at end of file
	// pad with 0 until we reach the file alignment

	fmt.Println(sizeOfImage)
	fmt.Println(sizeOfHeaders)
	fmt.Println(numberOfSections)
	fmt.Println(sectionAlignment)
	fmt.Println(fileAlignment)

	nextVirtualAddress := alignAddress(
		base + uint64(lastSection.VirtualAddress) + uint64(lastSection.VirtualSize),
		sectionAlignment,
	)

	// nextPhysicalAddress := lastSection.Offset + lastSection.Size
	header = SectionHeader32 {
		Name                 [8]uint8
		VirtualSize          uint32
		VirtualAddress       uint32
		Size                 uint32
		Offset               uint32
		PointerToRelocations uint32
		PointerToLineNumbers uint32
		NumberOfRelocations  uint16
		NumberOfLineNumbers  uint16
		Characteristics      uint32
	}

	//p.nextVirtualAddress = alignAddress(p.nextVirtualAddress + uint64(size), p.sectionAlignment)
	return errors.New("Not finished")
	*/
}
