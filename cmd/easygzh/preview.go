package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/easygzh/easygzh/internal/cli"
	"github.com/spf13/cobra"
)

func newPreviewCmd() *cobra.Command {
	var theme string
	cmd := &cobra.Command{
		Use:   "preview <markdown-file>",
		Short: "Render Markdown to a temporary HTML file and open it in the browser.",
		Long: `Renders the article with the given theme and opens it in your default
browser so you can eyeball the result before publishing. No uploads, no drafts.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if theme != "" {
				globalFlags.theme = theme
			}
			tmp, err := os.CreateTemp("", "easygzh-preview-*.html")
			if err != nil {
				return emit(cmd, cli.Fail("PREVIEW_FAILED", err.Error()))
			}
			tmpPath := tmp.Name()
			tmp.Close()
			// Defer cleanup? Keep the file so the browser still has it.

			// Reuse convert logic by invoking it with --output.
			convert := newConvertCmd()
			convert.SetArgs([]string{args[0], "--output", tmpPath, "--theme", globalFlags.theme})
			convert.SetOut(os.Stderr)
			convert.SetErr(os.Stderr)
			if err := convert.Execute(); err != nil {
				return err
			}

			if openErr := openBrowser(tmpPath); openErr != nil {
				fmt.Fprintf(os.Stderr, "wrote %s (open it manually)\n", tmpPath)
			}
			return emit(cmd, cli.OK("PREVIEW_READY", map[string]any{
				"file": tmpPath,
			}, "review in browser", "publish when satisfied"))
		},
	}
	cmd.Flags().StringVar(&theme, "theme", "", "theme name to preview")
	return cmd
}

func openBrowser(path string) error {
	abs, _ := filepath.Abs(path)
	url := "file://" + abs
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default:
		return exec.Command("xdg-open", url).Start()
	}
}
