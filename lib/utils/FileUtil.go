package utils

import (
	"os"
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
