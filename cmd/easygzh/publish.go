package main

import (
	"github.com/spf13/cobra"
)

// publish orchestrates convert → image upload → draft add (→ optional freepublish).
// Phase 3 implements the full WeChat publishing path. Phase 1 stub returns guidance.
func newPublishCmd() *cobra.Command {
	var (
		saveDraft string
		dryRun    bool
		account   string
		title     string
		author    string
		digest    string
		cover     string
	)
	cmd := &cobra.Command{
		Use:   "publish <markdown-file>",
		Short: "Render, upload images, and push a draft to the WeChat draft box.",
		Long: `Full publish pipeline: convert Markdown → inline HTML, upload images and
cover to WeChat permanent material, create a draft. With --save-draft, the draft
JSON is written locally instead of pushed (useful for testing).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPublish(cmd, args[0], publishOpts{
				saveDraft: saveDraft,
				dryRun:    dryRun,
				account:   account,
				title:     title,
				author:    author,
				digest:    digest,
				cover:     cover,
			})
		},
	}
	cmd.Flags().StringVar(&saveDraft, "save-draft", "", "write the draft JSON to this file instead of pushing to WeChat")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "prepare the draft but do not push")
	cmd.Flags().StringVar(&account, "account", "", "named account from config (default: wechat.default_account)")
	cmd.Flags().StringVar(&title, "title", "", "article title (required for draft)")
	cmd.Flags().StringVar(&author, "author", "", "article author")
	cmd.Flags().StringVar(&digest, "digest", "", "article digest/summary")
	cmd.Flags().StringVar(&cover, "cover", "", "cover image path")
	return cmd
}

type publishOpts struct {
	saveDraft string
	dryRun    bool
	account   string
	title     string
	author    string
	digest    string
	cover     string
}
