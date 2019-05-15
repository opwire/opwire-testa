package utils

import (
	"fmt"
	"regexp"
	"strings"
)

const VERSION_PATTERN string = `[v]?((\\d+\\.)?(\\d+\\.)?(\\*|\\d+))`

var versionRe = regexp.MustCompile(VERSION_PATTERN)

func StandardizeVersion(version string) string {
	return versionRe.ReplaceAllString(version, `$1`)
}

var tagRe = regexp.MustCompile(`^` + strings.ReplaceAll(TAG_PATTERN, `\\`, `\`) + `$`)

var tagCharRe = regexp.MustCompile(TAG_CHAR_PATTERN)

func StandardizeTagLabel(tag string) (string, error) {
	tag = tagCharRe.ReplaceAllString(tag, "_")
	if !tagRe.MatchString(tag) {
		return tag, fmt.Errorf("Tag label [%s] is invalid", tag)
	}
	return tag, nil
}

func ConvertTabToSpaces(block string, dedent int) string {
	lines := strings.Split(block, "\n")
	// determines the indent length
	indent := -1
	for _, line := range lines {
		tablen := strings.IndexFunc(line, func(c rune) bool {
			return c != '\t'
		})
		if tablen < 0 {
			continue
		}
		if indent < 0 || indent > tablen {
			indent = tablen
		}
	}
	if indent < 0 {
		indent = 0
	}
	// update dedent length
	if dedent == 0 || dedent > indent {
		dedent = indent
	}
	// de-indent the text block
	lines = Map(lines, func(line string, number int) string {
		var count int
		line = strings.TrimLeftFunc(line, func(r rune) bool {
			count += 1
			return count <= dedent
		})
		return strings.ReplaceAll(line, "\t", "  ")
	})
	// remove the first blank line
	if len(lines) > 0 && lines[0] == "" {
		lines = lines[1:]
	}
	return strings.Join(lines, "\n")
}

type DevNull int

func (DevNull) Write(p []byte) (int, error) {
	return len(p), nil
}

func (DevNull) WriteString(s string) (int, error) {
	return len(s), nil
}
