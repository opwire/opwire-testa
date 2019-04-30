package utils

import (
	"fmt"
	"path/filepath"
	"github.com/opwire/opwire-testa/lib/storage"
)

func IsDir(name string) bool {
	fs := storage.GetFs()
	if stat, err := fs.Stat(name); !fs.IsNotExist(err) {
		return stat.IsDir()
	}
	return false
}

func FindWorkingDir() string {
	fs := storage.GetFs()
	dir, err := fs.Getwd()
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

func DetectRelativePaths(p []string) []string {
	return Map(p, func(dir string, i int) string {
		newPath, err := DetectRelativePath(dir)
		if err == nil {
			return newPath
		}
		return dir
	})
}