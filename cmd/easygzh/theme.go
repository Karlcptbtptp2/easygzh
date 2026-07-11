package main

import (
	"fmt"

	"github.com/easygzh/easygzh/internal/cli"
	"github.com/spf13/cobra"
)

func newThemeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "theme",
		Short: "Manage built-in and user themes.",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List available built-in themes.",
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := themeManager().List()
			if err != nil {
				return emit(cmd, cli.Fail("THEME_LIST_FAILED", err.Error()))
			}
			if globalFlags.json {
				return emit(cmd, cli.OK("THEME_LIST", map[string]any{"themes": names}))
			}
			for _, n := range names {
				fmt.Fprintln(cmd.OutOrStdout(), n)
			}
			return nil
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "show <name>",
		Short: "Print the CSS of a built-in theme.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			css, err := themeManager().Load(args[0])
			if err != nil {
				return emit(cmd, cli.Fail("THEME_NOT_FOUND", err.Error()))
			}
			fmt.Fprint(cmd.OutOrStdout(), css)
			return nil
		},
	})
	return cmd
}
