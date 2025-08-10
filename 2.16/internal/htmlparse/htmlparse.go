package htmlparse

import (
	"bytes"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type Page struct {
	Data    []byte
	Links   []string // href/src links found
	Rewrite func(map[string]string) []byte
	IsHTML  bool
}

// Extract parses HTML, collects links (href, src) and returns original bytes.
// If content is not HTML, returns the data as-is and no links.
func Extract(u *url.URL, contentType string, body []byte) (Page, error) {
	isHTML := strings.Contains(strings.ToLower(contentType), "text/html")
	if !isHTML {
		return Page{Data: body, IsHTML: false}, nil
	}
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		// broken HTML, still save original
		return Page{Data: body, IsHTML: true}, nil
	}
	var links []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, attr := range n.Attr {
				switch attr.Key {
				case "href", "src":
					if attr.Val != "" {
						links = append(links, attr.Val)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	rewriter := func(repl map[string]string) []byte {
		// produce new HTML with replaced attributes from localMap
		var walk func(*html.Node)
		walk = func(n *html.Node) {
			if n.Type == html.ElementNode {
				for i := range n.Attr {
					if n.Attr[i].Key == "href" || n.Attr[i].Key == "src" {
						if newVal, ok := repl[n.Attr[i].Val]; ok {
							n.Attr[i].Val = newVal
						}
					}
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				walk(c)
			}
		}
		walk(doc)
		var buf bytes.Buffer
		_ = html.Render(&buf, doc)
		return buf.Bytes()
	}
	return Page{Data: body, Links: links, Rewrite: rewriter, IsHTML: true}, nil
}

// IsProbablyBinary helps decide to treat response as asset even if content-type missing
func IsProbablyBinary(ct string, u *url.URL) bool {
	if strings.Contains(ct, "text/html") {
		return false
	}
	lower := strings.ToLower(u.Path)
	for _, ext := range []string{".css", ".js", ".png", ".jpg", ".jpeg", ".gif", ".svg", ".ico", ".woff", ".woff2", ".ttf"} {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}
