// Package wechat wraps the WeChat Official Account SDK (silenceper/wechat/v2)
// with: file-backed access_token persistence (fixing the in-memory-only weakness
// of the upstream memory cache), full error-code mapping, and a publish pipeline
// that does NOT require any third-party paid key — only the user's own
// appid/secret.
package wechat

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// FileCache implements silenceper/wechat's cache.Cache interface backed by a
// JSON file, so the access_token survives across CLI invocations (tokens live
// ~2h; WeChat rate-limits fetches to 2000/day per appid).
type FileCache struct {
	mu   sync.Mutex
	Path string
}

type cacheEntry struct {
	Val    string    `json:"val"`
	Expire time.Time `json:"expire"`
}
type cacheFile map[string]cacheEntry

// NewFileCache returns a FileCache at the given path (e.g. ~/.easygzh/token.json).
func NewFileCache(path string) *FileCache {
	return &FileCache{Path: path}
}

func (f *FileCache) read() cacheFile {
	b, err := os.ReadFile(f.Path)
	if err != nil {
		return cacheFile{}
	}
	var m cacheFile
	if json.Unmarshal(b, &m) != nil {
		return cacheFile{}
	}
	return m
}

func (f *FileCache) write(m cacheFile) error {
	if err := os.MkdirAll(filepath.Dir(f.Path), 0o700); err != nil {
		return err
	}
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	// Atomic write: temp file + rename.
	tmp := f.Path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, f.Path)
}

// Get returns the cached value if present and unexpired, else nil.
func (f *FileCache) Get(key string) interface{} {
	f.mu.Lock()
	defer f.mu.Unlock()
	m := f.read()
	e, ok := m[key]
	if !ok || time.Now().After(e.Expire) {
		return nil
	}
	return e.Val
}

// Set stores val with the given timeout (from now).
func (f *FileCache) Set(key string, val interface{}, timeout time.Duration) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	m := f.read()
	s, _ := val.(string)
	m[key] = cacheEntry{Val: s, Expire: time.Now().Add(timeout)}
	return f.write(m)
}

// IsExist reports whether a non-expired entry exists (silenceper cache.Cache).
func (f *FileCache) IsExist(key string) bool {
	return f.Get(key) != nil
}

// Has is an alias kept for readability in non-SDK contexts.
func (f *FileCache) Has(key string) bool {
	return f.IsExist(key)
}

// Delete removes an entry.
func (f *FileCache) Delete(key string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	m := f.read()
	delete(m, key)
	return f.write(m)
}
