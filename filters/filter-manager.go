package filters

import (
	"sync"
	"golang.org/x/net/html"
	)

type FilterManager struct {
	filters []SeoFilter
	tags []string
	parseGroup sync.WaitGroup
}

func NewManager() *FilterManager {
	fm := FilterManager{}
	fm.Configure()

	return &fm
}

func (fm *FilterManager) Add(filter SeoFilter) *FilterManager {
	fm.filters = append(fm.filters, filter)

	return fm
}

func (fm *FilterManager) Configure() {
	fm.Add(NewImgTag())
}

func (fm *FilterManager) Parse(node *html.Node) {
	if node != nil {
		for _, filter := range fm.filters {
			filter.Parse(node)
		}
	}
}

func (fm *FilterManager) clearTagsList() {
	fm.tags = []string{}
}

func (fm *FilterManager) createTagsList() {
	for _, tag := range fm.filters {
		fm.tags = append(fm.tags, tag.Tags()...)
	}
}
