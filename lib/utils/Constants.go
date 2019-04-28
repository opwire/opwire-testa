package utils

import (
	"regexp"
)

const TAG_PATTERN string = `[a-zA-Z][a-zA-Z0-9]*([_-][a-zA-Z0-9]*)*`
const TIME_RFC3339 string = `([0-9]+)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|60)(\\.[0-9]+)?(([Zz])|([\\+|\\-]([01][0-9]|2[0-3]):[0-5][0-9]))`

const TEST_CASE_TITLE_PATTERN string = `[a-zA-Z][\w\-. ]*`
var TEST_CASE_TITLE_REGEXP *regexp.Regexp = regexp.MustCompile(`^` + TEST_CASE_TITLE_PATTERN + `$`)
