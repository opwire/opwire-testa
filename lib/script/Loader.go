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

func (l *Loader) ReadDir(sourceDir string, ext string) ([]string, error) {
	files := []string{}
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
