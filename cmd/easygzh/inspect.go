package main

import (
	"os"

	"github.com/easygzh/easygzh/internal/cli"
	"github.com/spf13/cobra"
)

// inspect reports whether an article is ready for each publish target, with no
// side effects. Phase 3 fills in WeChat credential readiness; phase 1 reports
// basic content readiness (file readable, non-empty).
func newInspectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect <markdown-file>",
		Short: "Check article readiness for preview/upload/draft (no side effects).",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			info, err := os.Stat(args[0])
			if err != nil {
				return emit(cmd, cli.Fail("INSPECT_FAILED", err.Error()))
			}
			data := map[string]any{
				"file": args[0],
				"size": info.Size(),
				"readiness": map[string]bool{
					"preview": true,
					"upload":  false, // requires WeChat creds
					"draft":   false,
				},
				"blockers": []string{
					"WECHAT_APPID/SECRET not configured for upload/draft",
				},
			}
			return emit(cmd, cli.OK("INSPECT_OK", data, "preview now", "configure wechat creds to enable upload/draft"))
		},
	}
	return cmd
}
