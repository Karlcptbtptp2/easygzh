package main

import (
	"github.com/easygzh/easygzh/internal/render"
)

// renderPipeline runs the standard Markdown → inline-styled HTML conversion with
// footnotes on. Used by publish and preview helpers.
func renderPipeline(markdown, themeCSS string) (string, error) {
	return render.Render(markdown, render.PipelineOptions{
		ThemeCSS:      themeCSS,
		LinkFootnotes: true,
	})
}

// renderPipelineWithTemplate is like renderPipeline but applies a structure
// template. Used when globalFlags or publish options specify a template.
func renderPipelineWithTemplate(markdown, themeCSS string, tmplHTML, tmplCSS string, slotData render.SlotData) (string, error) {
	return render.Render(markdown, render.PipelineOptions{
		ThemeCSS:      themeCSS,
		LinkFootnotes: true,
		TemplateHTML:  tmplHTML,
		TemplateCSS:   tmplCSS,
		SlotData:      slotData,
	})
}
