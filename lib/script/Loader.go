package script

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"gopkg.in/yaml.v2"
	"github.com/opwire/opwire-testa/lib/engine"
	"github.com/opwire/opwire-testa/lib/storages"
)

type LoaderOptions interface {}

type Loader struct {
}

func NewLoader(opts LoaderOptions) (*Loader, error) {
	l := new(Loader)
	return l, nil
}

func (l *Loader) LoadScripts(sourceDirs []string) (map[string]*Descriptor, error) {
	locators, _ := l.ReadDirs(sourceDirs, ".yml")
	descriptors, _ := l.LoadFiles(locators)
	return descriptors, nil
}

func (l *Loader) LoadFiles(locators []*Locator) (descriptors map[string]*Descriptor, err error) {
	descriptors = make(map[string]*Descriptor, 0)
	for _, locator := range locators {
		descriptor, err := l.LoadFile(locator)
		if err == nil {
			descriptors[locator.FullPath] = descriptor
		}
	}
	return descriptors, nil
}

func (l *Loader) LoadFile(locator *Locator) (*Descriptor, error) {
	if locator == nil {
		return nil, fmt.Errorf("Descriptor must not be nil")
	}

	descriptor := &Descriptor{}

	fs := storages.GetFs()
	file, err := fs.Open(locator.FullPath)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	parser := yaml.NewDecoder(file)
	err = parser.Decode(descriptor)
	if err != nil {
		return nil, err
	}

	descriptor.Locator = locator

	return descriptor, nil
}

func (l *Loader) ReadDirs(sourceDirs []string, ext string) (locators []*Locator, err error) {
	locators = make([]*Locator, 0)
	for _, sourceDir := range sourceDirs {
		items, err := l.ReadDir(sourceDir, ext)
		if err == nil {
			locators = append(locators, items...)
		}
	}
	return locators, nil
}

func (l *Loader) ReadDir(sourceDir string, ext string) ([]*Locator, error) {
	locators := make([]*Locator, 0)
	err := filepath.Walk(sourceDir, func(path string, f os.FileInfo, err error) error {
		if err == nil && !f.IsDir() {
			r, err := regexp.MatchString(ext, f.Name())
			if err == nil && r {
				locator := &Locator{}
				locator.FullPath = path
				locator.Home = sourceDir
				locator.Path = strings.TrimPrefix(path, sourceDir)
				locators = append(locators, locator)
			}
		}
		return nil
	})
	return locators, err
}

type Locator struct {
	FullPath string
	Home string
	Path string
	Error error
}

type Descriptor struct {
	Locator *Locator
	Scenarios []*engine.Scenario `yaml:"scenarios"`
}
