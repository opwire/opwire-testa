package bootstrap

import (
	"path/filepath"
	"strings"
	"github.com/opwire/opwire-testa/lib/script"
)

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
