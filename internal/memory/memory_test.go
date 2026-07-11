package memory

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestInitStoreAndValidate(t *testing.T) {
	store := filepath.Join(t.TempDir(), "memory")
	result, err := InitStoreWithOptions(store, false)
	if err != nil {
		t.Fatal(err)
	}
	if result.Dir == "" || result.BackupDir != "" {
		t.Fatalf("unexpected init result: %+v", result)
	}
	report, err := ValidateStore(store)
	if err != nil {
		t.Fatal(err)
	}
	if !report.Exists || !report.Valid {
		t.Fatalf("expected valid store, got %+v", report)
	}
	if len(report.Profiles) != 0 {
		t.Fatalf("unexpected profiles: %v", report.Profiles)
	}
}

func TestInitStoreForcePreservesBackup(t *testing.T) {
	store := filepath.Join(t.TempDir(), "memory")
	if err := InitStore(store); err != nil {
		t.Fatal(err)
	}
	preferences := filepath.Join(store, "preferences.md")
	if err := os.WriteFile(preferences, []byte("sentinel\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	result, err := InitStoreWithOptions(store, true)
	if err != nil {
		t.Fatal(err)
	}
	if result.BackupDir == "" {
		t.Fatal("expected a backup directory")
	}
	backupBody, err := os.ReadFile(filepath.Join(result.BackupDir, "preferences.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(backupBody) != "sentinel\n" {
		t.Fatalf("backup did not preserve old content: %q", backupBody)
	}
	report, err := ValidateStore(store)
	if err != nil || !report.Valid {
		t.Fatalf("replacement store invalid: report=%+v err=%v", report, err)
	}
}

func TestAddProfileUpdatesStoreAndValidates(t *testing.T) {
	store := filepath.Join(t.TempDir(), "memory")
	if err := InitStore(store); err != nil {
		t.Fatal(err)
	}
	path, err := AddProfile(store, "tech-notes")
	if err != nil {
		t.Fatal(err)
	}
	if path != filepath.Join(store, "profiles", "tech-notes") {
		t.Fatalf("unexpected path: %s", path)
	}
	body, err := os.ReadFile(filepath.Join(path, "identity.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(body), "main-account") || !strings.Contains(string(body), "account: tech-notes") {
		t.Fatalf("profile placeholders were not rewritten:\n%s", body)
	}
	folderMeta, err := os.ReadFile(filepath.Join(path, ".ok", "frontmatter.yml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(folderMeta), "tags: [profile, tech-notes]") || strings.Contains(string(folderMeta), "template") {
		t.Fatalf("profile folder metadata was not activated:\n%s", folderMeta)
	}
	if _, err := AddProfile(store, "product-news"); err != nil {
		t.Fatal(err)
	}
	report, err := ValidateStore(store)
	if err != nil || !report.Valid {
		t.Fatalf("store invalid after profile add: report=%+v err=%v", report, err)
	}
	want := []string{"product-news", "tech-notes"}
	if strings.Join(report.Profiles, ",") != strings.Join(want, ",") {
		t.Fatalf("profiles=%v want=%v", report.Profiles, want)
	}
	logBody, err := os.ReadFile(filepath.Join(store, "log.md"))
	if err != nil {
		t.Fatal(err)
	}
	heading := "## " + time.Now().UTC().Format("2006-01-02")
	if strings.Count(string(logBody), heading) != 1 {
		t.Fatalf("expected one daily log heading, got:\n%s", logBody)
	}
}

func TestAddProfileRejectsUnsafeName(t *testing.T) {
	store := filepath.Join(t.TempDir(), "memory")
	if err := InitStore(store); err != nil {
		t.Fatal(err)
	}
	if _, err := AddProfile(store, "../escape"); err == nil {
		t.Fatal("expected unsafe account name to be rejected")
	}
}

func TestValidateStoreFindsBrokenLinkAndSecret(t *testing.T) {
	store := filepath.Join(t.TempDir(), "memory")
	if err := InitStore(store); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(store, "preferences.md")
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	body = append(body, []byte("\n[missing](./does-not-exist.md)\napi_key = abcdefghijklmnopqrstuvwxyz\n")...)
	if err := os.WriteFile(path, body, 0o644); err != nil {
		t.Fatal(err)
	}
	report, err := ValidateStore(store)
	if err != nil {
		t.Fatal(err)
	}
	if report.Valid {
		t.Fatal("expected validation failure")
	}
	joined := strings.Join(report.Issues, "\n")
	if !strings.Contains(joined, "broken link") || !strings.Contains(joined, "secret-shaped") {
		t.Fatalf("expected broken link and secret issues, got:\n%s", joined)
	}
}
