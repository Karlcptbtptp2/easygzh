// easygzh is the CLI entrypoint for the easyGZH WeChat-public-account formatting
// and publishing tool. It is the deterministic engine that an agent (driven by
// SKILL.md) orchestrates: rendering, image handling, WeChat publishing, diagnostics.
package main

import (
	"fmt"
	"os"

	"github.com/easygzh/easygzh/internal/theme"
	"github.com/spf13/cobra"
)

// Version is set at build time via -ldflags. Defaults to "dev".
var Version = "dev"

var globalFlags struct {
	json  bool
	quiet bool
	theme string
}

func main() {
	root := newRootCmd()
	if err := root.Execute(); err != nil {
		// cobra already prints the error; exit non-zero.
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:     "easygzh",
		Short:   "Format Markdown into WeChat-public-account HTML and publish drafts.",
		Long:    longDesc,
		Version: Version,
	}
	root.PersistentFlags().BoolVar(&globalFlags.json, "json", false, "emit structured JSON output (for agent consumption)")
	root.PersistentFlags().BoolVarP(&globalFlags.quiet, "quiet", "q", false, "suppress informational output")
	root.PersistentFlags().StringVar(&globalFlags.theme, "theme", "", "theme name (default: default)")

	root.AddCommand(newConvertCmd())
	root.AddCommand(newPreviewCmd())
	root.AddCommand(newVersionCmd())
	root.AddCommand(newDoctorCmd())
	root.AddCommand(newThemeCmd())
	root.AddCommand(newMemoryCmd())
	root.AddCommand(newInspectCmd())
	root.AddCommand(newPublishCmd())
	root.AddCommand(newImageCmd())
	root.AddCommand(newSkillsCmd())
	return root
}

var longDesc = `easygzh — format Markdown into WeChat-public-account HTML, with a stable
per-account tone, and publish to the draft box.

Pipeline: Markdown → goldmark → CSS inlined (go-premailer) → WeChat-safe HTML.

Run with --json for structured output consumable by an AI agent.`

// helper: build a theme manager rooted at the default themes dir.
func themeManager() *theme.Manager {
	return &theme.Manager{ThemesDir: theme.DefaultThemesDir()}
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the easygzh version.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Version)
		},
	}
}
