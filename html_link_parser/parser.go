package html_link_parser

import (
	"golang.org/x/net/html"
	"io"
	"strings"
)

type Link struct {
	Href string
	Text string
}

func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	linkNodes := extractLinkNodes(doc)

	var links []Link
	for _, node := range linkNodes {
		links = append(links, newLink(node))
	}

	return links, nil
}

func extractLinkNodes(tree *html.Node) []*html.Node {
	// if we found the "a" tag then return the node instantly
	if tree.Type == html.ElementNode && tree.Data == "a" {
		return []*html.Node{tree}
	}

	// loop over all the children of current node to figure out
	// if any of them may be an "a" tag
	var linkNodes []*html.Node
	for child := tree.FirstChild; child != nil; child = child.NextSibling {
		linkNodes = append(linkNodes, extractLinkNodes(child)...)
	}

	return linkNodes
}

func newLink(node *html.Node) Link {
	var link Link

	// extract the value of href from the given node if
	// it exists
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			link.Href = attr.Val
			break
		}
	}

	// extract the text from all children recursively
	link.Text = extractLinkText(node)

	return link
}

func extractLinkText(node *html.Node) string {
	// if the node is a text node we found the target
	if node.Type == html.TextNode {
		return strings.TrimSpace(node.Data)
	}

	// if the node is anything other than an html element
	// then we can ignore it (eg. Comment, Doctype etc.)
	if node.Type != html.ElementNode {
		return ""
	}

	// iterate over all the children of current element node and
	// extract text from all nodes
	var text string
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		text += extractLinkText(child)
	}

	return text
}
