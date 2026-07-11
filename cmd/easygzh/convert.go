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
		outFile     string
		noFootnotes bool
		themeCSS    string // raw CSS override
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

			out, err := render.Render(string(md), render.PipelineOptions{
				ThemeCSS:      css,
				LinkFootnotes: !noFootnotes,
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
