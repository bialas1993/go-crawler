package filters

import "golang.org/x/net/html"

type HTag NodeFilter

func (h *HTag) Tags() []string {
	return []string{
		"h1", "h2", "h3",
	}
}

func (h *HTag) Filter(attributes []html.Attribute) {
	panic("implement me")
}

