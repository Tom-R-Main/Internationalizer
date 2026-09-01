package config

import (
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
