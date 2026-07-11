package main

// runPublish — phase 3 implements the full WeChat publish path. Phase 1 stub
// returns readiness guidance so the CLI compiles end-to-end.
import (
	"os"

	"github.com/easygzh/easygzh/internal/cli"
	"github.com/easygzh/easygzh/internal/wechat"
	"github.com/spf13/cobra"
)

func runPublish(cmd *cobra.Command, markdownFile string, opts publishOpts) error {
	// Read markdown.
	md, err := os.ReadFile(markdownFile)
	if err != nil {
		return emit(cmd, cli.Fail("READ_FAILED", err.Error()))
	}

	// Resolve theme CSS and render.
	css, err := themeManager().Load(globalFlags.theme)
	if err != nil {
		return emit(cmd, cli.Fail("THEME_NOT_FOUND", err.Error()))
	}
	html, err := renderPipeline(string(md), css)
	if err != nil {
		return emit(cmd, cli.Fail("CONVERT_FAILED", err.Error()))
	}

	// Build the draft article (local only at this stage).
	article := wechat.Article{
		Title:   opts.title,
		Author:  opts.author,
		Digest:  opts.digest,
		Content: html,
	}
	if opts.saveDraft != "" {
		if err := wechat.WriteDraftJSON(opts.saveDraft, &article); err != nil {
			return emit(cmd, cli.Fail("WRITE_FAILED", err.Error()))
		}
		return emit(cmd, cli.OK("DRAFT_SAVED_LOCAL", map[string]any{
			"file":  opts.saveDraft,
			"title": opts.title,
			"bytes": len(html),
		}, "review the draft JSON", "push to WeChat once creds are configured"))
	}

	// Live publish requires credentials.
	if err := wechat.ValidateConfig(); err != nil {
		return emit(cmd, cli.Fail("WECHAT_NOT_CONFIGURED",
			"live publish requires WECHAT_APPID and WECHAT_SECRET. "+err.Error()))
	}
	result, err := wechat.Publish(string(md), html, wechat.PublishOptions{
		Title:   opts.title,
		Author:  opts.author,
		Digest:  opts.digest,
		Cover:   opts.cover,
		Account: opts.account,
	})
	if err != nil {
		return emit(cmd, cli.Fail("PUBLISH_FAILED", err.Error()))
	}
	return emit(cmd, cli.OK("DRAFT_CREATED", result, "publish via: easygzh publish --release (freepublish)"))
}
