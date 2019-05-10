package bootstrap

import(
	"path/filepath"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/opwire/opwire-testa/lib/storage"
)

func Test_checkFilePathMatchPattern(t *testing.T) {
	fs := storage.GetFs()
	cwd, err := fs.Getwd()
	if err != nil {
		t.Errorf("Cannot get current working directory path")
		return
	}
	TESTCASES := []struct {
		filePath string
		pattern string
		matched bool
	}{
		// try matching as a suffix string
		{
			filePath: filepath.Join(cwd, "test/relative/path/to/file.yml"),
			pattern: "to/file.yml",
			matched: true,
		},
		{
			filePath: filepath.Join(cwd, "test/relative/path/to/file.yml"),
			pattern: "to/file.json",
			matched: false,
		},
		{
			filePath: filepath.Join(cwd, "test/relative/path/to/file.yml"),
			pattern: "from/file.yml",
			matched: false,
		},
		// try matching using filepath
		{
			filePath: filepath.Join(cwd, "test/relative/path/to/file.yml"),
			pattern: "test/relative/*",
			matched: true,
		},
		{
			filePath: filepath.Join(cwd, "test/relative/path/to/file.yml"),
			pattern: "test/relative/path/to",
			matched: true,
		},
		{
			filePath: filepath.Join(cwd, "test/relative/path/to/file.yml"),
			pattern: "test/relative/path/*/file.yml",
			matched: true,
		},
		// try matching with a regexp
		{
			filePath: filepath.Join(cwd, "test/relative/path/to/file.yml"),
			pattern: ".*/relative/.*",
			matched: true,
		},
	}
	for _, TEST := range TESTCASES {
		assert.Equal(t, checkFilePathMatchPattern(TEST.filePath, TEST.pattern), TEST.matched)
	}
}
