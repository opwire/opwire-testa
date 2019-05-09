package utils

import(
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestStandardizeVersion(t *testing.T) {
	TESTCASES := []struct {
		version string
		result string
	}{
		{
			version: "v1",
			result: "1",
		},
		{
			version: "2",
			result: "2",
		},
		{
			version: "v1.3",
			result: "1.3",
		},
		{
			version: "1.4",
			result: "1.4",
		},
		{
			version: "v1.5.1",
			result: "1.5.1",
		},
		{
			version: "1.6.1",
			result: "1.6.1",
		},
		{
			version: "v1.5.1-hotfix1",
			result: "1.5.1-hotfix1",
		},
		{
			version: "1.6.1-fixbug",
			result: "1.6.1-fixbug",
		},
	}
	for _, TEST := range TESTCASES {
		assert.Equal(t, StandardizeVersion(TEST.version), TEST.result)
	}
}

func TestStandardizeTagLabel(t *testing.T) {
	TESTCASES := []struct {
		label string
		tag string
		hasError bool
	}{
		{
			label: `valid-tag_01`,
			tag: `valid-tag_01`,
			hasError: false,
		},
		{
			label: `Windows\User`,
			tag: `Windows_User`,
			hasError: false,
		},
		{
			label: `0`,
			tag: `0`,
			hasError: true,
		},
		{
			label: `01-invalid`,
			tag: `01-invalid`,
			hasError: true,
		},
		{
			label: `-invalid`,
			tag: `-invalid`,
			hasError: true,
		},
		{
			label: `_invalid`,
			tag: `_invalid`,
			hasError: true,
		},
	}
	for _, TEST := range TESTCASES {
		tag, err := StandardizeTagLabel(TEST.label)
		assert.Equal(t, tag, TEST.tag)
		if TEST.hasError {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}
