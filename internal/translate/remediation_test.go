package translate

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Tom-R-Main/Internationalizer/internal/llm"
	"github.com/Tom-R-Main/Internationalizer/internal/state"
	"gopkg.in/yaml.v3"
)

func TestRunNewJSONTargetPreservesNonStringSourceSchema(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	if err := os.WriteFile(sourcePath, []byte(`{"label":"Save","enabled":true,"limit":3,"empty":null}`), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := testConfig(dir, sourcePath)
	provider := &fakeProvider{response: &llm.TranslateResponse{Translations: map[string]string{"label": "Enregistrer"}}}

	if _, err := Run(context.Background(), cfg, provider, Options{}); err != nil {
		t.Fatal(err)
	}
	output, err := os.ReadFile(filepath.Join(dir, "fr.json"))
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]interface{}
	if err := json.Unmarshal(output, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded["enabled"] != true || decoded["limit"] != float64(3) {
		t.Fatalf("new target lost non-string source fields: %#v", decoded)
	}
	if value, exists := decoded["empty"]; !exists || value != nil {
		t.Fatalf("new target lost null source field: %#v", decoded)
	}
}

func TestRunNewYAMLTargetPreservesSequencesAndNonStringFields(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.yaml")
	if err := os.WriteFile(sourcePath, []byte("steps:\n  - First\n  - Second\nenabled: true\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := testConfig(dir, sourcePath)
	provider := &fakeProvider{response: &llm.TranslateResponse{Translations: map[string]string{
		"steps.0": "Premier",
		"steps.1": "Deuxième",
	}}}

	if _, err := Run(context.Background(), cfg, provider, Options{}); err != nil {
		t.Fatal(err)
	}
	output, err := os.ReadFile(filepath.Join(dir, "fr.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	var decoded struct {
		Steps   []string `yaml:"steps"`
		Enabled bool     `yaml:"enabled"`
	}
	if err := yaml.Unmarshal(output, &decoded); err != nil {
		t.Fatal(err)
	}
	if len(decoded.Steps) != 2 || decoded.Steps[1] != "Deuxième" || !decoded.Enabled {
		t.Fatalf("new target did not preserve source schema: %#v; output = %s", decoded, output)
	}
}

func TestRunRejectsBlankProviderTranslation(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	if err := os.WriteFile(sourcePath, []byte(`{"label":"Save"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := testConfig(dir, sourcePath)
	provider := &fakeProvider{response: &llm.TranslateResponse{Translations: map[string]string{"label": "   "}}}

	if _, err := Run(context.Background(), cfg, provider, Options{}); err == nil {
		t.Fatal("Run accepted a blank translation for nonblank source text")
	}
	if _, err := os.Stat(filepath.Join(dir, "fr.json")); !os.IsNotExist(err) {
		t.Fatalf("target was written after invalid provider output; stat error = %v", err)
	}
}

func TestRunRejectsProviderTranslationThatDropsProtectedStructure(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	if err := os.WriteFile(sourcePath, []byte(`{"label":"<strong>Save</strong>"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := testConfig(dir, sourcePath)
	provider := &fakeProvider{response: &llm.TranslateResponse{Translations: map[string]string{"label": "Enregistrer"}}}

	if _, err := Run(context.Background(), cfg, provider, Options{}); err == nil {
		t.Fatal("Run accepted a translation that dropped HTML tags")
	}
}

func TestAdoptExistingReconcilesReviewedManualEditAndReportsStaleness(t *testing.T) {
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
	if err := os.WriteFile(sourcePath, []byte(`{"a":"New"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(targetPath, []byte(`{"a":"Révision humaine"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	results, err := Run(context.Background(), cfg, nil, Options{DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	if results[0].KeysManualEdit != 1 || results[0].KeysSourceStale != 1 {
		t.Fatalf("manual source-stale edit was not reported independently: %#v", results[0])
	}
	if _, err := Run(context.Background(), cfg, nil, Options{AdoptExisting: true}); err != nil {
		t.Fatalf("adopting reviewed manual edit: %v", err)
	}
	results, err = Run(context.Background(), cfg, nil, Options{DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	if results[0].KeysCurrent != 1 || results[0].KeysManualEdit != 0 || results[0].KeysSourceStale != 0 {
		t.Fatalf("adopted manual edit did not become current: %#v", results[0])
	}
}

func TestRunRecordsConcreteProviderProvenance(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	if err := os.WriteFile(sourcePath, []byte(`{"a":"Save"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := testConfig(dir, sourcePath)
	provider := &fakeProvider{response: &llm.TranslateResponse{Translations: map[string]string{"a": "Enregistrer"}}}

	if _, err := Run(context.Background(), cfg, provider, Options{}); err != nil {
		t.Fatal(err)
	}
	manifest, err := state.Load(cfg.ManifestPath)
	if err != nil {
		t.Fatal(err)
	}
	entry, ok := manifest.Get("default", "a", "fr")
	if !ok {
		t.Fatal("translation was not recorded in the manifest")
	}
	if entry.Provider != provider.Name() {
		t.Fatalf("manifest provider = %q, want %q", entry.Provider, provider.Name())
	}
	if entry.Origin != "provider" {
		t.Fatalf("manifest origin = %q, want provider", entry.Origin)
	}

	if err := os.Remove(filepath.Join(dir, "fr.json")); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(cfg.ManifestPath); err != nil {
		t.Fatal(err)
	}
	results, err := Run(context.Background(), cfg, provider, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if results[0].KeysCached != 1 || provider.calls != 1 {
		t.Fatalf("translation memory was not reused: result = %#v, provider calls = %d", results[0], provider.calls)
	}
	manifest, err = state.Load(cfg.ManifestPath)
	if err != nil {
		t.Fatal(err)
	}
	entry, ok = manifest.Get("default", "a", "fr")
	if !ok {
		t.Fatal("cached translation was not recorded in the manifest")
	}
	if entry.Origin != "tm" || entry.Provider != provider.Name() || entry.Model != cfg.LLM.Model {
		t.Fatalf("cached provenance = %#v", entry)
	}
}

func TestFormatResultsSeparatesObservedStateFromCompletedWork(t *testing.T) {
	output := FormatResults([]Result{{
		Bundle:         "app",
		Locale:         "fr",
		TargetPath:     "fr.json",
		KeysMissing:    1,
		KeysTranslated: 1,
	}}, 0)
	if !strings.Contains(output, "Observed before run:") {
		t.Fatalf("output does not distinguish observed pre-run state: %q", output)
	}
	if strings.Contains(output, "\nState:") || strings.Contains(output, "\nPlan:") {
		t.Fatalf("output presents pre-run observations as post-run state or plan: %q", output)
	}
}
