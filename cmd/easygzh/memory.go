package main

import (
	"fmt"

	"github.com/easygzh/easygzh/internal/cli"
	"github.com/easygzh/easygzh/internal/memory"
	"github.com/spf13/cobra"
)

func newMemoryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "memory",
		Short: "Manage the local OpenKnowledge-compatible tone memory store.",
	}

	var force bool
	initCmd := &cobra.Command{
		Use:   "init [target-dir]",
		Short: "Initialize a validated memory store from the embedded scaffold.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := memory.DefaultDir()
			if len(args) > 0 {
				target = args[0]
			}
			result, err := memory.InitStoreWithOptions(target, force)
			if err != nil {
				return emit(cmd, cli.Fail("MEMORY_INIT_FAILED", err.Error()))
			}
			report, err := memory.ValidateStore(result.Dir)
			if err != nil {
				return emit(cmd, cli.Fail("MEMORY_VALIDATE_FAILED", err.Error()))
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "✓ memory store initialized at %s\n", result.Dir)
			return emit(cmd, cli.OK("MEMORY_INIT_OK", map[string]any{
				"dir":        result.Dir,
				"backup_dir": result.BackupDir,
				"validation": report,
			}, "Edit a profile or run: easygzh memory profile add <account>"))
		},
	}
	initCmd.Flags().BoolVar(&force, "force", false, "replace an existing store after preserving a timestamped backup")
	cmd.AddCommand(initCmd)

	cmd.AddCommand(&cobra.Command{
		Use:   "profiles",
		Short: "List account profiles in the memory store.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			store := memory.DefaultDir()
			names, err := memory.ListProfiles(store)
			if err != nil {
				return emit(cmd, cli.Fail("MEMORY_LIST_FAILED", err.Error()))
			}
			if globalFlags.json {
				return emit(cmd, cli.OK("MEMORY_PROFILES", map[string]any{"profiles": names, "dir": store}))
			}
			for _, name := range names {
				fmt.Fprintln(cmd.OutOrStdout(), name)
			}
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "path",
		Short: "Print the memory store path.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), memory.DefaultDir())
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Report whether the memory store exists and passes validation.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			report, err := memory.ValidateStore(memory.DefaultDir())
			if err != nil {
				return emit(cmd, cli.Fail("MEMORY_STATUS_FAILED", err.Error()))
			}
			if globalFlags.json {
				return emit(cmd, cli.OK("MEMORY_STATUS", report))
			}
			fmt.Fprintf(cmd.OutOrStdout(), "dir: %s\nexists: %t\nvalid: %t\nprofiles: %v\n", report.Dir, report.Exists, report.Valid, report.Profiles)
			for _, issue := range report.Issues {
				fmt.Fprintln(cmd.OutOrStdout(), "issue: "+issue)
			}
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "validate",
		Short: "Validate metadata, internal links, profile shape and secret safety.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			report, err := memory.ValidateStore(memory.DefaultDir())
			if err != nil {
				return emit(cmd, cli.Fail("MEMORY_VALIDATE_FAILED", err.Error()))
			}
			if !report.Valid {
				return emit(cmd, cli.Fail("MEMORY_INVALID", fmt.Sprintf("memory validation failed: %v", report.Issues)))
			}
			if !globalFlags.json {
				fmt.Fprintf(cmd.OutOrStdout(), "✓ memory store valid (%d profiles)\n", len(report.Profiles))
			}
			return emit(cmd, cli.OK("MEMORY_VALID", report))
		},
	})

	profileCmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage account profiles.",
	}
	profileCmd.AddCommand(&cobra.Command{
		Use:   "add <account>",
		Short: "Create and validate a lowercase-kebab account profile.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store := memory.DefaultDir()
			path, err := memory.AddProfile(store, args[0])
			if err != nil {
				return emit(cmd, cli.Fail("MEMORY_PROFILE_ADD_FAILED", err.Error()))
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "✓ profile created at %s\n", path)
			return emit(cmd, cli.OK("MEMORY_PROFILE_ADDED", map[string]any{"account": args[0], "path": path, "dir": store}))
		},
	})
	cmd.AddCommand(profileCmd)

	return cmd
}
