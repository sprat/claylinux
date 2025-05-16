package imager

import "path/filepath"

func (i Image) getCmdline() string {
	// TODO: space-separated
	return filepath.Join(i.RootFsDir, "boot", "cmdline")
}
