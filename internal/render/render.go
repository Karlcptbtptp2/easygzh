// Package render implements the easyGZH rendering pipeline.
//
// Pipeline (pure, deterministic, no network):
//
//	Markdown ──goldmark──▶ HTML
//	           theme CSS ──┐
//	                       ▼
//	               go-premailer (CSS inline)
//	                       ▼
//	               WeChat-safe post-process
//	                       ▼
//	               inline-styled HTML fragment
//
// Determinism is the foundation of "stable tone": the same input always yields
// the same output, so per-account visual identity never drifts.
package render

import "fmt"

// PipelineOptions configures a full render.
type PipelineOptions struct {
	// ThemeCSS is the concatenated CSS to inline (base theme + profile overrides).
	ThemeCSS string
	// LinkFootnotes converts external links to references. Default true.
	LinkFootnotes bool
}

// Render runs the full pipeline and returns WeChat-safe inline-styled HTML.
func Render(markdown string, opts PipelineOptions) (string, error) {
	if opts.ThemeCSS == "" {
		return "", fmt.Errorf("render: ThemeCSS is required")
	}

	htmlOut, _, err := RenderMarkdown(markdown, MarkdownOptions{
		LinkFootnotes: opts.LinkFootnotes,
	})
	if err != nil {
		return "", err
	}

	// Scope the fragment BEFORE inlining, so `#easygzh-root h1` selectors match.
	scoped := fmt.Sprintf(`<section id="easygzh-root">%s</section>`, htmlOut)

	inlined, err := InlineCSS(scoped, opts.ThemeCSS, InlineOptions{PreserveImportant: true})
	if err != nil {
		return "", err
	}

	// Post-process for WeChat safety WITHOUT re-wrapping (already wrapped above).
	safe, err := WeChatSafe(inlined, WeChatOptions{WrapSection: false, wrapSet: true})
	if err != nil {
		return "", err
	}
	return safe, nil
}
