package render

// WeChat HTML post-processing: makes the inlined HTML safe to paste into the
// WeChat public account editor. Port of scripts/lib/wechat-postprocess.mjs.
//
// The WeChat editor strips <style>/<script>, drops unknown tags, resets list
// markers, caps same-tag nesting at <15, and needs absolute image URLs.

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

var allowedTags = map[string]bool{
	"p": true, "span": true, "div": true, "section": true, "img": true, "a": true,
	"ul": true, "ol": true, "li": true, "strong": true, "em": true, "b": true, "i": true,
	"br": true, "hr": true,
	"table": true, "thead": true, "tbody": true, "tr": true, "td": true, "th": true,
	"pre": true, "code": true, "blockquote": true,
	"h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true,
	"sup": true, "sub": true, "mark": true,
	// Structural tags produced by goquery's full-document wrapping. They are
	// never part of the final WeChat fragment (we serialize only <body>'s
	// children), so keeping them here just prevents unwrapUnknownTags from
	// dismantling the document during post-processing.
	"html": true, "head": true, "body": true, "title": true, "meta": true,
}

const nestingDepthLimit = 15

// WeChatOptions controls the post-processing step.
type WeChatOptions struct {
	// WrapSection, when true (default), wraps the whole fragment in a single
	// <section id="easygzh-root">. Themes scope their selectors under this id.
	WrapSection bool
	wrapSet     bool
}

// WeChatSafe makes inlined HTML safe for the WeChat editor and returns the
// fragment (inner HTML, no <html>/<body> wrappers).
func WeChatSafe(inlinedHTML string, opts WeChatOptions) (string, error) {
	wrap := true
	if opts.wrapSet {
		wrap = opts.WrapSection
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(inlinedHTML)))
	if err != nil {
		return "", fmt.Errorf("parse html: %w", err)
	}

	stripDangerous(doc)
	if err := unwrapUnknownTags(doc); err != nil {
		return "", err
	}
	enforceNestingDepth(doc)
	absolutizeImages(doc)
	hardenLists(doc)

	// Serialize via html.Render. Find the <body> element by walking the parsed
	// tree (goquery's NewDocumentFromReader always wraps a fragment in a full
	// document, so a <body> exists somewhere; we just locate it).
	bodyNode := findBody(doc.Nodes[0])
	if bodyNode == nil {
		// No body found — serialize the whole thing.
		bodyNode = doc.Nodes[0]
	}
	var b bytes.Buffer
	for c := bodyNode.FirstChild; c != nil; c = c.NextSibling {
		if err := html.Render(&b, c); err != nil {
			return "", fmt.Errorf("serialize: %w", err)
		}
	}
	inner := b.String()
	if !wrap {
		return inner, nil
	}
	return fmt.Sprintf(`<section id="easygzh-root">%s</section>`, inner), nil
}

func stripDangerous(doc *goquery.Document) {
	doc.Find("script, style, link, iframe, object, embed, noscript").Remove()
}

// unwrapUnknownTags replaces disallowed tags with their children. Bounded passes
// (port of the Node fix that avoided the unbounded fixpoint loop).
func unwrapUnknownTags(doc *goquery.Document) error {
	for pass := 0; pass < 8; pass++ {
		var targets []*html.Node
		doc.Find("*").Each(func(_ int, s *goquery.Selection) {
			if n := s.Nodes[0]; n.Type == html.ElementNode && !allowedTags[n.Data] {
				targets = append(targets, n)
			}
		})
		if len(targets) == 0 {
			break
		}
		for _, n := range targets {
			// Replace the element with its children.
			parent := n.Parent
			if parent == nil {
				continue
			}
			// Collect children first (the list mutates as we move them).
			var children []*html.Node
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				children = append(children, c)
			}
			for _, c := range children {
				n.RemoveChild(c)
				parent.InsertBefore(c, n)
			}
			parent.RemoveChild(n)
		}
	}
	return nil
}

func enforceNestingDepth(doc *goquery.Document) {
	doc.Find("*").Each(func(_ int, s *goquery.Selection) {
		n := s.Nodes[0]
		if n.Type != html.ElementNode {
			return
		}
		tag := n.Data
		depth := 0
		for a := n.Parent; a != nil && a.Type == html.ElementNode; a = a.Parent {
			if a.Data == tag {
				depth++
				if depth >= nestingDepthLimit {
					// Flatten: replace n with its children.
					parent := n.Parent
					var children []*html.Node
					for c := n.FirstChild; c != nil; c = c.NextSibling {
						children = append(children, c)
					}
					for _, c := range children {
						n.RemoveChild(c)
						parent.InsertBefore(c, n)
					}
					parent.RemoveChild(n)
					return
				}
			}
		}
	})
}

func absolutizeImages(doc *goquery.Document) {
	doc.Find("img").Each(func(_ int, s *goquery.Selection) {
		src, ok := s.Attr("src")
		src = strings.TrimSpace(src)
		if !ok || src == "" {
			s.Remove()
			return
		}
		if strings.HasPrefix(src, "//") {
			src = "https:" + src
		}
		if !strings.HasPrefix(src, "http://") && !strings.HasPrefix(src, "https://") && !strings.HasPrefix(src, "data:") {
			s.SetAttr("data-bad-src", src)
			s.SetAttr("src", "")
		} else {
			s.SetAttr("src", src)
		}
		s.SetAttr("referrerpolicy", "no-referrer")
	})
}

func hardenLists(doc *goquery.Document) {
	doc.Find("ol > li").Each(func(_ int, s *goquery.Selection) {
		appendStyle(s, "list-style: decimal inside;")
	})
	doc.Find("ul > li").Each(func(_ int, s *goquery.Selection) {
		appendStyle(s, "list-style: disc inside;")
	})
}

func appendStyle(s *goquery.Selection, decl string) {
	existing, _ := s.Attr("style")
	if existing != "" && !strings.HasSuffix(strings.TrimSpace(existing), ";") {
		existing = strings.TrimSpace(existing) + ";"
	}
	combined := strings.TrimSpace(existing + " " + decl)
	s.SetAttr("style", combined)
}

// findBody walks the node tree and returns the first <body> element, or nil.
func findBody(n *html.Node) *html.Node {
	if n == nil {
		return nil
	}
	if n.Type == html.ElementNode && n.Data == "body" {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if found := findBody(c); found != nil {
			return found
		}
	}
	return nil
}
