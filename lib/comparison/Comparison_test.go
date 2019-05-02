package comparison

import(
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestConvertNumber(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		assert.True(t, IsEqual(nil, nil))

		var b bool = false
		assert.True(t, IsEqual(false, b))

		var x int = 1024
		var y float64 = 1024
		assert.True(t, IsEqual(x, y))
	})
}
