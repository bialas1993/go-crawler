package filters

import (
	"golang.org/x/net/html"
	)

type SeoFilter interface {
	Tags() []string
	Filter(attributes []html.Attribute)
	Parse(s SeoFilter, node *html.Node)
}

type NodeFilter struct {
	SeoFilter
}

func (nf *NodeFilter) Filter(attributes []html.Attribute) {
	hasAlt := false

	for _, a := range attributes {
		if a.Key == "alt" && a.Val != "" {
			hasAlt = true
			break
		}
	}

	if !hasAlt {
		// log
	}
}


func (nf *NodeFilter) Parse(s SeoFilter, node *html.Node) {
	nf.parse(s, node)
}

func (nf *NodeFilter) parse(s SeoFilter, node *html.Node) {
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && supportTag(s, n.Data) {
			s.Filter(n.Attr)
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

func supportTag(sf SeoFilter, nodeName string) bool {
	for _, tag := range sf.Tags(){
		if nodeName == tag {
			return true
		}
	}

	return false
}




