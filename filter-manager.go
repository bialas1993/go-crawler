package main

import "github.com/bialas1993/go-crawler/filters"

type FilterManager struct {
	filters []filters.SeoFilter
	tags []string
}

func (fm *FilterManager) Add(filter filters.SeoFilter) *FilterManager {
	fm.filters = append(fm.filters, filter)

	return fm
}

func (fm *FilterManager) Configure() {
	fm.createTagsList()
}

func (fm *FilterManager) clearTagsList() {
	fm.tags = []string{}
}

func (fm *FilterManager) createTagsList() {
	for _, tag := range fm.filters {
		fm.tags = append(fm.tags, tag.Tag()...)
	}
}
