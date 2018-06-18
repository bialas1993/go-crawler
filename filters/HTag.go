package filters

type HTag struct{
	SeoFilter
}

func (h *HTag) Tag() []string {
	return []string{
		"h1", "h2", "h3",
	}
}

func (h *HTag) Filter() {
	panic("implement me")
}

