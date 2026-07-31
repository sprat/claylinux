package imager

import "path/filepath"

// Find the Linux kernel file in the root filesystem
func (i Image) findKernel() (string, error) {
	pattern := filepath.Join(i.RootFsDir, "boot", "vmlinu*")
	return findSingleFile(pattern)
}

// Get the kernel release information
func (i Image) getKernelRelease() (string, error) {
	pattern := filepath.Join(i.RootFsDir, "lib", "modules", "*")
	modulesBase, err := findSingleFile(pattern)
	if err != nil {
		return "", err
	}
	return filepath.Base(modulesBase), nil
}
