package translate

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
	"github.com/Tom-R-Main/Internationalizer/internal/llm"
	"github.com/Tom-R-Main/Internationalizer/internal/policy"
	"github.com/Tom-R-Main/Internationalizer/internal/state"
	"github.com/Tom-R-Main/Internationalizer/internal/tm"
	validation "github.com/Tom-R-Main/Internationalizer/internal/validate"
)

type fakeProvider struct {
	name     string
	response *llm.TranslateResponse
	calls    int
	requests []llm.TranslateRequest
}

func (p *fakeProvider) Name() string {
	if p.name != "" {
		return p.name
	}
	return "fake"
}

func (p *fakeProvider) Translate(_ context.Context, request llm.TranslateRequest) (*llm.TranslateResponse, error) {
	p.calls++
	p.requests = append(p.requests, request)
	return p.response, nil
}

func TestRunTranslatesFluentUnitsWithDeveloperContext(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.ftl")
	source := "# A command, not a noun.\nopen-button = Open { -brand-short-name }\n    .aria-label = Open the application\n\n-brand-short-name = Acme\n"
	if err := os.WriteFile(sourcePath, []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := testConfig(dir, sourcePath)
	provider := &fakeProvider{response: &llm.TranslateResponse{Translations: map[string]string{
		"open-button":            "Ouvrir { -brand-short-name }",
		"open-button.aria-label": "Ouvrir l’application",
		"-brand-short-name":      "Acme",
	}}}

	if _, err := Run(context.Background(), cfg, provider, Options{}); err != nil {
		t.Fatal(err)
	}
	if len(provider.requests) != 1 {
		t.Fatalf("provider requests = %#v, want one request", provider.requests)
	}
	context := ""
	for _, entry := range provider.requests[0].Entries {
		if entry.Key == "open-button" {
			context = entry.Context
		}
	}
	if len(provider.requests[0].Entries) != 3 || context != "# A command, not a noun." {
		t.Fatalf("provider requests = %#v", provider.requests)
	}
	if !strings.Contains(provider.requests[0].SystemPrompt, "Fluent variables") {
		t.Fatalf("provider did not receive Fluent prompt: %s", provider.requests[0].SystemPrompt)
	}
	output, err := os.ReadFile(filepath.Join(dir, "fr.ftl"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(output)
	for _, expected := range []string{"# A command, not a noun.", "open-button = Ouvrir { -brand-short-name }", ".aria-label = Ouvrir l’application", "-brand-short-name = Acme"} {
		if !strings.Contains(text, expected) {
			t.Fatalf("Fluent output lacks %q:\n%s", expected, text)
		}
	}
	changedContext := strings.Replace(source, "# A command, not a noun.", "# A noun shown in the File menu.", 1)
	if err := os.WriteFile(sourcePath, []byte(changedContext), 0o644); err != nil {
		t.Fatal(err)
	}
	reports, err := validation.ValidateWithOptions(cfg, validation.Options{RequireState: true})
	if err != nil {
		t.Fatal(err)
	}
	foundStale := false
	for _, finding := range reports[0].Findings {
		if finding.Key == "open-button" && finding.Code == validation.CodeSourceStale {
			foundStale = true
		}
	}
	if !foundStale {
		t.Fatalf("developer-context change did not stale provenance: %#v", reports[0].Findings)
	}
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

func TestRunRejectsInvalidICUProviderResponseWithoutWritingTarget(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	if err := os.WriteFile(sourcePath, []byte(`{"items":"{count, plural, one {One item} other {# items}}"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := testConfig(dir, sourcePath)
	provider := &fakeProvider{response: &llm.TranslateResponse{Translations: map[string]string{
		"items": "{count, plural, one {Un article} other {# articles}",
	}}}

	results, err := Run(context.Background(), cfg, provider, Options{})
	if err == nil {
		t.Fatal("Run returned nil error for malformed ICU provider output")
	}
	if len(results) != 1 || len(results[0].Errors) == 0 {
		t.Fatalf("Run results = %#v, want one locale error", results)
	}
	if _, statErr := os.Stat(filepath.Join(dir, "fr.json")); !os.IsNotExist(statErr) {
		t.Fatalf("target was written after malformed ICU output; stat error = %v", statErr)
	}
}

func TestRunRejectsMalformedICUSourceWithoutCallingProvider(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	if err := os.WriteFile(sourcePath, []byte(`{"items":"{count, plural, one {One}}"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := testConfig(dir, sourcePath)
	provider := &fakeProvider{response: &llm.TranslateResponse{Translations: map[string]string{
		"items": "{count, plural, one {Un} other {Autres}}",
	}}}

	results, err := Run(context.Background(), cfg, provider, Options{})
	if err == nil {
		t.Fatal("Run returned nil error for malformed ICU source")
	}
	if provider.calls != 0 {
		t.Fatalf("provider called %d times for malformed ICU source", provider.calls)
	}
	if len(results) != 1 || len(results[0].Errors) == 0 {
		t.Fatalf("Run results = %#v, want source error", results)
	}
}

func TestRunDoesNotReuseStructurallyInvalidTMRecord(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	source := "{count, plural, one {One item} other {# items}}"
	if err := os.WriteFile(sourcePath, []byte(`{"items":"{count, plural, one {One item} other {# items}}"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := testConfig(dir, sourcePath)
	resolved, err := policy.Resolve(cfg, "fr", "json", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	memory, err := tm.Load(cfg.TMPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := memory.Add(tm.Record{
		Bundle:     "default",
		Key:        "items",
		Source:     source,
		Target:     "{count, plural, one {Un article} other {# articles}",
		Locale:     "fr",
		Hash:       state.SourceHash("json", source),
		PolicyHash: resolved.Hash,
		Timestamp:  time.Now(),
	}); err != nil {
		t.Fatal(err)
	}
	provider := &fakeProvider{response: &llm.TranslateResponse{Translations: map[string]string{
		"items": "{count, plural, one {Un article} other {# articles}}",
	}}}

	results, err := Run(context.Background(), cfg, provider, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].KeysCached != 0 || results[0].KeysTranslated != 1 || provider.calls != 1 {
		t.Fatalf("Run results = %#v, provider calls = %d; want invalid TM bypass and fresh translation", results, provider.calls)
	}
}

func TestRunSelectsCanonicalEquivalentConfiguredLocale(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	if err := os.WriteFile(sourcePath, []byte(`{"save":"Save"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := testConfig(dir, sourcePath)
	cfg.TargetLocales = []string{"pt-br"}

	results, err := Run(context.Background(), cfg, nil, Options{DryRun: true, Locales: []string{"pt-BR"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Locale != "pt-br" || results[0].TargetPath != filepath.Join(dir, "pt-br.json") {
		t.Fatalf("Run results = %#v, want configured locale spelling and path", results)
	}
}

func TestRunProducesTargetLocalePluralFormsThatPassStrictValidation(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	if err := os.WriteFile(sourcePath, []byte(`{"items_one":"{{count}} item","items_other":"{{count}} items"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := testConfig(dir, sourcePath)
	cfg.Validation.PluralStyle = "i18next-v4"
	provider := &fakeProvider{response: &llm.TranslateResponse{Translations: map[string]string{
		"items_many":  "{{count}} articles",
		"items_one":   "{{count}} article",
		"items_other": "{{count}} articles",
	}}}

	if _, err := Run(context.Background(), cfg, provider, Options{}); err != nil {
		t.Fatal(err)
	}
	reports, err := validation.ValidateWithOptions(cfg, validation.Options{RequireState: true})
	if err != nil {
		t.Fatal(err)
	}
	if validation.HasFailures(reports) {
		t.Fatalf("translated plural target failed state validation: %#v", reports[0])
	}
	reports, err = validation.ValidateWithOptions(cfg, validation.Options{Strict: true})
	if err != nil {
		t.Fatal(err)
	}
	if validation.HasFailures(reports) {
		t.Fatalf("translated plural target failed content validation: %#v", reports[0])
	}
}

func TestRunOmitsSourceOnlyPluralFormsForTargetLocale(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	if err := os.WriteFile(sourcePath, []byte(`{"items_one":"{{count}} item","items_other":"{{count}} items"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := testConfig(dir, sourcePath)
	cfg.TargetLocales = []string{"ja"}
	cfg.Validation.PluralStyle = "i18next-v4"
	provider := &fakeProvider{response: &llm.TranslateResponse{Translations: map[string]string{
		"items_other": "{{count}} 個",
	}}}

	results, err := Run(context.Background(), cfg, provider, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].KeysTotal != 1 || results[0].KeysTranslated != 1 {
		t.Fatalf("Japanese plural result = %#v, want only items_other", results)
	}
	data, err := os.ReadFile(filepath.Join(dir, "ja.json"))
	if err != nil {
		t.Fatal(err)
	}
	var target map[string]string
	if err := json.Unmarshal(data, &target); err != nil {
		t.Fatal(err)
	}
	if target["items_other"] != "{{count}} 個" {
		t.Fatalf("Japanese target = %#v, want items_other", target)
	}
	if _, ok := target["items_one"]; ok {
		t.Fatalf("Japanese target retained source-only plural: %#v", target)
	}
}

func TestRunProducesTargetLocalePluralFormsForYAML(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.yml")
	if err := os.WriteFile(sourcePath, []byte("# source\nitems_one: '{{count}} item'\nitems_other: '{{count}} items'\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := testConfig(dir, sourcePath)
	cfg.Validation.PluralStyle = "i18next-v4"
	provider := &fakeProvider{response: &llm.TranslateResponse{Translations: map[string]string{
		"items_many":  "{{count}} articles",
		"items_one":   "{{count}} article",
		"items_other": "{{count}} articles",
	}}}

	if _, err := Run(context.Background(), cfg, provider, Options{}); err != nil {
		t.Fatal(err)
	}
	reports, err := validation.ValidateWithOptions(cfg, validation.Options{Strict: true, RequireState: true})
	if err != nil {
		t.Fatal(err)
	}
	if validation.HasFailures(reports) {
		t.Fatalf("translated YAML plural target failed validation: %#v", reports[0])
	}
}

func TestRunRoutesLocaleToConfiguredProvider(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	if err := os.WriteFile(sourcePath, []byte(`{"hello":"Hello"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := testConfig(dir, sourcePath)
	cfg.TargetLocales = []string{"fr", "ja"}
	cfg.LLM.LocaleOverrides = map[string]config.LLMOverride{
		"ja": {
			Provider: "openrouter",
			Model:    "sakana/sakana-namazu",
		},
	}
	defaultProvider := &fakeProvider{name: "gemini", response: &llm.TranslateResponse{
		Translations: map[string]string{"hello": "Bonjour"},
	}}
	japaneseProvider := &fakeProvider{name: "openrouter", response: &llm.TranslateResponse{
		Translations: map[string]string{"hello": "こんにちは"},
	}}

	results, err := Run(context.Background(), cfg, defaultProvider, Options{
		LocaleProviders: map[string]llm.Provider{"ja": japaneseProvider},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 || defaultProvider.calls != 1 || japaneseProvider.calls != 1 {
		t.Fatalf("results = %#v, default calls = %d, Japanese calls = %d", results, defaultProvider.calls, japaneseProvider.calls)
	}
	if got, err := os.ReadFile(filepath.Join(dir, "ja.json")); err != nil {
		t.Fatal(err)
	} else if string(got) != "{\n  \"hello\": \"こんにちは\"\n}\n" {
		t.Fatalf("Japanese target = %q", got)
	}
}

func TestRunDoesNotFallBackWhenLocaleProviderIsMissing(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	if err := os.WriteFile(sourcePath, []byte(`{"hello":"Hello"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := testConfig(dir, sourcePath)
	cfg.TargetLocales = []string{"ja"}
	cfg.LLM.LocaleOverrides = map[string]config.LLMOverride{
		"ja": {Provider: "openrouter", Model: "sakana/sakana-namazu"},
	}
	defaultProvider := &fakeProvider{response: &llm.TranslateResponse{
		Translations: map[string]string{"hello": "wrong provider"},
	}}

	results, err := Run(context.Background(), cfg, defaultProvider, Options{})
	if err == nil {
		t.Fatal("Run accepted a global provider for a locale with an override")
	}
	if len(results) != 1 || len(results[0].Errors) == 0 || defaultProvider.calls != 0 {
		t.Fatalf("results = %#v, global provider calls = %d", results, defaultProvider.calls)
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

func TestRunAdoptsMarkdownWithRebasedLinksAndCodeBraces(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "README.md")
	targetPath := filepath.Join(dir, "docs", "i18n", "fr.md")
	source := "# Project\n\n[License](LICENSE)\n\n## Configuration\n\n```json\n{\"target\":\"{locale}\"}\n```\n"
	target := "> [English](../../README.md)\n\n# Projet\n\n[Licence](../../LICENSE)\n\n## Configuration\n\n```json\n{\"target\":\"{locale}\"}\n```\n"
	if err := os.WriteFile(sourcePath, []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(targetPath, []byte(target), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := testConfig(dir, "")
	cfg.Bundles = []config.Bundle{{
		ID:     "docs",
		Source: sourcePath,
		Target: filepath.Join(dir, "docs", "i18n", "{locale}.md"),
		Format: "markdown",
	}}

	results, err := Run(context.Background(), cfg, nil, Options{AdoptExisting: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || len(results[0].Errors) != 0 || results[0].KeysUntracked != 2 {
		t.Fatalf("results = %#v, want two adopted Markdown units", results)
	}
}

func TestRunRetranslatesOnlyChangedMarkdownSectionAndNeverWritesGuide(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "README.md")
	targetPath := filepath.Join(dir, "docs", "i18n", "fr.md")
	guidePath := filepath.Join(dir, "style-guides", "fr.md")
	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(guidePath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(sourcePath, []byte("# Project\n\n## Install\n\nOld install.\n\n## Usage\n\nStable usage.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(targetPath, []byte("# Projet\n\n## Installation\n\nAncienne installation.\n\n## Utilisation\n\nUtilisation stable.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	guide := []byte("Use clear, professional French.\n")
	if err := os.WriteFile(guidePath, guide, 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := testConfig(dir, "")
	cfg.Bundles = []config.Bundle{{
		ID:     "docs",
		Source: sourcePath,
		Target: filepath.Join(dir, "docs", "i18n", "{locale}.md"),
		Format: "markdown",
	}}
	if _, err := Run(context.Background(), cfg, nil, Options{AdoptExisting: true}); err != nil {
		t.Fatalf("adopting existing translation: %v", err)
	}
	if err := os.WriteFile(sourcePath, []byte("# Project\n\n## Install\n\nNew install.\n\n## Usage\n\nStable usage.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	provider := &fakeProvider{response: &llm.TranslateResponse{Translations: map[string]string{
		"markdown:install": "## Installation\n\nNouvelle installation.\n",
	}}}

	results, err := Run(context.Background(), cfg, provider, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if results[0].KeysSourceStale != 1 || results[0].KeysTranslated != 1 || provider.calls != 1 {
		t.Fatalf("results = %#v, provider calls = %d", results, provider.calls)
	}
	if len(provider.requests) != 1 || len(provider.requests[0].Entries) != 1 || provider.requests[0].Entries[0].Key != "markdown:install" {
		t.Fatalf("provider requests = %#v, want only markdown:install", provider.requests)
	}
	got, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatal(err)
	}
	want := "# Projet\n\n<!-- internationalizer:unit markdown:install -->\n## Installation\n\nNouvelle installation.\n<!-- internationalizer:unit markdown:usage -->\n## Utilisation\n\nUtilisation stable.\n"
	if string(got) != want {
		t.Fatalf("target = %q, want %q", got, want)
	}
	guideAfter, err := os.ReadFile(guidePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(guideAfter) != string(guide) {
		t.Fatalf("style guide changed: got %q, want %q", guideAfter, guide)
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

func TestRunTreatsReasoningEffortAsTranslationPolicy(t *testing.T) {
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
	cfg.LLM.Provider = "openai"
	cfg.LLM.Model = "gpt-5.6-luna"
	cfg.LLM.ReasoningEffort = "low"
	if _, err := Run(context.Background(), cfg, nil, Options{AdoptExisting: true}); err != nil {
		t.Fatal(err)
	}

	cfg.LLM.ReasoningEffort = "max"
	results, err := Run(context.Background(), cfg, nil, Options{DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	if results[0].KeysPolicyStale != 1 || results[0].KeysCurrent != 0 {
		t.Fatalf("reasoning policy change was not reported as stale: %#v", results[0])
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
