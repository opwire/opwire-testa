package utils

import(
	"testing"
	"github.com/stretchr/testify/assert"
)

func GetOut(out string, err error) string {
	return out
}

func GetErr(out string, err error) error {
	return err
}

func TestUrlJoin(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		assert.Equal(t, GetOut(UrlJoin("http://localhost", "")), "http://localhost")
		assert.Equal(t, GetOut(UrlJoin("http://localhost/", "")), "http://localhost/")
		assert.Equal(t, GetOut(UrlJoin("http://localhost", "/")), "http://localhost/")
		assert.Equal(t, GetOut(UrlJoin("http://localhost", "foo")), "http://localhost/foo")
		assert.Equal(t, GetOut(UrlJoin("http://localhost/", "foo")), "http://localhost/foo")
		assert.Equal(t, GetOut(UrlJoin("http://localhost", "/foo")), "http://localhost/foo")
		assert.Equal(t, GetOut(UrlJoin("http://localhost/", "/foo")), "http://localhost/foo")
		assert.Equal(t, GetOut(UrlJoin("http://localhost", "foo/bar")), "http://localhost/foo/bar")
	})
}
