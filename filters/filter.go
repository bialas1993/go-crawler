package filters

import (
	"golang.org/x/net/html"
	)

type SeoFilter interface {
	Tags() []string
	Filter(attributes []html.Attribute)
	Parse(node *html.Node)
}

type NodeFilter struct {
	SeoFilter
}

func (nf *NodeFilter) Filter(attributes []html.Attribute) {
	panic("Not implemented yet.")
}

func (nf *NodeFilter) Tags() []string {
	panic("Not implemented yet.")
}

func (nf *NodeFilter) Parse(node *html.Node) {
	nf.parse(nf.Tags(), node)
}

func (nf *NodeFilter) parse(tags []string, node *html.Node) {
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && supportTag(nf, n.Data) {
			nf.Filter(n.Attr)
		}
	}

	forEachNode(node, visitNode, nil)
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

func supportTag(sf* NodeFilter, nodeName string) bool {
	for _, tag := range sf.Tags(){
		if nodeName == tag {
			return true
		}
	}

	return false
}




