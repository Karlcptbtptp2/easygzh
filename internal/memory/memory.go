// Package memory manages the local easyGZH tone memory store.
//
// The store is a small OpenKnowledge-compatible Markdown bundle. The canonical
// scaffold is generated into the binary, so `easygzh memory init` works from a
// standalone release without depending on the source repository.
package memory

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

var (
	accountNamePattern  = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	markdownLinkPattern = regexp.MustCompile(`\[[^\]]+\]\(([^)]+\.md(?:#[^)]+)?)\)`)
	secretPatterns      = []*regexp.Regexp{
		regexp.MustCompile(`(?i)sk-[a-z0-9_-]{20,}`),
		regexp.MustCompile(`(?i)bearer\s+[a-z0-9._~+/-]{20,}={0,2}`),
		regexp.MustCompile(`(?i)(api[_-]?key|access[_-]?token|password|secret)\s*[:=]\s*["']?[a-z0-9._~+/-]{12,}`),
	}
)

// DefaultDir returns the conventional memory store location.
func DefaultDir() string {
	if d := os.Getenv("EASYGZH_MEMORY_DIR"); d != "" {
		return d
	}
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".easygzh", "memory")
	}
	return filepath.Join(".", ".easygzh", "memory")
}

// InitResult describes a completed store initialization.
type InitResult struct {
	Dir       string `json:"dir"`
	BackupDir string `json:"backup_dir,omitempty"`
}

// ValidationReport is returned by ValidateStore and Status.
type ValidationReport struct {
	Dir      string   `json:"dir"`
	Exists   bool     `json:"exists"`
	Valid    bool     `json:"valid"`
	Profiles []string `json:"profiles"`
	Issues   []string `json:"issues"`
}

// InitStore creates a new store without overwriting an existing one.
func InitStore(target string) error {
	_, err := InitStoreWithOptions(target, false)
	return err
}

// InitStoreWithOptions initializes target from the scaffold embedded in the
// binary. When force is true, an existing store is moved to a timestamped
// sibling backup before the replacement is installed.
func InitStoreWithOptions(target string, force bool) (InitResult, error) {
	abs, err := filepath.Abs(filepath.Clean(target))
	if err != nil {
		return InitResult{}, fmt.Errorf("resolve memory target: %w", err)
	}
	parent := filepath.Dir(abs)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return InitResult{}, fmt.Errorf("create memory parent: %w", err)
	}

	exists, err := pathState(abs)
	if err != nil {
		return InitResult{}, err
	}
	if exists {
		if !force {
			return InitResult{}, fmt.Errorf("target already exists: %s", abs)
		}
		info, err := os.Lstat(abs)
		if err != nil {
			return InitResult{}, err
		}
		if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
			return InitResult{}, fmt.Errorf("refusing to replace non-directory or symlink target: %s", abs)
		}
	}

	tmp, err := os.MkdirTemp(parent, ".easygzh-memory-init-")
	if err != nil {
		return InitResult{}, fmt.Errorf("create temporary memory store: %w", err)
	}
	cleanupTmp := true
	defer func() {
		if cleanupTmp {
			_ = os.RemoveAll(tmp)
		}
	}()

	if err := writeEmbeddedScaffold(tmp); err != nil {
		return InitResult{}, err
	}
	report, err := ValidateStore(tmp)
	if err != nil {
		return InitResult{}, err
	}
	if !report.Valid {
		return InitResult{}, fmt.Errorf("embedded memory scaffold is invalid: %s", strings.Join(report.Issues, "; "))
	}

	result := InitResult{Dir: abs}
	if exists {
		result.BackupDir = nextBackupPath(abs, time.Now().UTC())
		if err := os.Rename(abs, result.BackupDir); err != nil {
			return InitResult{}, fmt.Errorf("backup existing memory store: %w", err)
		}
	}

	if err := os.Rename(tmp, abs); err != nil {
		if result.BackupDir != "" {
			_ = os.Rename(result.BackupDir, abs)
		}
		return InitResult{}, fmt.Errorf("install memory store: %w", err)
	}
	cleanupTmp = false
	return result, nil
}

// StoreExists reports whether store is an existing directory.
func StoreExists(store string) (bool, error) {
	info, err := os.Stat(store)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.IsDir(), nil
}

// ListProfiles returns sorted account profile names in <store>/profiles/.
func ListProfiles(store string) ([]string, error) {
	exists, err := StoreExists(store)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("memory store not initialized: %s", store)
	}
	dir := filepath.Join(store, "profiles")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	names := []string{}
	for _, e := range entries {
		if e.IsDir() && !strings.HasPrefix(e.Name(), ".") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}

// ProfilePath returns the directory of a named profile.
func ProfilePath(store, account string) string {
	return filepath.Join(store, "profiles", account)
}

