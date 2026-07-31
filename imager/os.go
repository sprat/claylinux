package imager

import (
	"os"
	"path/filepath"
)

// Get the kernel release information
func (i Image) getOSRelease() (string, error) {
	path := filepath.Join(i.RootFsDir, "etc", "os-release")
	data, err := os.ReadFile(path)
	return string(data), err
}
