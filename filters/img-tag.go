package filters

func NewImgTag() *ImgTag{
	return new(ImgTag)
}

type ImgTag struct {
	NodeFilter
}

func (i *ImgTag) Tags() []string {
	return []string{"img"}
}
