// Package render implements the easyGZH rendering pipeline:
// Markdown -> HTML (goldmark) -> CSS inlined HTML (go-premailer) -> WeChat-safe HTML.
//
// This file wraps go-premailer for CSS inlining. The companion test
// (inline_test.go + testdata golden fixtures) verifies output parity with the
// Node juice baseline that the original easyGZH shipped.
package render

import (
	"fmt"

	premailer "github.com/vanng822/go-premailer/premailer"
)

// InlineOptions controls CSS-inlining behavior.
type InlineOptions struct {
	// PreserveImportant keeps !important declarations when true (default true).
	// go-premailer parses !important but its precedence handling is weaker than
	// juice; we keep this on to match the Node baseline as closely as possible.
	PreserveImportant bool
}

// InlineCSS takes an HTML string and a CSS string, returns the HTML with every
// matching CSS rule inlined onto the matched elements' style="..." attributes.
//
// The HTML should already be wrapped in the #easygzh-root scope (themes target
// `#easygzh-root h1` etc.), because go-premailer resolves selectors via Cascadia
// — descendant selectors work, but only if the ancestor exists in the input.
func InlineCSS(html, css string, opts InlineOptions) (string, error) {
	// go-premailer expects the CSS to live in a <style> block within the HTML it
	// parses. We splice it in at the top of the fragment. Premailer then moves
	// those rules inline and drops the <style> block.
	wrapped := fmt.Sprintf("<style>%s</style>%s", css, html)

	prem, err := premailer.NewPremailerFromString(wrapped, premailer.NewOptions())
	if err != nil {
		return "", fmt.Errorf("premailer init: %w", err)
	}

	out, err := prem.Transform()
	if err != nil {
		return "", fmt.Errorf("premailer transform: %w", err)
	}

	return out, nil
}
