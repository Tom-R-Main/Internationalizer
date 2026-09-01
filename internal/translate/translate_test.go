package translate

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
	"github.com/Tom-R-Main/Internationalizer/internal/llm"
)

type fakeProvider struct {
	response *llm.TranslateResponse
	calls    int
}

func (p *fakeProvider) Name() string { return "fake" }

func (p *fakeProvider) Translate(_ context.Context, _ llm.TranslateRequest) (*llm.TranslateResponse, error) {
	p.calls++
	return p.response, nil
}

func TestRunRejectsIncompleteProviderResponseWithoutWritingTarget(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	if err := os.WriteFile(sourcePath, []byte(`{"a":"A","b":"B"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := testConfig(dir, sourcePath)
	provider := &fakeProvider{response: &llm.TranslateResponse{
		Translations: map[string]string{"a": "Un"},
	}}

	results, err := Run(context.Background(), cfg, provider, Options{})
	if err == nil {
		t.Fatal("Run returned nil error for an incomplete provider response")
	}
	if len(results) != 1 || len(results[0].Errors) == 0 {
		t.Fatalf("Run results = %#v, want one locale error", results)
	}
	if _, statErr := os.Stat(filepath.Join(dir, "fr.json")); !os.IsNotExist(statErr) {
		t.Fatalf("target was written after incomplete response; stat error = %v", statErr)
	}
}

func TestRunRejectsMalformedTargetWithoutOverwritingIt(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	targetPath := filepath.Join(dir, "fr.json")
	if err := os.WriteFile(sourcePath, []byte(`{"a":"A"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	malformed := []byte(`{"a":`)
	if err := os.WriteFile(targetPath, malformed, 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := testConfig(dir, sourcePath)
	provider := &fakeProvider{response: &llm.TranslateResponse{
		Translations: map[string]string{"a": "Un"},
	}}

	results, err := Run(context.Background(), cfg, provider, Options{})
	if err == nil {
		t.Fatal("Run returned nil error for a malformed target")
	}
	if len(results) != 1 || len(results[0].Errors) == 0 {
		t.Fatalf("Run results = %#v, want one locale error", results)
	}
	if provider.calls != 0 {
		t.Fatalf("provider called %d times for malformed target, want 0", provider.calls)
	}
	got, readErr := os.ReadFile(targetPath)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(got) != string(malformed) {
		t.Fatalf("malformed target changed: got %q, want %q", got, malformed)
	}
}

func TestRunRejectsInterpolationMismatchWithoutWritingTarget(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	if err := os.WriteFile(sourcePath, []byte(`{"welcome":"Hello, {{name}}"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := testConfig(dir, sourcePath)
	provider := &fakeProvider{response: &llm.TranslateResponse{
		Translations: map[string]string{"welcome": "Bonjour"},
	}}

	results, err := Run(context.Background(), cfg, provider, Options{})
	if err == nil {
		t.Fatal("Run returned nil error for interpolation mismatch")
	}
	if len(results) != 1 || len(results[0].Errors) == 0 {
		t.Fatalf("results = %#v, want interpolation error", results)
	}
	if _, statErr := os.Stat(filepath.Join(dir, "fr.json")); !os.IsNotExist(statErr) {
		t.Fatalf("target was written after interpolation mismatch; stat error = %v", statErr)
	}
}

func TestRunUsesBundleTargetTemplateAndAdoptsExistingTranslation(t *testing.T) {
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

	cfg := testConfig(dir, "")
	cfg.Bundles = []config.Bundle{{
		ID:     "docs",
		Source: sourcePath,
		Target: filepath.Join(dir, "docs", "i18n", "{locale}.md"),
		Format: "markdown",
	}}

	results, err := Run(context.Background(), cfg, nil, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if results[0].KeysUntracked != 1 {
		t.Fatalf("ordinary run results = %#v, want untracked entry", results)
	}
	if _, statErr := os.Stat(cfg.ManifestPath); !os.IsNotExist(statErr) {
		t.Fatalf("ordinary run implicitly adopted existing translation; stat error = %v", statErr)
	}

	results, err = Run(context.Background(), cfg, nil, Options{AdoptExisting: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].KeysUntracked != 1 || results[0].KeysMissing != 0 {
		t.Fatalf("first run results = %#v, want one adopted entry", results)
	}

	results, err = Run(context.Background(), cfg, nil, Options{DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	if results[0].KeysCurrent != 1 {
		t.Fatalf("second run results = %#v, want current entry", results)
	}
}

func TestAdoptExistingRejectsInterpolationMismatch(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	targetPath := filepath.Join(dir, "fr.json")
	if err := os.WriteFile(sourcePath, []byte(`{"welcome":"Hello, {{name}}"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(targetPath, []byte(`{"welcome":"Bonjour"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := testConfig(dir, sourcePath)

	results, err := Run(context.Background(), cfg, nil, Options{AdoptExisting: true})
	if err == nil {
		t.Fatal("adoption accepted an interpolation mismatch")
	}
	if len(results) != 1 || len(results[0].Errors) == 0 {
		t.Fatalf("results = %#v, want adoption error", results)
	}
	if _, statErr := os.Stat(cfg.ManifestPath); !os.IsNotExist(statErr) {
		t.Fatalf("manifest was written after invalid adoption; stat error = %v", statErr)
	}
}

func TestRunRetranslatesSourceStaleEntry(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	targetPath := filepath.Join(dir, "fr.json")
	if err := os.WriteFile(sourcePath, []byte(`{"a":"Old"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(targetPath, []byte(`{"a":"Ancien"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := testConfig(dir, sourcePath)

	if _, err := Run(context.Background(), cfg, nil, Options{AdoptExisting: true}); err != nil {
		t.Fatalf("adopting existing translation: %v", err)
	}
	if err := os.WriteFile(sourcePath, []byte(`{"a":"New"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	provider := &fakeProvider{response: &llm.TranslateResponse{
		Translations: map[string]string{"a": "Nouveau"},
	}}

	results, err := Run(context.Background(), cfg, provider, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if results[0].KeysSourceStale != 1 || results[0].KeysTranslated != 1 {
		t.Fatalf("results = %#v, want one source-stale translation", results)
	}
	got, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "{\n  \"a\": \"Nouveau\"\n}\n" {
		t.Fatalf("target = %q", got)
	}
}

func TestRunProtectsManualEditEvenWhenSourceChanges(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	targetPath := filepath.Join(dir, "fr.json")
	if err := os.WriteFile(sourcePath, []byte(`{"a":"Old"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(targetPath, []byte(`{"a":"Ancien"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := testConfig(dir, sourcePath)
	if _, err := Run(context.Background(), cfg, nil, Options{AdoptExisting: true}); err != nil {
		t.Fatal(err)
	}
	manual := []byte(`{"a":"Révision humaine"}`)
	if err := os.WriteFile(targetPath, manual, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(sourcePath, []byte(`{"a":"New"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	provider := &fakeProvider{response: &llm.TranslateResponse{
		Translations: map[string]string{"a": "Nouveau"},
	}}

	results, err := Run(context.Background(), cfg, provider, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if results[0].KeysManualEdit != 1 || provider.calls != 0 {
		t.Fatalf("results = %#v, provider calls = %d", results, provider.calls)
	}
	got, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(manual) {
		t.Fatalf("manual target changed: got %q, want %q", got, manual)
	}
}

func TestRunReportsPolicyStaleUntilRefreshIsExplicit(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	targetPath := filepath.Join(dir, "fr.json")
	if err := os.WriteFile(sourcePath, []byte(`{"a":"Save"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(targetPath, []byte(`{"a":"Enregistrer"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := testConfig(dir, sourcePath)
	if _, err := Run(context.Background(), cfg, nil, Options{AdoptExisting: true}); err != nil {
		t.Fatal(err)
	}
	cfg.LLM.Model = "test-v2"
	provider := &fakeProvider{response: &llm.TranslateResponse{
		Translations: map[string]string{"a": "Sauvegarder"},
	}}

	results, err := Run(context.Background(), cfg, provider, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if results[0].KeysPolicyStale != 1 || provider.calls != 0 {
		t.Fatalf("results = %#v, provider calls = %d", results, provider.calls)
	}

	results, err = Run(context.Background(), cfg, provider, Options{RefreshPolicy: true})
	if err != nil {
		t.Fatal(err)
	}
	if results[0].KeysPolicyStale != 1 || results[0].KeysTranslated != 1 || provider.calls != 1 {
		t.Fatalf("refresh results = %#v, provider calls = %d", results, provider.calls)
	}
}

func testConfig(dir, sourcePath string) *config.Config {
	return &config.Config{
		SourceLocale:   "en",
		TargetLocales:  []string{"fr"},
		SourcePath:     sourcePath,
		BatchSize:      40,
		Concurrency:    1,
		StyleGuidesDir: filepath.Join(dir, "style-guides"),
		GlossaryDir:    filepath.Join(dir, "glossary"),
		TMPath:         filepath.Join(dir, ".internationalizer", "tm.jsonl"),
		ManifestPath:   filepath.Join(dir, ".internationalizer", "manifest.json"),
		LLM: config.LLM{
			Provider: "test",
			Model:    "test-v1",
		},
	}
}
