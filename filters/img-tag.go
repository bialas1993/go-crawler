package filters

import (
	"golang.org/x/net/html"
	)

func NewImgTag() *ImgTag{
	return new(ImgTag)
}

type ImgTag struct {
	NodeFilter
}

func (i *ImgTag) Tags() []string {
	return []string{"img"}
}

func (i *ImgTag) Filter(attributes []html.Attribute) {
	panic("implement me")
}

func (i *ImgTag) Parse(node *html.Node) {
	i.parse(i.Tags(), node)
}