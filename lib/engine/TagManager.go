package engine

import(
	"strings"
	"github.com/opwire/opwire-testa/lib/utils"
)

type TagManagerOptions interface {
	GetConditionalTags() []string
}

type TagManager struct {
	includedTags []string
	excludedTags []string
}

func NewTagManager(opts TagManagerOptions) (ref *TagManager, err error) {
	ref = &TagManager{}
	var conditionalTags []string
	if opts != nil {
		conditionalTags = opts.GetConditionalTags()
	}
	ref.Initialize(conditionalTags)
	return ref, err
}

func (g *TagManager) IsActive(tags []string) bool {
	if len(g.excludedTags) > 0 {
		for _, tag := range tags {
			if utils.Contains(g.excludedTags, tag) {
				return false
			}
		}
	}
	if len(g.includedTags) > 0 {
		for _, tag := range tags {
			if utils.Contains(g.includedTags, tag) {
				return true
			}
		}
		return false
	}
	return true
}

func (g *TagManager) Initialize(tagExps []string) {
	pTags := make([]string, 0)
	nTags := make([]string, 0)
	for _, exp := range tagExps {
		signedTags := utils.Split(exp, ",")
		for _, tag := range signedTags {
			if strings.HasPrefix(tag, "-") {
				nTags = append(nTags, strings.TrimPrefix(tag, "-"))
			} else {
				pTags = append(pTags, strings.TrimPrefix(tag, "+"))
			}
		}
	}
	g.includedTags = pTags
	g.excludedTags = nTags
}