// AddProfile creates an account profile from the bundled hidden template,
// updates the profile index and log, and validates the resulting store.
func AddProfile(store, account string) (string, error) {
	if !accountNamePattern.MatchString(account) {
		return "", fmt.Errorf("account name must be lowercase kebab-case: %q", account)
	}
	report, err := ValidateStore(store)
	if err != nil {
		return "", err
	}
	if !report.Exists {
		return "", fmt.Errorf("memory store not initialized: %s", store)
	}
	if !report.Valid {
		return "", fmt.Errorf("memory store is invalid: %s", strings.Join(report.Issues, "; "))
	}

	src := filepath.Join(store, "profiles", ".template")
	dst := ProfilePath(store, account)
	if _, err := os.Stat(dst); err == nil {
		return "", fmt.Errorf("profile already exists: %s", account)
	} else if !os.IsNotExist(err) {
		return "", err
	}
	indexPath := filepath.Join(store, "profiles", "index.md")
	logPath := filepath.Join(store, "log.md")
	indexBefore, err := os.ReadFile(indexPath)
	if err != nil {
		return "", err
	}
	logBefore, err := os.ReadFile(logPath)
	if err != nil {
		return "", err
	}
	rollback := func() {
		_ = os.RemoveAll(dst)
		_ = atomicWriteFile(indexPath, indexBefore, 0o644)
		_ = atomicWriteFile(logPath, logBefore, 0o644)
	}

	if err := copyTree(src, dst); err != nil {
		return "", fmt.Errorf("copy profile template: %w", err)
	}
	if err := rewriteProfile(dst, account); err != nil {
		rollback()
		return "", err
	}
	if err := appendProfileIndex(store, account); err != nil {
		rollback()
		return "", err
	}
	if err := prependLogEntry(store, account); err != nil {
		rollback()
		return "", err
	}

	report, err = ValidateStore(store)
	if err != nil {
		rollback()
		return "", err
	}
	if !report.Valid {
		rollback()
		return "", fmt.Errorf("profile created but validation failed: %s", strings.Join(report.Issues, "; "))
	}
	return dst, nil
}

// ValidateStore checks required files, Markdown frontmatter, internal links,
// profile completeness and common secret-shaped values.
func ValidateStore(store string) (ValidationReport, error) {
	abs, err := filepath.Abs(filepath.Clean(store))
	if err != nil {
		return ValidationReport{}, err
	}
	report := ValidationReport{Dir: abs, Profiles: []string{}, Issues: []string{}}
	exists, err := StoreExists(abs)
	if err != nil {
		return report, err
	}
	report.Exists = exists
	if !exists {
		report.Issues = append(report.Issues, "memory store does not exist")
		return report, nil
	}

	required := []string{
		".ok/config.yml",
		"index.md",
		"log.md",
		"preferences.md",
		"profiles/index.md",
		"themes/index.md",
	}
	for _, rel := range required {
		if info, err := os.Stat(filepath.Join(abs, filepath.FromSlash(rel))); err != nil || info.IsDir() {
			report.Issues = append(report.Issues, "missing required file: "+rel)
		}
	}

	err = filepath.WalkDir(abs, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".md") {
			return nil
		}
		rel, _ := filepath.Rel(abs, path)
		rel = filepath.ToSlash(rel)
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		text := string(body)
		frontmatter, ok := parseFrontmatter(text)
		if !ok {
			report.Issues = append(report.Issues, rel+": missing YAML frontmatter")
		} else {
			for _, key := range []string{"type", "title", "description"} {
				if strings.TrimSpace(frontmatter[key]) == "" {
					report.Issues = append(report.Issues, rel+": missing frontmatter field "+key)
				}
			}
		}
		for _, pattern := range secretPatterns {
			if pattern.MatchString(text) {
				report.Issues = append(report.Issues, rel+": contains a secret-shaped value")
				break
			}
		}
		for _, match := range markdownLinkPattern.FindAllStringSubmatch(text, -1) {
			href := strings.SplitN(match[1], "#", 2)[0]
			if strings.Contains(href, "://") {
				continue
			}
			var target string
			if strings.HasPrefix(href, "/") {
				target = filepath.Join(abs, filepath.FromSlash(strings.TrimPrefix(href, "/")))
			} else {
				target = filepath.Join(filepath.Dir(path), filepath.FromSlash(href))
			}
			target = filepath.Clean(target)
			if !withinRoot(abs, target) {
				report.Issues = append(report.Issues, rel+": link escapes store: "+href)
				continue
			}
			if info, err := os.Stat(target); err != nil || info.IsDir() {
				report.Issues = append(report.Issues, rel+": broken link: "+href)
			}
		}
		return nil
	})
	if err != nil {
		return report, err
	}

	profiles, err := ListProfiles(abs)
	if err != nil {
		report.Issues = append(report.Issues, err.Error())
	} else {
		report.Profiles = profiles
		for _, account := range profiles {
			for _, rel := range []string{".ok/frontmatter.yml", "index.md", "identity.md", "visual-tone.md", "structure-tone.md", "current-theme.md", "samples/index.md"} {
				path := filepath.Join(ProfilePath(abs, account), filepath.FromSlash(rel))
				if info, err := os.Stat(path); err != nil || info.IsDir() {
					report.Issues = append(report.Issues, fmt.Sprintf("profile %s missing %s", account, rel))
					continue
				}
				body, err := os.ReadFile(path)
				if err != nil {
					return report, err
				}
				if strings.HasSuffix(rel, ".yml") {
					if !strings.Contains(string(body), account) || strings.Contains(string(body), "tags: [profile, template]") {
						report.Issues = append(report.Issues, fmt.Sprintf("profile %s has template folder metadata", account))
					}
					continue
				}
				frontmatter, ok := parseFrontmatter(string(body))
				if ok && frontmatter["account"] != account {
					report.Issues = append(report.Issues, fmt.Sprintf("profile %s has mismatched account in %s", account, rel))
				}
			}
		}
	}
	sort.Strings(report.Issues)
	report.Valid = len(report.Issues) == 0
	return report, nil
}

