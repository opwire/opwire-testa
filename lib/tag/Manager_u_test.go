package tag

import(
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestManager_Initialize(t *testing.T) {
	ref, err := NewManager(nil)
	assert.NotNil(t, ref)
	assert.Nil(t, err)
	t.Run("Ok", func(t *testing.T) {
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
			assert.Equal(t, ref.inclusiveTags, TEST.inclusiveTags)
			assert.Equal(t, ref.exclusiveTags, TEST.exclusiveTags)
		}
	})
}