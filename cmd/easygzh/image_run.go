package main

import (
	"github.com/easygzh/easygzh/internal/cli"
	"github.com/easygzh/easygzh/internal/image"
	"github.com/spf13/cobra"
)

// runImage compresses an image and optionally uploads it to WeChat material.
func runImage(cmd *cobra.Command, path string, upload bool) error {
	result, err := image.Process(path)
	if err != nil {
		return emit(cmd, cli.Fail("IMAGE_PROCESS_FAILED", err.Error()))
	}
	data := map[string]any{
		"original":   result.Original,
		"compressed": result.Compressed,
		"width":      result.Width,
		"height":     result.Height,
		"bytes":      result.Bytes,
	}
	if !upload {
		return emit(cmd, cli.OK("IMAGE_PROCESSED", data, "upload with --upload when ready"))
	}

	if err := ValidateUploadConfig(); err != nil {
		return emit(cmd, cli.Fail("WECHAT_NOT_CONFIGURED", err.Error()))
	}
	mediaID, url, err := image.Upload(result.Compressed)
	if err != nil {
		return emit(cmd, cli.Fail("IMAGE_UPLOAD_FAILED", err.Error()))
	}
	data["media_id"] = mediaID
	data["wechat_url"] = url
	return emit(cmd, cli.OK("IMAGE_UPLOADED", data, "use the media_id in: easygzh publish --cover-media-id"))
}

// ValidateUploadConfig mirrors wechat.ValidateConfig (kept here to avoid a cmd→wechat import cycle for the image flow).
func ValidateUploadConfig() error {
	return wechatValidateConfig()
}
