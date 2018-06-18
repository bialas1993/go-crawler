package filters

type ImgTag struct {
	SeoFilter
}

func (i ImgTag) Tag() []string {
	return []string{"img"}
}

func (i ImgTag) Filter() {
	panic("implement me")
}



