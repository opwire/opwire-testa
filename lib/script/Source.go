package script

import (
	"path/filepath"
	"github.com/opwire/opwire-testa/lib/utils"
)

type Source interface {
	GetTestDirs() []string
	GetInclFiles() []string
	GetExclFiles() []string
	GetTestName() string
	GetConditionalTags() []string
}

func NewSource(opts Source) (ref *SourceBuffer, err error) {
	buf := &SourceBuffer{}
	if opts != nil {
		buf.TestDirs = opts.GetTestDirs()
		buf.InclFiles = opts.GetInclFiles()
		buf.ExclFiles = opts.GetExclFiles()
		buf.TestName = opts.GetTestName()
		buf.Tags = opts.GetConditionalTags()
	}
	return buf, err
}

type SourceBuffer struct {
	TestDirs []string
	InclFiles []string
	ExclFiles []string
	TestName string
	Tags []string
}

func (a *SourceBuffer) GetTestDirs() []string {
	a.TestDirs = initDefaultDirs(a.TestDirs)
	return a.TestDirs
}

func (a *SourceBuffer) GetInclFiles() []string {
	return a.InclFiles
}

func (a *SourceBuffer) GetExclFiles() []string {
	return a.ExclFiles
}

func (a *SourceBuffer) GetTestName() string {
	return a.TestName
}

func (a *SourceBuffer) GetConditionalTags() []string {
	return a.Tags
}

func initDefaultDirs(testDirs []string) []string {
	if testDirs == nil || len(testDirs) == 0 {
		testDir := filepath.Join(utils.FindWorkingDir(), "tests")
		if utils.IsDir(testDir) {
			testDirs = []string{testDir}
		}
	}
	return testDirs
}
