package main

import (
	"github.com/spf13/cobra"
)

// image: process and upload images to WeChat permanent material (phase 4).
// No AI generation — generation is delegated to the agent.
func newImageCmd() *cobra.Command {
	var upload bool
	cmd := &cobra.Command{
		Use:   "image <image-path>",
		Short: "Process an image (compress) and optionally upload to WeChat material.",
		Long: `Compresses and validates an image for WeChat. With --upload, uploads it to
the permanent material library and returns the media_id + URL. Image GENERATION
is not performed here — delegate that to your agent.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runImage(cmd, args[0], upload)
		},
	}
	cmd.Flags().BoolVar(&upload, "upload", false, "upload to WeChat permanent material (requires creds)")
	return cmd
}
