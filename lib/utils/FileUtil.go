package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func IsDir(name string) bool {
	if stat, err := os.Stat(name); !os.IsNotExist(err) {
		return stat.IsDir()
	}
	return false
}

func FindWorkingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}

func DetectRelativePath(p string) (string, error) {
	cwd := FindWorkingDir()
	if len(cwd) > 0 {
		rel, err := filepath.Rel(cwd, p)
		if err == nil {
			return rel, nil
		}
	}
	return p, fmt.Errorf("Relative path detecting failed")
}
