package wechat

import (
	"fmt"

	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	"github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/draft"
	"github.com/silenceper/wechat/v2/officialaccount/material"
)

// Client wraps silenceper/wechat's OfficialAccount with a file-backed token cache.
type Client struct {
	oa *officialaccount.OfficialAccount
}

// NewClient builds a Client from the given credentials, persisting the
// access_token to TokenCachePath().
func NewClient(cfg *Config) (*Client, error) {
	if cfg == nil {
		c, err := LoadConfig()
		if err != nil {
			return nil, err
		}
		cfg = c
	}
	fc := NewFileCache(TokenCachePath())
	// silenceper's cache.Cache interface — our FileCache satisfies it.
	var _ cache.Cache = (*FileCache)(nil)

	wc := officialaccount.NewOfficialAccount(&config.Config{
		AppID:     cfg.AppID,
		AppSecret: cfg.AppSecret,
		Cache:     fc,
	})
	return &Client{oa: wc}, nil
}

// UploadImage uploads a local image to permanent material and returns media_id + url.
func (c *Client) UploadImage(path string) (mediaID, url string, err error) {
	mat := c.oa.GetMaterial()
	mediaID, url, err = mat.AddMaterial(material.MediaTypeImage, path)
	if err != nil {
		return "", "", mapErr(err)
	}
	return mediaID, url, nil
}

// CreateDraft pushes one article to the draft box. Returns the draft media_id.
func (c *Client) CreateDraft(art *Article) (string, error) {
	d := c.oa.GetDraft()
	mediaID, err := d.AddDraft([]*draft.Article{{
		Title:              art.Title,
		Author:             art.Author,
		Digest:             art.Digest,
		Content:            art.Content,
		ThumbMediaID:       art.ThumbMediaID,
		ShowCoverPic:       1,
		NeedOpenComment:    0,
		OnlyFansCanComment: 0,
	}})
	if err != nil {
		return "", mapErr(err)
	}
	return mediaID, nil
}

// mapErr wraps a raw SDK error into our *Error with a human message when possible.
func mapErr(err error) error {
	if err == nil {
		return nil
	}
	// silenceper returns errors that embed the errcode as a string; attempt to
	// extract. For simplicity, return a generic Error carrying the message.
	return &Error{ErrMsg: err.Error()}
}

var _ = fmt.Sprintf
