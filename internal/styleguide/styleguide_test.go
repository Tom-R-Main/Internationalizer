package styleguide

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadCombinesSharedAndLocaleGuidance(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "_conventions.md"), []byte("  Keep product names.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "fr.md"), []byte("\nUse formal French.  "), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := Load(dir, "fr")
	if err != nil {
		t.Fatal(err)
	}
	want := "Keep product names.\n\n---\n\nUse formal French."
	if got != want {
		t.Fatalf("Load() = %q, want %q", got, want)
	}
}

func TestLoadAllowsMissingGuides(t *testing.T) {
	got, err := Load(t.TempDir(), "fr")
	if err != nil {
		t.Fatal(err)
	}
	if got != "" {
		t.Fatalf("Load() = %q, want empty", got)
	}
}

func TestLoadReportsUnreadableGuide(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "_conventions.md"), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(dir, "fr"); err == nil {
		t.Fatal("Load() accepted a directory as a style-guide file")
	}
}
