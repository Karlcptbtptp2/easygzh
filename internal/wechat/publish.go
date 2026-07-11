package wechat

import (
	"encoding/json"
	"os"
)

// Article is the draft article payload.
type Article struct {
	Title        string `json:"title"`
	Author       string `json:"author,omitempty"`
	Digest       string `json:"digest,omitempty"`
	Content      string `json:"content"`
	ThumbMediaID string `json:"thumb_media_id,omitempty"`
}

// PublishOptions configures a publish run.
type PublishOptions struct {
	Title   string
	Author  string
	Digest  string
	Cover   string // path to cover image (uploaded if set)
	Account string
}

// PublishResult describes the outcome.
type PublishResult struct {
	MediaID      string `json:"media_id"`
	CoverMediaID string `json:"cover_media_id,omitempty"`
	Title        string `json:"title"`
}

// Publish runs: (optional cover upload) → draft add. Returns the draft media_id.
// Image generation is NOT done here (delegated to the agent); images inside the
// content HTML are expected to already be absolute URLs (WeChat re-hosts them on
// paste only via the web editor — via the API they must be uploaded first; for
// the MVP we support cover upload and content with already-uploaded/absolute URLs).
func Publish(markdown, html string, opts PublishOptions) (*PublishResult, error) {
	cfg, _ := LoadConfig()
	client, err := NewClient(cfg)
	if err != nil {
		return nil, err
	}

	art := &Article{
		Title:   opts.Title,
		Author:  opts.Author,
		Digest:  opts.Digest,
		Content: html,
	}

	if opts.Cover != "" {
		mediaID, _, err := client.UploadImage(opts.Cover)
		if err != nil {
			return nil, err
		}
		art.ThumbMediaID = mediaID
	}

	mediaID, err := client.CreateDraft(art)
	if err != nil {
		return nil, err
	}
	return &PublishResult{MediaID: mediaID, Title: opts.Title}, nil
}

// WriteDraftJSON writes the article as a draft JSON file locally (no network).
// Used for testing the pipeline without WeChat credentials.
func WriteDraftJSON(path string, art *Article) error {
	b, err := json.MarshalIndent(art, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}
