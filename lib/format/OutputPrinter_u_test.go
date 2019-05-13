package format

import(
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestColorlessPen(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		assert.Equal(t, ColorlessPen("abc"), "abc")
		assert.Equal(t, ColorlessPen("abc", 123), "abc 123")
		assert.Equal(t, ColorlessPen("abc", true, 0), "abc true 0")
		assert.Equal(t, ColorlessPen("abc", 3.14, nil), "abc 3.14 <nil>")
		assert.Equal(t, ColorlessPen("a", "", "z"), "a z")
		assert.Equal(t, ColorlessPen("", "", "", ""), "")
	})
}
