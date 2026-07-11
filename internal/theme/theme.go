// Package theme loads built-in and user CSS themes.
//
// In the full design (phase 2), a theme is loaded by name: built-in themes ship
// embedded in the binary; profile overrides layer on top. This MVP phase-1
// implementation loads built-in themes from the themes/ directory relative to
// the repo root (and is replaced by //go:embed in phase 4).
package theme

import (
	"fmt"
	"os"
	"path/filepath"
)

// Manager resolves a theme name to its CSS.
type Manager struct {
	// ThemesDir is where built-in theme .css files live (e.g. "./themes").
	ThemesDir string
}

// Load returns the CSS for the named theme, read from ThemesDir/<name>.css.
func (m *Manager) Load(name string) (string, error) {
	if name == "" {
		name = "default"
	}
	path := filepath.Join(m.ThemesDir, name+".css")
	css, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("load theme %q from %s: %w", name, path, err)
	}
	return string(css), nil
}

// List returns the names of available built-in themes.
func (m *Manager) List() ([]string, error) {
	entries, err := os.ReadDir(m.ThemesDir)
	if err != nil {
		return nil, fmt.Errorf("list themes in %s: %w", m.ThemesDir, err)
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if filepath.Ext(name) == ".css" {
			names = append(names, name[:len(name)-len(".css")])
		}
	}
	return names, nil
}

// DefaultThemesDir returns the conventional location: repo root / themes.
// The CLI resolves the repo root via the EASYGZH_THEMES_DIR env var or defaults
// to "./themes" (relative to the working directory).
func DefaultThemesDir() string {
	if d := os.Getenv("EASYGZH_THEMES_DIR"); d != "" {
		return d
	}
	return "./themes"
}
