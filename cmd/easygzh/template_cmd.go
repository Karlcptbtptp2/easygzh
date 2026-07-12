package main

import (
	"fmt"

	"github.com/easygzh/easygzh/internal/cli"
	"github.com/spf13/cobra"
)

func newTemplateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Manage built-in and user structure templates.",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List available structure templates.",
		RunE: func(cmd *cobra.Command, args []string) error {
			infos, err := templateManager().List()
			if err != nil {
				return emit(cmd, cli.Fail("TEMPLATE_LIST_FAILED", err.Error()))
			}
			if globalFlags.json {
				return emit(cmd, cli.OK("TEMPLATE_LIST", map[string]any{"templates": infos}))
			}
			for _, info := range infos {
				if info.Description != "" {
					fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", info.Name, info.Description)
				} else {
					fmt.Fprintln(cmd.OutOrStdout(), info.Name)
				}
			}
			return nil
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "show <name>",
		Short: "Print the HTML of a structure template.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tmpl, err := templateManager().Load(args[0])
			if err != nil {
				return emit(cmd, cli.Fail("TEMPLATE_NOT_FOUND", err.Error()))
			}
			if globalFlags.json {
				return emit(cmd, cli.OK("TEMPLATE_SHOW", map[string]any{
					"name":        tmpl.Name,
					"description": tmpl.Description,
					"css":         tmpl.CSS,
					"html":        tmpl.HTML,
				}))
			}
			fmt.Fprintf(cmd.OutOrStdout(), "name: %s\n", tmpl.Name)
			if tmpl.Description != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "description: %s\n", tmpl.Description)
			}
			fmt.Fprintln(cmd.OutOrStdout())
			if tmpl.CSS != "" {
				fmt.Fprintln(cmd.OutOrStdout(), "--- CSS ---")
				fmt.Fprint(cmd.OutOrStdout(), tmpl.CSS)
				fmt.Fprintln(cmd.OutOrStdout())
			}
			fmt.Fprintln(cmd.OutOrStdout(), "--- HTML ---")
			fmt.Fprint(cmd.OutOrStdout(), tmpl.HTML)
			return nil
		},
	})
	return cmd
}
