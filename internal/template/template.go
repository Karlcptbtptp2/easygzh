// Package template loads built-in and user content-structure templates.
//
// A structure template (.html file) wraps goldmark-rendered HTML in a richer
// narrative layout (brand header, hook, body, CTA, footer) that Markdown alone
// cannot express. This package mirrors internal/theme: Load by name, List
// available templates. Parsing is delegated to render.ParseTemplate.
package template

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/easygzh/easygzh/internal/render"
)

// Manager resolves a template name to its parsed StructureTemplate.
type Manager struct {
	// TemplatesDir is where .html template files live (e.g. "./templates").
	TemplatesDir string
}

// TemplateInfo is a list entry describing one available template.
type TemplateInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Load reads and parses the named template from TemplatesDir/<name>.html.
func (m *Manager) Load(name string) (*render.StructureTemplate, error) {
	if name == "" {
		return nil, fmt.Errorf("template: name is required")
	}
	path := filepath.Join(m.TemplatesDir, name+".html")
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("load template %q from %s: %w", name, path, err)
	}
	tmpl, err := render.ParseTemplate(string(raw))
	if err != nil {
		return nil, fmt.Errorf("parse template %q: %w", name, err)
	}
	return tmpl, nil
}

// List returns info about all available templates in TemplatesDir.
func (m *Manager) List() ([]TemplateInfo, error) {
	entries, err := os.ReadDir(m.TemplatesDir)
	if err != nil {
		return nil, fmt.Errorf("list templates in %s: %w", m.TemplatesDir, err)
	}
	var infos []TemplateInfo
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if filepath.Ext(name) != ".html" {
			continue
		}
		baseName := name[:len(name)-len(".html")]
		// Read the file to extract description.
		desc := ""
		if raw, err := os.ReadFile(filepath.Join(m.TemplatesDir, name)); err == nil {
			if tmpl, err := render.ParseTemplate(string(raw)); err == nil {
				desc = tmpl.Description
			}
		}
		infos = append(infos, TemplateInfo{
			Name:        baseName,
			Description: desc,
		})
	}
	return infos, nil
}

// DefaultTemplatesDir returns the conventional location: repo root / templates.
// The CLI resolves via the EASYGZH_TEMPLATES_DIR env var or defaults to
// "./templates" (relative to the working directory).
func DefaultTemplatesDir() string {
	if d := os.Getenv("EASYGZH_TEMPLATES_DIR"); d != "" {
		return d
	}
	return "./templates"
}
