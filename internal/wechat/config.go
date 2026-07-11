package wechat

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the WeChat account credentials.
type Config struct {
	AppID     string
	AppSecret string
	Account   string // named account, for multi-account resolution
}

// LoadConfig resolves credentials from env (and, in phase 3+, from a config
// file). Env vars: WECHAT_APPID, WECHAT_SECRET, WECHAT_ACCOUNT.
func LoadConfig() (*Config, error) {
	cfg := &Config{
		AppID:     os.Getenv("WECHAT_APPID"),
		AppSecret: os.Getenv("WECHAT_SECRET"),
		Account:   os.Getenv("WECHAT_ACCOUNT"),
	}
	return cfg, nil
}

// ValidateConfig returns an error if credentials required for live publish are
// missing. This is the gate md2wechat gates behind a PAID key — easygzh only
// needs the user's own appid/secret.
func ValidateConfig() error {
	cfg, _ := LoadConfig()
	if cfg.AppID == "" || cfg.AppSecret == "" {
		return fmt.Errorf("WECHAT_APPID and WECHAT_SECRET must be set for live publish")
	}
	return nil
}

// TokenCachePath returns the file used for access_token persistence.
func TokenCachePath() string {
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".easygzh", "token.json")
	}
	return "./token.json"
}
