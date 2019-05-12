package tag

import(
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestManager_Initialize(t *testing.T) {
	ref, err := NewManager(nil)
	assert.NotNil(t, ref)
	assert.Nil(t, err)
	t.Run("Default", func(t *testing.T) {
		TESTCASES := []struct {
			tagExpression []string
			inclusiveTags []string
			exclusiveTags []string
		}{
			// try matching as a suffix string
			{
				tagExpression: nil,
				inclusiveTags: []string{},
				exclusiveTags: []string{},
			},
			{
				tagExpression: []string{ "" },
				inclusiveTags: []string{},
				exclusiveTags: []string{},
			},
			{
				tagExpression: []string{ "abc, +def" },
				inclusiveTags: []string{ "abc", "def" },
				exclusiveTags: []string{},
			},
			{
				tagExpression: []string{ "abc, -def", "+abc,-def-ghi,-xyz, +xyz" },
				inclusiveTags: []string{ "abc", "xyz" },
				exclusiveTags: []string{ "def", "def-ghi", "xyz" },
			},
			{
				tagExpression: []string{ "abc, +, -, +abc-1, " },
				inclusiveTags: []string{ "abc", "abc-1" },
				exclusiveTags: []string{},
			},
		}
		for _, TEST := range TESTCASES {
			ref.Initialize(TEST.tagExpression)
			assert.Equal(t, TEST.inclusiveTags, ref.inclusiveTags)
			assert.Equal(t, TEST.exclusiveTags, ref.exclusiveTags)
		}
	})
}

func TestManager_IsActive(t *testing.T) {
	ref, err := NewManager(nil)
	assert.NotNil(t, ref)
	assert.Nil(t, err)
	t.Run("Default", func(t *testing.T) {
		TESTCASES := []struct {
			tagExpression []string
			tags []string
			ok bool
			mark map[string]int8
		}{
			{
				tagExpression: []string{},
				tags: []string{},
				ok: true,
				mark: map[string]int8{},
			},
			{
				tagExpression: []string{},
				tags: []string{ "abc", "def" },
				ok: true,
				mark: map[string]int8{},
			},
			{
				tagExpression: []string{ "+abc, -xyz" },
				tags: []string{ "abc", "def" },
				ok: true,
				mark: map[string]int8{
					"abc": 1,
				},
			},
			{
				tagExpression: []string{ "+abc, -def, -xyz" },
				tags: []string{ "abc", "def" },
				ok: false,
				mark: map[string]int8{
					"def": -1,
				},
			},
		}
		for _, TEST := range TESTCASES {
			ref.Initialize(TEST.tagExpression)
			ok, mark := ref.IsActive(TEST.tags)
			assert.Equal(t, TEST.ok, ok)
			assert.Equal(t, TEST.mark, mark)
		}
	})
}
