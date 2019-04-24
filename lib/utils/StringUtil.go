package utils

import (
	"regexp"
)

const VERSION_PATTERN string = `[v]?((\\d+\\.)?(\\d+\\.)?(\\*|\\d+))`

var re = regexp.MustCompile(VERSION_PATTERN)

func StandardizeVersion(version string) string {
	return re.ReplaceAllString(version, `$1`)
}
