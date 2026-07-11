// Package image handles WeChat image processing: compress/resize to stay within
// WeChat limits (max width ~1920px, max size ~5MB), and upload to the permanent
// material library. AI generation is NOT done here — it's delegated to the agent.
package image

import (
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/easygzh/easygzh/internal/wechat"
)

const (
	MaxWidth  = 1920
	MaxSizeMB = 5
)

// Result describes a processed image.
type Result struct {
	Original   string
	Compressed string
	Width      int
	Height     int
	Bytes      int
}

// Process reads the image at path, downscales it to MaxWidth if needed, and
// writes a compressed copy to a temp file. Returns metadata.
func Process(path string) (*Result, error) {
	src, err := imaging.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open image: %w", err)
	}
	bounds := src.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	// Downscale if wider than MaxWidth.
	out := src
	if w > MaxWidth {
		out = imaging.Resize(src, MaxWidth, 0, imaging.Lanczos)
		w, h = out.Bounds().Dx(), out.Bounds().Dy()
	}

	// Write to a temp file in the same format (fallback to jpeg).
	ext := strings.ToLower(filepath.Ext(path))
	tmp, err := os.CreateTemp("", "easygzh-img-*"+ext)
	if err != nil {
		return nil, fmt.Errorf("create temp: %w", err)
	}
	tmpPath := tmp.Name()
	tmp.Close()

	if err := encodeImage(out, tmpPath, ext); err != nil {
		os.Remove(tmpPath)
		return nil, err
	}
	info, _ := os.Stat(tmpPath)

	return &Result{
		Original:   path,
		Compressed: tmpPath,
		Width:      w,
		Height:     h,
		Bytes:      int(info.Size()),
	}, nil
}

func encodeImage(img image.Image, path, ext string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	switch ext {
	case ".jpg", ".jpeg":
		return jpeg.Encode(f, img, &jpeg.Options{Quality: 85})
	case ".png":
		return png.Encode(f, img)
	default:
		return jpeg.Encode(f, img, &jpeg.Options{Quality: 85})
	}
}

// Upload pushes a local image to the WeChat permanent material library.
// Returns (media_id, url).
func Upload(path string) (mediaID, url string, err error) {
	cfg, _ := wechat.LoadConfig()
	client, err := wechat.NewClient(cfg)
	if err != nil {
		return "", "", err
	}
	return client.UploadImage(path)
}
