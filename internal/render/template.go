package render

// Structure templates add a content-structure layer on top of visual themes.
// While a theme (CSS) controls how Markdown tags look (h1, p, blockquote...),
// a structure template controls the *narrative layout*: brand header, hook,
// body, CTA, footer. It wraps the goldmark-rendered HTML in richer structural
// HTML that Markdown alone cannot express.
//
// Template files (.html) contain:
//   - <!-- template: name / description --> frontmatter in HTML comments
//   - <style data-template-css>...</style> component CSS (inlined with theme CSS)
//   - <!-- slot: xxx --> markers where the engine injects content
//
// Slot injection is a goquery HTML post-process, mirroring the TransformLinks
// pattern (see markdown.go). Template HTML goes through the same CSS inlining
// and WeChat safety pipeline downstream.

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// StructureTemplate is a parsed content-structure template.
type StructureTemplate struct {
	Name        string // template name (e.g. "mindful-journal")
	Description string // human-readable description
	HTML        string // template body HTML (style block extracted, slots intact)
	CSS         string // component CSS extracted from <style data-template-css>
}

// SlotData provides values for placeholder substitution.
type SlotData struct {
	Title       string // document title (from first h1 or --title flag)
	BrandLabel  string // brand label (e.g. "( 每月正念 )")
	BrandDate   string // date string (auto-filled if empty)
	Subtitle    string // subtitle/lead text
	CTAText     string // call-to-action text
	BrandFooter string // footer brand name
}

// slotCommentRe matches <!-- slot: body --> style markers.
var slotCommentRe = regexp.MustCompile(`<!--\s*slot:\s*(\w+)\s*-->`)

// templateMetaRe extracts name/description from <!-- template: ... --> comments.
// The comment may span multiple lines; we match the opening tag and capture
// everything until the closing -->.
var templateMetaRe = regexp.MustCompile(`(?s)<!--\s*template:\s*(.*?)\s*-->`)

// ApplyStructureTemplate wraps goldmark-rendered HTML in a structure template.
// If tmpl is nil or its HTML is empty, htmlContent is returned unchanged
// (backward compatibility — no template means linear Markdown rendering).
func ApplyStructureTemplate(htmlContent string, tmpl *StructureTemplate, data SlotData) (string, error) {
	if tmpl == nil || tmpl.HTML == "" {
		return htmlContent, nil
	}

	// Fill placeholder values with defaults where empty.
	if data.BrandDate == "" {
		data.BrandDate = time.Now().Format("2006年1月")
	}

	// Extract title from the goldmark HTML if not provided: first <h1> text.
	if data.Title == "" {
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(htmlContent)))
		if err == nil {
			data.Title = doc.Find("h1").First().Text()
		}
	}

	// Resolve slots: replace <!-- slot: body --> with the rendered HTML,
	// and fill other slot regions with template-defined content.
	result := injectBodySlot(tmpl.HTML, htmlContent)

	// Replace placeholders {{TITLE}}, {{BRAND_LABEL}}, etc.
	result = fillPlaceholders(result, data)

	return result, nil
}

// injectBodySlot replaces the <!-- slot: body --> marker with the rendered
// HTML content. All other slots (brand, hook, cta...) remain as-is — they
// are part of the template's static HTML.
func injectBodySlot(templateHTML, bodyHTML string) string {
	return slotCommentRe.ReplaceAllStringFunc(templateHTML, func(match string) string {
		subs := slotCommentRe.FindStringSubmatch(match)
		if len(subs) < 2 {
			return match
		}
		slotName := subs[1]
		if slotName == "body" {
			return bodyHTML
		}
		// Non-body slots: keep the comment marker (it's a no-op in HTML).
		return match
	})
}

// fillPlaceholders replaces {{KEY}} tokens with SlotData values.
func fillPlaceholders(html string, data SlotData) string {
	replacements := map[string]string{
		"{{TITLE}}":        data.Title,
		"{{BRAND_LABEL}}":  data.BrandLabel,
		"{{BRAND_DATE}}":   data.BrandDate,
		"{{SUBTITLE}}":     data.Subtitle,
		"{{CTA_TEXT}}":     data.CTAText,
		"{{BRAND_FOOTER}}": data.BrandFooter,
	}
	result := html
	for key, val := range replacements {
		result = strings.ReplaceAll(result, key, val)
	}
	return result
}

// ParseTemplate parses a raw .html template file into a StructureTemplate,
// extracting the frontmatter (name/description), the component CSS
// (<style data-template-css>), and the body HTML (with style block removed).
func ParseTemplate(raw string) (*StructureTemplate, error) {
	tmpl := &StructureTemplate{}

	// 1. Extract frontmatter from <!-- template: ... --> comment.
	if m := templateMetaRe.FindStringSubmatch(raw); len(m) >= 2 {
		parseTemplateFrontmatter(m[1], tmpl)
	}

	// 2. Extract CSS from <style data-template-css>...</style>.
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(raw)))
	if err != nil {
		return nil, fmt.Errorf("parse template html: %w", err)
	}

	styleSel := doc.Find(`style[data-template-css]`)
	if styleSel.Length() > 0 {
		tmpl.CSS, _ = styleSel.First().Html()
		styleSel.Remove()
	}

	// 3. Get the body HTML (goquery wraps fragments in <html><body>).
	tmpl.HTML = bodyInner(doc)

	// Default name if not set.
	if tmpl.Name == "" {
		tmpl.Name = "unnamed"
	}

	return tmpl, nil
}

// parseTemplateFrontmatter fills name/description from key:value lines.
func parseTemplateFrontmatter(meta string, tmpl *StructureTemplate) {
	for _, line := range strings.Split(meta, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Split on first ":" — values may contain colons.
		idx := strings.Index(line, ":")
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		switch key {
		case "name":
			tmpl.Name = val
		case "description":
			tmpl.Description = val
		}
	}
}
