package wechat

// Error maps WeChat errcode values to actionable messages. Covers the full
// common set — md2wechat only mapped 45002-45005.
type Error struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (e *Error) Error() string {
	if msg, ok := errcodeMessages[e.ErrCode]; ok {
		return msg
	}
	if e.ErrMsg != "" {
		return e.ErrMsg
	}
	return "unknown wechat error"
}

// IsCredentialError reports whether the error indicates the access_token is bad
// or expired (caller should refresh token and retry once).
func (e *Error) IsCredentialError() bool {
	switch e.ErrCode {
	case 40001, 40014, 42001:
		return true
	}
	return false
}

var errcodeMessages = map[int]string{
	-1:    "system busy, please retry",
	0:     "ok",
	40001: "invalid credential / access_token invalid — refresh token and retry",
	40002: "invalid grant type",
	40003: "invalid openid",
	40004: "invalid media type",
	40013: "invalid appid",
	40125: "invalid appsecret — check your WECHAT_APPID/WECHAT_SECRET",
	40164: "invalid source IP — add your IP to the WeChat IP whitelist",
	45002: "content too long (max 20000 chars / title 64)",
	45003: "title missing or too long (max 64)",
	45004: "author too long (max 8) — note: SDK field limit differs",
	45005: "digest too long (max 120)",
	45007: "image url too long",
	45009: "reach max api daily limit",
	45010: "content missing",
	45166: "reply content missing",
	48001: "API unauthorized — your account lacks this API (verify it is a verified service account)",
	48002: "fan must be in sandbox/test list",
	48004: "API requires verified account",
	48006: "forbidden: function not enabled",
	50001: "user not in app scope",
	50002: "user restricted",
	81000: "url domain not in jsapi safe domain list",
	// draft-specific
	40118: "invalid media_id for thumb (cover image not uploaded as permanent material)",
	45154: "appid not bound to this draft",
	46003: "draft not found / media_id mismatch",
	40007: "invalid media_id",
}
