package main

import "github.com/easygzh/easygzh/internal/wechat"

// wechatValidateConfig delegates to the wechat package, kept here to avoid
// import-cycle noise in the image flow.
func wechatValidateConfig() error {
	return wechat.ValidateConfig()
}
