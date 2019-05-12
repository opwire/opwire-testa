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

func (g *Manager) Parse(tagexps []string) {
	if g.exclusiveTags == nil {
		g.exclusiveTags = make([]string, 0)
	}
	if g.inclusiveTags == nil {
		g.inclusiveTags = make([]string, 0)
	}
	for _, tagexp := range tagexps {
		signedTags := utils.Split(tagexp, ",")
		for _, tag := range signedTags {
			if strings.HasPrefix(tag, "-") {
				g.exclusiveTags = appendTag(g.exclusiveTags, strings.TrimPrefix(tag, "-"))
			} else {
				g.inclusiveTags = appendTag(g.inclusiveTags, strings.TrimPrefix(tag, "+"))
			}
		}
	}
}

func (g *Manager) Reset() {
	g.inclusiveTags = nil
	g.exclusiveTags = nil
}

func (g *Manager) Initialize(tagexps []string) {
	g.Reset()
	g.Parse(tagexps)
}

func (g *Manager) GetInclusiveTags() []string {
	return g.inclusiveTags
}

func (g *Manager) GetExclusiveTags() []string {
	return g.exclusiveTags
}

func appendTag(tagStore []string, tags ...string) []string {
	for _, tag := range tags {
		if len(tag) > 0 && !utils.Contains(tagStore, tag) {
			tagStore = append(tagStore, tag)
		}
	}
	return tagStore
}
