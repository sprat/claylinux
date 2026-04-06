package imager

import (
	"os"
	"path/filepath"
	"strings"
)

func (i Image) getCmdline() (string, error) {
	path := filepath.Join(i.RootFsDir, "boot", "cmdline")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	cmdline := strings.Replace(string(data), "\n", " ", -1)
	return cmdline, nil
}
