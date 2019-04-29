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
	Tags []string
	TestFile string
	TestName string
}

func (a *SourceBuffer) GetTestDirs() []string {
	return a.TestDirs
}

func (a *SourceBuffer) GetTestName() string {
	return a.TestName
}

func (a *SourceBuffer) GetConditionalTags() []string {
	return a.Tags
}
