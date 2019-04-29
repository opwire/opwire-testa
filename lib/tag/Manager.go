package tag

import(
	"strings"
	"github.com/opwire/opwire-testa/lib/utils"
)

type ManagerOptions interface {
	GetConditionalTags() []string
}

type Manager struct {
	includedTags []string
	excludedTags []string
}

func NewManager(opts ManagerOptions) (ref *Manager, err error) {
	ref = &Manager{}
	var conditionalTags []string
	if opts != nil {
		conditionalTags = opts.GetConditionalTags()
	}
	ref.Initialize(conditionalTags)
	return ref, err
}

func (g *Manager) IsActive(tags []string) bool {
	if len(tags) == 0 {
		return true
	}
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

func (g *Manager) Initialize(tagexps []string) {
	pTags := make([]string, 0)
	nTags := make([]string, 0)
	for _, tagexp := range tagexps {
		signedTags := utils.Split(tagexp, ",")
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

func (g *Manager) GetIncludedTags() []string {
	return g.includedTags
}

func (g *Manager) GetExcludedTags() []string {
	return g.excludedTags
}
