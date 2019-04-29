package bootstrap

import (
	"path/filepath"
	"strings"
	"github.com/opwire/opwire-testa/lib/format"
	"github.com/opwire/opwire-testa/lib/script"
	"github.com/opwire/opwire-testa/lib/tag"
	"github.com/opwire/opwire-testa/lib/utils"
)

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
		return "(" + strings.Join(tags, ", ") + ")"
	}
	return ""
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

	inclTags := tagManager.GetIncludedTags()
	if inclTags != nil && len(inclTags) > 0 {
		outputPrinter.Println(outputPrinter.ContextInfo("Selected tags", strings.Join(inclTags, ", ")))
	}

	exclTags := tagManager.GetExcludedTags()
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

func filterDescriptorsByFilePattern(src map[string]*script.Descriptor, suffix string) map[string]*script.Descriptor {
	if len(suffix) == 0 {
		return src
	}
	dst := make(map[string]*script.Descriptor, 0)
	for key, d := range src {
		if strings.HasSuffix(d.Locator.AbsolutePath, suffix) {
			dst[key] = d
			continue
		}
		matched, err := filepath.Match(suffix, d.Locator.AbsolutePath)
		if (err == nil && matched) {
			dst[key] = d
			continue
		}
	}
	return dst
}
