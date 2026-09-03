package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestApplyDefaultsUsesCurrentProviderModels(t *testing.T) {
	tests := []struct {
		provider string
		want     string
	}{
		{provider: "anthropic", want: DefaultAnthropicModel},
		{provider: "openai", want: DefaultOpenAIModel},
		{provider: "gemini", want: DefaultGeminiModel},
		{provider: "openrouter", want: DefaultOpenRouterModel},
	}

	for _, test := range tests {
		t.Run(test.provider, func(t *testing.T) {
			cfg := &Config{LLM: LLM{Provider: test.provider}}
			cfg.ApplyDefaults()
			if cfg.LLM.Model != test.want {
				t.Fatalf("model = %q, want %q", cfg.LLM.Model, test.want)
			}
			if test.provider == "openai" && cfg.LLM.ReasoningEffort != DefaultOpenAIReasoning {
				t.Fatalf("reasoning effort = %q, want %q", cfg.LLM.ReasoningEffort, DefaultOpenAIReasoning)
			}
		})
	}
}

func TestLLMForLocaleUsesProviderSpecificOverride(t *testing.T) {
	cfg := &Config{
		LLM: LLM{
			Provider:  "gemini",
			Model:     "gemini-3.8-flash",
			APIKeyEnv: "GOOGLE_AI_STUDIO_API_KEY",
			LocaleOverrides: map[string]LLMOverride{
				"ja": {
					Provider: "openrouter",
					Model:    "sakana/sakana-namazu",
				},
			},
		},
	}
	cfg.ApplyDefaults()

	if got := cfg.LLMForLocale("fr"); got.Provider != "gemini" || got.Model != "gemini-3.8-flash" || got.APIKeyEnv != "GOOGLE_AI_STUDIO_API_KEY" {
		t.Fatalf("French LLM = %#v, want global Gemini settings", got)
	}
	if got := cfg.LLMForLocale("ja"); got.Provider != "openrouter" || got.Model != "sakana/sakana-namazu" || got.APIKeyEnv != "OPENROUTER_API_KEY" {
		t.Fatalf("Japanese LLM = %#v, want Sakana via OpenRouter", got)
	}
}

func TestLLMForLocaleInheritsGlobalProviderSettings(t *testing.T) {
	cfg := &Config{
		LLM: LLM{
			Provider:        "openai",
			Model:           "gpt-5.6-luna",
			APIKeyEnv:       "CUSTOM_OPENAI_KEY",
			BaseURL:         "https://example.com/v1",
			ReasoningEffort: "max",
			LocaleOverrides: map[string]LLMOverride{
				"ja": {Model: "gpt-5.6-luna-japanese"},
			},
		},
	}

	got := cfg.LLMForLocale("ja")
	if got.Provider != "openai" || got.Model != "gpt-5.6-luna-japanese" || got.APIKeyEnv != "CUSTOM_OPENAI_KEY" || got.BaseURL != "https://example.com/v1" || got.ReasoningEffort != "max" {
		t.Fatalf("Japanese LLM = %#v, want inherited OpenAI settings", got)
	}
}

func TestValidateCredentialsForLocalesChecksOnlySelectedProviders(t *testing.T) {
	const globalKey = "TEST_INTERNATIONALIZER_GLOBAL_KEY"
	const japaneseKey = "TEST_INTERNATIONALIZER_JAPANESE_KEY"
	t.Setenv(globalKey, "")
	t.Setenv(japaneseKey, "available")

	cfg := &Config{
		TargetLocales: []string{"fr", "ja"},
		LLM: LLM{
			Provider:  "gemini",
			APIKeyEnv: globalKey,
			LocaleOverrides: map[string]LLMOverride{
				"ja": {Provider: "openrouter", APIKeyEnv: japaneseKey},
			},
		},
	}
	if err := cfg.ValidateCredentialsForLocales([]string{"ja"}); err != nil {
		t.Fatalf("validating Japanese credentials: %v", err)
	}
	if err := cfg.ValidateCredentialsForLocales([]string{"fr"}); err == nil {
		t.Fatalf("validation accepted missing credential %s", globalKey)
	}

	if got := os.Getenv(japaneseKey); got != "available" {
		t.Fatalf("test credential changed to %q", got)
	}
}

func TestLoadResolvesLocaleOverrideFromYAML(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".internationalizer.yml")
	data := []byte(`source_locale: en
target_locales: [fr, ja]
source_path: locales/en.json
llm:
  provider: gemini
  locale_overrides:
    ja:
      provider: openrouter
      model: sakana/sakana-namazu
`)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	got := cfg.LLMForLocale("ja")
	if got.Provider != "openrouter" || got.Model != "sakana/sakana-namazu" || got.APIKeyEnv != "OPENROUTER_API_KEY" {
		t.Fatalf("Japanese LLM = %#v, want YAML OpenRouter override", got)
	}
}

func TestEffectiveBundlesPreservesLegacySourcePathContract(t *testing.T) {
	cfg := &Config{SourcePath: filepath.Join("locales", "en.json")}
	bundles := cfg.EffectiveBundles()
	if len(bundles) != 1 {
		t.Fatalf("EffectiveBundles returned %d bundles, want 1", len(bundles))
	}
	if bundles[0].ID != "default" || bundles[0].Source != cfg.SourcePath {
		t.Fatalf("legacy bundle = %#v", bundles[0])
	}
	target, err := bundles[0].TargetPath("pt-BR")
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join("locales", "pt-BR.json"); target != want {
		t.Fatalf("target = %q, want %q", target, want)
	}
}

func TestBundleTargetPathUsesConfiguredTemplate(t *testing.T) {
	bundle := Bundle{ID: "docs", Source: "README.md", Target: "docs/i18n/{locale}.md"}
	target, err := bundle.TargetPath("fr")
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join("docs", "i18n", "fr.md"); target != want {
		t.Fatalf("target = %q, want %q", target, want)
	}
}

func TestValidateProjectRejectsTargetWithoutLocalePlaceholder(t *testing.T) {
	cfg := &Config{
		TargetLocales: []string{"fr"},
		Bundles: []Bundle{{
			ID:     "docs",
			Source: "README.md",
			Target: "docs/i18n/fr.md",
		}},
	}
	if err := cfg.ValidateProject(); err == nil {
		t.Fatal("ValidateProject accepted a fixed target path")
	}
}

func TestValidateProjectRequiresStableExplicitBundleID(t *testing.T) {
	cfg := &Config{
		TargetLocales: []string{"fr"},
		Bundles: []Bundle{{
			Source: "locales/en.json",
			Target: "locales/{locale}.json",
		}},
	}
	if err := cfg.ValidateProject(); err == nil {
		t.Fatal("ValidateProject accepted an explicit bundle without a stable id")
	}
}

func TestValidateProjectRejectsLLMOverrideForUnknownLocale(t *testing.T) {
	cfg := &Config{
		TargetLocales: []string{"fr"},
		SourcePath:    filepath.Join("locales", "en.json"),
		LLM: LLM{LocaleOverrides: map[string]LLMOverride{
			"ja": {Provider: "openrouter"},
		}},
	}
	if err := cfg.ValidateProject(); err == nil {
		t.Fatal("ValidateProject accepted an LLM override for an unconfigured locale")
	}
}
