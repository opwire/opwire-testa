package tag

import(
	"strings"
	"github.com/opwire/opwire-testa/lib/utils"
)

type ManagerOptions interface {
	GetConditionalTags() []string
}

type Manager struct {
	inclusiveTags []string
	exclusiveTags []string
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

func (g *Manager) IsActive(tags []string) (bool, map[string]int8) {
	mark := make(map[string]int8, 0)
	if len(tags) == 0 {
		return true, mark
	}
	if len(g.exclusiveTags) > 0 {
		for _, tag := range tags {
			if utils.Contains(g.exclusiveTags, tag) {
				mark[tag] = -1
				return false, mark
			}
		}
	}
	if len(g.inclusiveTags) > 0 {
		for _, tag := range tags {
			if utils.Contains(g.inclusiveTags, tag) {
				mark[tag] = +1
				return true, mark
			}
		}
		return false, mark
	}
	return true, mark
}

func (g *Manager) Initialize(tagexps []string) {
	pTags := make([]string, 0)
	nTags := make([]string, 0)
	for _, tagexp := range tagexps {
		signedTags := utils.Split(tagexp, ",")
		for _, tag := range signedTags {
			if strings.HasPrefix(tag, "-") {
				nTag := strings.TrimPrefix(tag, "-")
				if len(nTag) > 0 && !utils.Contains(nTags, nTag) {
					nTags = append(nTags, nTag)
				}
			} else {
				pTag := strings.TrimPrefix(tag, "+")
				if len(pTag) > 0 && !utils.Contains(pTags, pTag) {
					pTags = append(pTags, pTag)
				}
			}
		}
	}
	g.inclusiveTags = pTags
	g.exclusiveTags = nTags
}

func (g *Manager) GetInclusiveTags() []string {
	return g.inclusiveTags
}

func (g *Manager) GetExclusiveTags() []string {
	return g.exclusiveTags
}
