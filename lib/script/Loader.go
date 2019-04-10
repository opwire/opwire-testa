package script

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Loader struct {
}

func NewLoader() *Loader {
	l := new(Loader)
	return l
}

func (l *Loader) ReadDirs(sourceDirs []string, ext string) (files []string, err error) {
	files = []string{}
	for _, sourceDir := range sourceDirs {
		files, _ = l.appendDir(files, sourceDir, ext)
	}
	return files, nil
}

func (l *Loader) ReadDir(sourceDir string, ext string) ([]string, error) {
	return l.appendDir(nil, sourceDir, ext)
}

func (l *Loader) appendDir(files []string, sourceDir string, ext string) ([]string, error) {
	if files == nil {
		files = []string{}
	}
	err := filepath.Walk(sourceDir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(ext, f.Name())
			if err == nil && r {
				files = append(files, strings.TrimPrefix(path, sourceDir))
			}
		}
		return nil
	})
	return files, err
}
