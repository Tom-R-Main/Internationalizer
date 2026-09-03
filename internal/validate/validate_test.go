package validate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
)

func TestValidateUsesBundleTargetTemplate(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "README.md")
	targetPath := filepath.Join(dir, "docs", "i18n", "fr.md")
	if err := os.WriteFile(sourcePath, []byte("# Hello\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(targetPath, []byte("# Bonjour\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := &config.Config{
		SourceLocale:  "en",
		TargetLocales: []string{"fr"},
		Bundles: []config.Bundle{{
			ID:     "docs",
			Source: sourcePath,
			Target: filepath.Join(dir, "docs", "i18n", "{locale}.md"),
		}},
	}

	reports, err := Validate(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(reports) != 1 || reports[0].Bundle != "docs" || reports[0].Coverage != 100 {
		t.Fatalf("reports = %#v", reports)
	}
	if HasFailures(reports) {
		t.Fatalf("valid translated document reported as failure: %#v", reports)
	}
}

func TestValidateReportsMalformedTargetAsFailure(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	targetPath := filepath.Join(dir, "fr.json")
	if err := os.WriteFile(sourcePath, []byte(`{"a":"A"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(targetPath, []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := &config.Config{
		TargetLocales: []string{"fr"},
		SourcePath:    sourcePath,
	}
	reports, err := Validate(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(reports) != 1 || len(reports[0].Errors) != 1 || !HasFailures(reports) {
		t.Fatalf("reports = %#v, want malformed target failure", reports)
	}
}

func TestFormatHumanPreservesInterpolationMismatchDetails(t *testing.T) {
	reports := []Report{{
		Bundle: "app",
		Locale: "fr",
		Mismatches: []Mismatch{{
			Key:        "welcome",
			SourceVars: []string{"name"},
			TargetVars: nil,
		}},
	}}

	got := FormatHuman(reports)
	if !strings.Contains(got, "mismatch: welcome (source: [name], target: [])") {
		t.Fatalf("FormatHuman omitted mismatch details: %q", got)
	}
}
