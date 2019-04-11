package script

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"gopkg.in/yaml.v2"
	"github.com/opwire/opwire-qakit/lib/storages"
)

type Loader struct {
}

func NewLoader() *Loader {
	l := new(Loader)
	return l
}

func (l *Loader) LoadScript(filePath string) (*Descriptor, error) {
	descriptor := &Descriptor{}

	fs := storages.GetFs()
	file, err := fs.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	parser := yaml.NewDecoder(file)
	err = parser.Decode(descriptor)
	if err != nil {
		return nil, err
	}

	return descriptor, nil
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

type Descriptor struct {
	Scenarios []Scenario `yaml:"scenarios"`
}

type Scenario struct {
	Title *string `yaml:"title"`
	Method *string `yaml:"method"`
	Path *string `yaml:"path"`
}

// type HttpHeader struct {
// 	Name *string `yaml:"name"`
// 	Value *string `yaml:"value"`
// }

// type HttpRequest struct {
// 	Headers []HttpHeader `yaml:"headers"`
// 	Body *string `yaml:"body"`
// }

// type HttpMeasure struct {
// 	Headers []HttpHeader `yaml:"headers"`
// 	Body *string `yaml:"body"`
// }
