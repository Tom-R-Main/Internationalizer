package pseudo

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
	"github.com/Tom-R-Main/Internationalizer/internal/state"
)

func TestGenerateWritesTrackedPseudoArtifactAndCanRefreshIt(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "locales", "en.json")
	targetTemplate := filepath.Join(dir, "locales", "{locale}.json")
	manifestPath := filepath.Join(dir, ".internationalizer.lock")
	writeTestFile(t, sourcePath, `{"greeting":"Hello {name}","link":"Read <a href=\"/docs\">docs</a>"}`)
	cfg := pseudoTestConfig(sourcePath, targetTemplate, manifestPath)

	results, err := Generate(cfg, GenerateOptions{Strategy: Accented})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if len(results) != 1 || !results[0].Written || results[0].Locale != "en-XA" || results[0].Units != 2 {
		t.Fatalf("Generate() results = %#v", results)
	}
	targetPath := filepath.Join(dir, "locales", "en-XA.json")
	data, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatal(err)
	}
	output := string(data)
	if !strings.Contains(output, "{name}") || !strings.Contains(output, `<a href=\"/docs\">`) || !strings.Contains(output, "[!!") {
		t.Fatalf("pseudo output did not preserve runtime syntax: %s", output)
	}
	manifest, err := state.Load(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	entry, ok := manifest.Get("app", "greeting", "en-XA")
	if !ok || entry.Origin != "pseudo" || entry.ReviewStatus != state.ReviewNeedsReview {
		t.Fatalf("pseudo manifest entry = %#v, present = %v", entry, ok)
	}

	// A tracked, byte-identical pseudo artifact is safe to regenerate.
	if _, err := Generate(cfg, GenerateOptions{Strategy: Accented}); err != nil {
		t.Fatalf("second Generate() error = %v", err)
	}
}

func TestGenerateRefusesUntrackedTargetUnlessForced(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	targetTemplate := filepath.Join(dir, "{locale}.json")
	writeTestFile(t, sourcePath, `{"greeting":"Hello"}`)
	writeTestFile(t, filepath.Join(dir, "en-XA.json"), `{"greeting":"Human text"}`)
	cfg := pseudoTestConfig(sourcePath, targetTemplate, filepath.Join(dir, "manifest.json"))

	_, err := Generate(cfg, GenerateOptions{Strategy: Accented})
	if err == nil || !strings.Contains(err.Error(), "refusing to overwrite") {
		t.Fatalf("Generate() error = %v, want overwrite refusal", err)
	}
	if _, err := Generate(cfg, GenerateOptions{Strategy: Accented, Force: true}); err != nil {
		t.Fatalf("forced Generate() error = %v", err)
	}
}

func TestGenerateDryRunDoesNotWriteArtifactOrManifest(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	targetTemplate := filepath.Join(dir, "{locale}.json")
	manifestPath := filepath.Join(dir, "manifest.json")
	writeTestFile(t, sourcePath, `{"greeting":"Hello"}`)
	cfg := pseudoTestConfig(sourcePath, targetTemplate, manifestPath)

	results, err := Generate(cfg, GenerateOptions{Strategy: Bidi, DryRun: true})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if len(results) != 1 || results[0].Written || results[0].Locale != "ar-XB" {
		t.Fatalf("Generate() results = %#v", results)
	}
	for _, path := range []string{filepath.Join(dir, "ar-XB.json"), manifestPath} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("dry run unexpectedly wrote %s", path)
		}
	}
}

func pseudoTestConfig(source, target, manifest string) *config.Config {
	return &config.Config{
		SourceLocale:  "en",
		TargetLocales: []string{"fr"},
		ManifestPath:  manifest,
		Bundles: []config.Bundle{{
			ID:     "app",
			Source: source,
			Target: target,
			Format: "json",
		}},
	}
}

func writeTestFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}
