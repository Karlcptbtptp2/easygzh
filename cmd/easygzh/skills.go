package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/easygzh/easygzh/internal/cli"
	"github.com/spf13/cobra"
)

// skills outputs the bundled SKILL.md so an agent can read it without searching
// the repo. Mirrors md2wechat's `skills read` command.
func newSkillsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skills",
		Short: "Read the bundled SKILL.md (the agent operating manual).",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "read easygzh",
		Short: "Print the SKILL.md content.",
		RunE: func(cmd *cobra.Command, args []string) error {
			content, err := readSkillMD()
			if err != nil {
				return emit(cmd, cli.Fail("SKILL_NOT_FOUND", err.Error()))
			}
			fmt.Fprint(cmd.OutOrStdout(), content)
			return nil
		},
	})
	return cmd
}

func readSkillMD() (string, error) {
	candidates := []string{"SKILL.md", "../SKILL.md", "../../SKILL.md"}
	for _, c := range candidates {
		if b, err := os.ReadFile(c); err == nil {
			return string(b), nil
		}
	}
	if exe, err := os.Executable(); err == nil {
		for d := filepath.Dir(exe); d != "/" && d != "."; d = filepath.Dir(d) {
			if b, err := os.ReadFile(filepath.Join(d, "SKILL.md")); err == nil {
				return string(b), nil
			}
		}
	}
	return "", fmt.Errorf("SKILL.md not found")
}
