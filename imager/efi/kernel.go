package efi

import "path/filepath"

func findKernel(rootfsDir string) (string, error) {
	kernelFiles, err := filepath.Glob(rootfsDir + "/boot/vmlinu*")
	if err != nil {
		return "", err
	}
	return kernelFiles[0], nil
}
