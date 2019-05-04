package comparison

import(
	"testing"
	"github.com/opwire/opwire-testa/lib/testutils"
	"github.com/stretchr/testify/assert"
)

func TestIsEqualTo(t *testing.T) {
	t.Run("Numbers comparison", func(t *testing.T) {
		assert.True(t, testutils.GetFirstResult_bool(IsEqualTo(nil, nil)))

		var b bool = false
		assert.True(t, testutils.GetFirstResult_bool(IsEqualTo(false, b)))

		var x int = 1024
		var y float64 = 1024
		assert.True(t, testutils.GetFirstResult_bool(IsEqualTo(x, y)))
	})
}
