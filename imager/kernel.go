package imager

import "path/filepath"

func findKernel(rootfsDir string) (string, error) {
	pattern := filepath.Join(rootfsDir, "boot", "vmlinu*")
	return findSingleFile(pattern)
}

func getKernelRelease(rootfsDir string) (string, error) {
	pattern := filepath.Join(rootfsDir, "lib", "modules", "*")
	modulesBase, err := findSingleFile(pattern)
	if err != nil {
		return "", err
	}
	return filepath.Base(modulesBase), nil
}