func parseFrontmatter(text string) (map[string]string, bool) {
	scanner := bufio.NewScanner(strings.NewReader(text))
	if !scanner.Scan() || strings.TrimSpace(scanner.Text()) != "---" {
		return nil, false
	}
	values := map[string]string{}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "---" {
			return values, true
		}
		if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") || strings.HasPrefix(strings.TrimSpace(line), "-") {
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		values[strings.TrimSpace(key)] = strings.Trim(strings.TrimSpace(value), `"'`)
	}
	return nil, false
}

func rewriteProfile(root, account string) error {
	return filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		text := strings.ReplaceAll(string(body), "main-account", account)
		text = strings.ReplaceAll(text, "主号", account)
		if strings.HasSuffix(filepath.ToSlash(path), ".ok/frontmatter.yml") {
			text = strings.ReplaceAll(text, "Profile 模板", account+" profile")
			text = strings.ReplaceAll(text, "仅供创建真实账号时复制的隐藏模板，不应作为激活 profile 读取。", "账号 "+account+" 的完整调性记忆。")
			text = strings.ReplaceAll(text, "tags: [profile, template]", "tags: [profile, "+account+"]")
		}
		return atomicWriteFile(path, []byte(text), 0o644)
	})
}

func appendProfileIndex(store, account string) error {
	path := filepath.Join(store, "profiles", "index.md")
	body, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	line := fmt.Sprintf("- [%s](%s/index.md) — 由 `easygzh memory profile add` 创建\n", account, account)
	if strings.Contains(string(body), "]("+account+"/index.md)") {
		return nil
	}
	text := strings.TrimRight(string(body), "\n") + "\n" + line
	return atomicWriteFile(path, []byte(text), 0o644)
}

func prependLogEntry(store, account string) error {
	path := filepath.Join(store, "log.md")
	body, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	marker := "# 变更日志\n"
	dateHeading := "## " + time.Now().UTC().Format("2006-01-02") + "\n"
	entry := fmt.Sprintf("\n**Profile** — 创建 `%s` profile，并通过记忆库校验。\n", account)
	text := string(body)
	if !strings.Contains(text, marker) {
		return errors.New("log.md is missing '# 变更日志' heading")
	}
	if strings.Contains(text, dateHeading) {
		text = strings.Replace(text, dateHeading, dateHeading+entry, 1)
	} else {
		text = strings.Replace(text, marker, marker+"\n"+dateHeading+entry, 1)
	}
	return atomicWriteFile(path, []byte(text), 0o644)
}

func writeEmbeddedScaffold(target string) error {
	paths := make([]string, 0, len(embeddedScaffoldFiles))
	for rel := range embeddedScaffoldFiles {
		paths = append(paths, rel)
	}
	sort.Strings(paths)
	for _, rel := range paths {
		clean := filepath.Clean(filepath.FromSlash(rel))
		if clean == "." || filepath.IsAbs(clean) || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
			return fmt.Errorf("invalid embedded scaffold path: %s", rel)
		}
		path := filepath.Join(target, clean)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := atomicWriteFile(path, []byte(embeddedScaffoldFiles[rel]), 0o644); err != nil {
			return fmt.Errorf("write embedded scaffold %s: %w", rel, err)
		}
	}
	return nil
}

func withinRoot(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func pathState(path string) (bool, error) {
	_, err := os.Lstat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func nextBackupPath(target string, now time.Time) string {
	base := target + ".backup-" + now.Format("20060102T150405Z")
	path := base
	for i := 2; ; i++ {
		if _, err := os.Lstat(path); os.IsNotExist(err) {
			return path
		}
		path = fmt.Sprintf("%s-%d", base, i)
	}
}
