package utils

import (
	"net/url"
	"path"
)

func UrlJoin(pdp string, basePath string) (string, error) {
	u, err := url.Parse(pdp)
	if err != nil {
		return path.Join(pdp, basePath), err
	}
	if len(basePath) > 0 {
		u.Path = path.Join(u.Path, basePath)
	}
	s := u.String()
	return s, nil
}
