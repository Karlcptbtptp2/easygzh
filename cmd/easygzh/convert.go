package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/easygzh/easygzh/internal/cli"
	"github.com/easygzh/easygzh/internal/render"
	"github.com/spf13/cobra"
)

func newConvertCmd() *cobra.Command {
	var (
		outFile      string
		noFootnotes  bool
		themeCSS     string // raw CSS override
		templateName string // structure template name
		title        string // override title for template slots
		brandLabel   string // brand label for template
		brandFooter  string // footer brand name for template
		subtitle     string // subtitle/lead for template
		ctaText      string // CTA text for template
	)
	cmd := &cobra.Command{
		Use:   "convert <markdown-file>",
		Short: "Convert Markdown to WeChat-safe inline-styled HTML.",
		Long: `Convert a Markdown file into WeChat-public-account HTML with all theme CSS
inlined as style="..." attributes. No network, no side effects. Output goes to
stdout unless --output is given.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			md, err := os.ReadFile(args[0])
			if err != nil {
				return emit(cmd, cli.Fail("READ_FAILED", err.Error()))
			}

			// Resolve theme CSS.
			var css string
			if themeCSS != "" {
				css = themeCSS
			} else {
				name := globalFlags.theme
				if name == "" {
					name = "default"
				}
				css, err = themeManager().Load(name)
				if err != nil {
					return emit(cmd, cli.Fail("THEME_NOT_FOUND", err.Error()))
				}
			}

			// Load structure template if specified.
			var tmplCSS, tmplHTML string
			var slotData render.SlotData
			if templateName != "" {
				tmpl, err := templateManager().Load(templateName)
				if err != nil {
					return emit(cmd, cli.Fail("TEMPLATE_NOT_FOUND", err.Error()))
				}
				tmplCSS = tmpl.CSS
				tmplHTML = tmpl.HTML
				slotData = render.SlotData{
					Title:       title,
					BrandLabel:  brandLabel,
					Subtitle:    subtitle,
					CTAText:     ctaText,
					BrandFooter: brandFooter,
				}
			}

			out, err := render.Render(string(md), render.PipelineOptions{
				ThemeCSS:      css,
				LinkFootnotes: !noFootnotes,
				TemplateHTML:  tmplHTML,
				TemplateCSS:   tmplCSS,
				SlotData:      slotData,
			})
			if err != nil {
				return emit(cmd, cli.Fail("CONVERT_FAILED", err.Error()))
			}

			if outFile != "" {
				if werr := os.WriteFile(outFile, []byte(out), 0o644); werr != nil {
					return emit(cmd, cli.Fail("WRITE_FAILED", werr.Error()))
				}
				return emit(cmd, cli.OK("CONVERT_COMPLETED", map[string]any{
					"output_file": outFile,
					"bytes":       len(out),
				}, "preview the file in a browser", "publish with: easygzh publish"))
			}

			if globalFlags.json {
				return emit(cmd, cli.OK("CONVERT_COMPLETED", map[string]any{
					"html": out,
				}))
			}
			fmt.Print(out)
			return nil
		},
	}
	cmd.Flags().StringVarP(&outFile, "output", "o", "", "write HTML to this file instead of stdout")
	cmd.Flags().BoolVar(&noFootnotes, "no-footnotes", false, "keep external links as-is instead of converting to references")
	cmd.Flags().StringVar(&themeCSS, "css", "", "raw CSS override (bypass theme loading)")
	cmd.Flags().StringVarP(&templateName, "template", "t", "", "structure template name (e.g. mindful-journal, book-club, product-launch)")
	cmd.Flags().StringVar(&title, "title", "", "title for template header (auto-extracted from h1 if omitted)")
	cmd.Flags().StringVar(&brandLabel, "brand-label", "", "brand label for template header")
	cmd.Flags().StringVar(&brandFooter, "brand-footer", "", "footer brand name")
	cmd.Flags().StringVar(&subtitle, "subtitle", "", "subtitle/lead text for template header")
	cmd.Flags().StringVar(&ctaText, "cta-text", "", "call-to-action text for template")
	return cmd
}

// emit prints a cli.Response (JSON when --json, else error-to-stderr on failure).
func emit(cmd *cobra.Command, resp cli.Response) error {
	if globalFlags.json {
		b, _ := json.Marshal(resp)
		fmt.Fprintln(cmd.OutOrStdout(), string(b))
		if !resp.Success {
			os.Exit(1)
		}
		return nil
	}
	if !resp.Success {
		fmt.Fprintln(os.Stderr, "error: "+resp.Error)
		os.Exit(1)
	}
	return nil
}
