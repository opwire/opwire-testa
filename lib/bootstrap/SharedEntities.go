package bootstrap

import (
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"github.com/opwire/opwire-testa/lib/engine"
	"github.com/opwire/opwire-testa/lib/format"
	"github.com/opwire/opwire-testa/lib/script"
	"github.com/opwire/opwire-testa/lib/tag"
	"github.com/opwire/opwire-testa/lib/utils"
)

func printUnmatchedPattern(outputPrinter *format.OutputPrinter, label string) string {
	if outputPrinter.IsColorized() {
		label = outputPrinter.NegativeTag(label)
	}
	return "(" + label + ")"
}

func printMarkedTags(outputPrinter *format.OutputPrinter, tags []string, mark map[string]int8) string {
	if len(tags) > 0 && len(mark) > 0 {
		noColor := !outputPrinter.IsColorized()
		tags = utils.Map(tags, func(tag string, pos int) string {
			if mark[tag] == -1 {
				if noColor {
					return "-" + tag
				}
				return outputPrinter.NegativeTag(tag)
			}
			if mark[tag] == +1 {
				if noColor {
					return "+" + tag
				}
				return outputPrinter.PositiveTag(tag)
			}
			return outputPrinter.RegularTag(tag)
		})
		return "(tags: " + strings.Join(tags, ", ") + ")"
	}
	return ""
}

func printDuration(outputPrinter *format.OutputPrinter, duration time.Duration) string {
	return "/ " + duration.String()
}

func printScriptSourceArgs(outputPrinter *format.OutputPrinter, scriptSource script.Source, scriptSelector *script.Selector, tagManager *tag.Manager) {
	var testDirs []string
	if scriptSource != nil {
		testDirs = scriptSource.GetTestDirs()
	}
	relaDirs := utils.DetectRelativePaths(testDirs)
	if relaDirs != nil && len(relaDirs) > 0 {
		outputPrinter.Println(outputPrinter.ContextInfo("Test directories", "", relaDirs...))
	} else {
		outputPrinter.Println(outputPrinter.ContextInfo("Test directories", "Unspecified"))
	}

	var inclFiles []string
	if scriptSource != nil {
		inclFiles = scriptSource.GetInclFiles()
	}
	if inclFiles != nil && len(inclFiles) > 0 {
		outputPrinter.Println(outputPrinter.ContextInfo("File inclusion patterns", "", inclFiles...))
	}

	var exclFiles []string
	if scriptSource != nil {
		exclFiles = scriptSource.GetExclFiles()
	}
	if exclFiles != nil && len(exclFiles) > 0 {
		outputPrinter.Println(outputPrinter.ContextInfo("File exclusion patterns", "", exclFiles...))
	}

	inclTags := tagManager.GetInclusiveTags()
	if inclTags != nil && len(inclTags) > 0 {
		outputPrinter.Println(outputPrinter.ContextInfo("Selected tags", strings.Join(inclTags, ", ")))
	}

	exclTags := tagManager.GetExclusiveTags()
	if exclTags != nil && len(exclTags) > 0 {
		outputPrinter.Println(outputPrinter.ContextInfo("Excluded tags", strings.Join(exclTags, ", ")))
	}

	testName := scriptSelector.GetTestNameFilter()
	if len(testName) > 0 {
		outputPrinter.Println(outputPrinter.ContextInfo("Name filter (" + scriptSelector.TypeOfTestNameFilter() + ")", testName))
	}
}

func filterInvalidDescriptors(src map[string]*script.Descriptor) (map[string]*script.Descriptor, []*script.Descriptor) {
	selected := make(map[string]*script.Descriptor, 0)
	rejected := make([]*script.Descriptor, 0)
	for key, d := range src {
		if d.Error == nil {
			selected[key] = d
		} else {
			rejected = append(rejected, d)
		}
	}
	return selected, rejected
}

func filterDescriptorsByInclusivePatterns(src map[string]*script.Descriptor, patterns []string) map[string]*script.Descriptor {
	if len(patterns) == 0 {
		return src
	}
	dst := make(map[string]*script.Descriptor, 0)
	for key, d := range src {
		for _, pattern := range patterns {
			if checkFilePathMatchPattern(d.Locator.AbsolutePath, pattern) {
				dst[key] = d
				continue
			}
		}
	}
	return dst
}

func filterDescriptorsByExclusivePatterns(src map[string]*script.Descriptor, patterns []string) map[string]*script.Descriptor {
	if len(patterns) == 0 {
		return src
	}
	dst := make(map[string]*script.Descriptor, 0)
	for key, d := range src {
		found := false
		for _, pattern := range patterns {
			if checkFilePathMatchPattern(d.Locator.AbsolutePath, pattern) {
				found = true
				break
			}
		}
		if !found {
			dst[key] = d
		}
	}
	return dst
}

func checkFilePathMatchPattern(fullPath string, pattern string) bool {
	srcPath, _ := utils.DetectRelativePath(fullPath)

	// try matching as a suffix string
	if strings.HasSuffix(srcPath, pattern) {
		return true
	}

	// try matching using filepath
	matched, err := filepath.Match(pattern, srcPath)
	if (err == nil && matched) {
		return true
	}

	// try matching with a regexp
	if reg, err := regexp.Compile(pattern); err == nil {
		if reg.MatchString(srcPath) {
			return true
		}
	}

	return false
}

func filterTestCasesByTags(tagManager *tag.Manager, testcases []*engine.TestCase) (accepted []*engine.TestCase, rejected []*engine.TestCase) {
	accepted = make([]*engine.TestCase, 0)
	rejected = make([]*engine.TestCase, 0)
	for _, testcase := range testcases {
		if active, _ := tagManager.IsActive(testcase.Tags); !active {
			rejected = append(rejected, testcase)
			continue
		}
		if testcase.Pending != nil && *testcase.Pending {
			rejected = append(rejected, testcase)
			continue
		}
		accepted = append(accepted, testcase)
	}
	return accepted, rejected
}
