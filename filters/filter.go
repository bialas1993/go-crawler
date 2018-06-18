package filters

type SeoFilter interface {
	Tag() []string
	Filter()
}
