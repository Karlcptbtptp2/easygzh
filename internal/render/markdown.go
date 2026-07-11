package render

// Markdown rendering via goldmark. External-link-to-footnote conversion is done
// as an HTML post-processing step with goquery (see TransformLinks), rather than
// as a goldmark AST transformer — this is more robust to goldmark AST API churn
// and easier to test.
//
// Port of scripts/lib/link-footnote.mjs.

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// Ref is one collected external-link reference.
type Ref struct {
	N    int    // 1-based reference number
	Text string // original link text
	URL  string // original href
}

// MarkdownOptions configures the Markdown stage.
type MarkdownOptions struct {
	// LinkFootnotes converts external <a> links into text + superscript [n] and
	// appends a numbered reference list. Default true.
	LinkFootnotes bool
	// LinkFootnotesExplicit records whether the caller set LinkFootnotes, so a
	// zero-value MarkdownOptions{} can mean "default on". Set via the With* helpers.
	linkFootnotesSet bool
}

// WithLinkFootnotes returns opts with LinkFootnotes explicitly set.
func WithLinkFootnotes(on bool) MarkdownOptions {
	return MarkdownOptions{LinkFootnotes: on, linkFootnotesSet: true}
}

// RenderMarkdown parses markdown into HTML using goldmark (GFM: tables,
// strikethrough, linkify, task list). When link footnotes are enabled, external
// links are converted to text + superscript [n] with an appended reference list.
func RenderMarkdown(markdown string, opts MarkdownOptions) (string, []Ref, error) {
	footnotesOn := true
	if opts.linkFootnotesSet {
		footnotesOn = opts.LinkFootnotes
	}

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithHardWraps()),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(markdown), &buf); err != nil {
		return "", nil, fmt.Errorf("goldmark convert: %w", err)
	}
	htmlOut := buf.String()

	if !footnotesOn {
		return htmlOut, nil, nil
	}
	return TransformLinks(htmlOut)
}

// TransformLinks walks the HTML, converts external <a> links into the original
// text + a superscript [n], dedups by URL, and appends a References section.
// Internal anchors (#...) and non-web links are left untouched.
func TransformLinks(htmlIn string) (string, []Ref, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(htmlIn)))
	if err != nil {
		return "", nil, fmt.Errorf("parse html: %w", err)
	}

	var refs []Ref
	seen := map[string]int{} // url -> ref number

	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || !isExternal(href) {
			return
		}
		text := s.Text()

		var n int
		if existing, ok := seen[href]; ok {
			n = existing
		} else {
			n = len(refs) + 1
			seen[href] = n
			refs = append(refs, Ref{N: n, Text: text, URL: href})
		}

		// Stamp a data-ref and insert <sup>[n]</sup> right after the link.
		s.SetAttr("data-ref", fmt.Sprintf("%d", n))
		s.AfterHtml(fmt.Sprintf(`<sup style="font-size:0.75em;color:#888;">[%d]</sup>`, n))
	})

	if len(refs) == 0 {
		// Return body inner HTML (goquery wraps fragments in <html><body>).
		return bodyInner(doc), refs, nil
	}

	// Append: <hr><h3>References / 引用</h3><p>[n] text url</p>...
	body := doc.Find("body")
	body.AppendHtml(`<hr>`)
	body.AppendHtml(`<h3>References / 引用</h3>`)
	for _, r := range refs {
		body.AppendHtml(fmt.Sprintf(`<p>[%d] %s %s</p>`, r.N, escapeHTML(r.Text), r.URL))
	}

	return bodyInner(doc), refs, nil
}

func bodyInner(doc *goquery.Document) string {
	out, err := doc.Find("body").Html()
	if err != nil {
		return ""
	}
	return out
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

func isExternal(href string) bool {
	if href == "" || strings.HasPrefix(href, "#") {
		return false
	}
	if strings.HasPrefix(href, "mailto:") || strings.HasPrefix(href, "tel:") {
		return false
	}
	if u, err := url.Parse(href); err == nil && u.Scheme != "" {
		return u.Scheme == "http" || u.Scheme == "https"
	}
	return strings.HasPrefix(href, "//") ||
		strings.HasPrefix(href, "http://") ||
		strings.HasPrefix(href, "https://")
}
