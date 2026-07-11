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
