package utils

import (
	"regexp"
	"strings"
)

const BLANK string = ""

const BODY_FORMAT_FLAT = `text`
const BODY_FORMAT_JSON = `json`
const BODY_FORMAT_YAML = `yaml`

const DEFAULT_PDP string = `http://localhost:17779`
const DEFAULT_PATH string = `/-`

const TAG_CHAR_PATTERN string = `[^a-zA-Z0-9_-]`
const TAG_PATTERN string = `[a-zA-Z][a-zA-Z0-9]*([_-][a-zA-Z0-9]*)*`
const TIME_RFC3339 string = `([0-9]+)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|60)(\\.[0-9]+)?(([Zz])|([\\+|\\-]([01][0-9]|2[0-3]):[0-5][0-9]))`
const TIMEOUT_PATTERN string = `([0-9]+h)?([0-9]+m)?([0-9]+s)?([0-9]+ms)?([0-9]+[uµ]s)?([0-9]+ns)?`

const TEST_CASE_TITLE_PATTERN string = `[a-zA-Z][\\w\\-\\s.:;,]*`
var TEST_CASE_TITLE_REGEXP *regexp.Regexp = regexp.MustCompile(`^` + strings.ReplaceAll(TEST_CASE_TITLE_PATTERN, `\\`, `\`) + `$`)
