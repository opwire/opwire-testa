package comparison

import(
	"testing"
	"reflect"
	"github.com/stretchr/testify/assert"
)

func TestIsZero(t *testing.T) {
	t.Run("Test zero values", func(t *testing.T) {
		assert.True(t, IsZero(reflect.ValueOf(nil)))
		assert.True(t, IsZero(reflect.ValueOf(0)))
		assert.True(t, IsZero(reflect.ValueOf(0.0)))
		assert.True(t, IsZero(reflect.ValueOf(false)))
		assert.True(t, IsZero(reflect.ValueOf("")))

		var n int
		assert.True(t, IsZero(reflect.ValueOf(n)))

		var v interface{}
		assert.True(t, IsZero(reflect.ValueOf(v)))
	})

	t.Run("Test non-zero values", func(t *testing.T) {
		assert.False(t, IsZero(reflect.ValueOf(1)))
		assert.False(t, IsZero(reflect.ValueOf(0.1)))
		assert.False(t, IsZero(reflect.ValueOf(true)))
		assert.False(t, IsZero(reflect.ValueOf("a")))
		v := struct{ f string }{ f: "a" }
		assert.False(t, IsZero(reflect.ValueOf(v)))
	})
}
