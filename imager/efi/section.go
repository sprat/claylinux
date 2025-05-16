package efi

import (
	"fmt"
	"os"
)

type Section struct {
	Name string
	Path string
}

func (s *Section) Size() (int64, error) {
	file, err := os.Open(s.Path)
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

func getSectionsArgs(sections []Section, firstAddress uint64, alignment uint64) ([]string, error) {
	var args []string
	nextAddress := firstAddress
	for _, section := range sections {
		args = append(
			args,
			"--add-section",
			fmt.Sprintf("%s=%s", section.Name, section.Path),
			"--change-section-vma",
			fmt.Sprintf("%s=0x%x", section.Name, nextAddress),
		)
		size, err := section.Size()
		if err != nil {
			return args, err
		}
		nextAddress = align(nextAddress + uint64(size), alignment)
	}
	return args, nil
}
