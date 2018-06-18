package filters

type ImgTag struct {
	SeoFilter
}

func (*ImgTag) Tag() []string {
	return []string{"img"}
}

func (*ImgTag) Filter() {
	panic("implement me")
}



