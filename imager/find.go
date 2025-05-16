package imager

import (
	"errors"
	"path/filepath"
)

func findSingleFile(pattern string) (string, error) {
    files, err := filepath.Glob(pattern)
    if err != nil {
        return "", err
    }
    if len(files) > 1 {
        return "", errors.New("more than one file found")
    }
    if len(files) == 0 {
        return "", nil
    }
    return files[0], nil
}
