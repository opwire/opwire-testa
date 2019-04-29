package script

type Source interface {
	GetTestDirs() []string
	GetTestName() string
	GetConditionalTags() []string
}

func NewSource(opts Source) (ref *SourceBuffer, err error) {
	buf := &SourceBuffer{}
	if opts != nil {
		buf.TestDirs = opts.GetTestDirs()
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
