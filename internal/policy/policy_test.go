package policy_test

import (
	"reflect"
	"testing"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
	"github.com/Tom-R-Main/Internationalizer/internal/glossary"
	"github.com/Tom-R-Main/Internationalizer/internal/policy"
)

func TestResolveIsDeterministic(t *testing.T) {
	cfg := baseConfig()
	terms := []glossary.Term{{Source: "workspace", Target: "espace de travail", WholeWord: true}}

	first, err := policy.Resolve(cfg, "fr", "json", "Use formal language.", terms)
	if err != nil {
		t.Fatal(err)
	}
	second, err := policy.Resolve(cfg, "fr", "json", "Use formal language.", terms)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("resolutions differ: first = %#v, second = %#v", first, second)
	}
	if want := cfg.LLMForLocale("fr"); !reflect.DeepEqual(first.LLM, want) {
		t.Fatalf("effective LLM = %#v, want %#v", first.LLM, want)
	}
	const wantHash = "fdfbfa87db81fd2fe792922f343342fe89705e51bb8bee61534a6813f17bf3d1"
	if first.Hash != wantHash {
		t.Fatalf("hash = %q, want legacy-compatible %q", first.Hash, wantHash)
	}
}

func TestResolveHashSensitivity(t *testing.T) {
	base, err := policy.Resolve(baseConfig(), "fr", "json", "Use formal language.", nil)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name         string
		mutateConfig func(*config.Config)
		targetLocale string
		format       string
		styleGuide   string
		terms        []glossary.Term
	}{
		{
			name:       "prompt",
			styleGuide: "Use informal language.",
		},
		{
			name:  "glossary prompt",
			terms: []glossary.Term{{Source: "workspace", Target: "espace de travail"}},
		},
		{
			name: "provider",
			mutateConfig: func(cfg *config.Config) {
				cfg.LLM.Provider = "anthropic"
			},
		},
		{
			name: "model",
			mutateConfig: func(cfg *config.Config) {
				cfg.LLM.Model = "gpt-5.6-terra"
			},
		},
		{
			name: "reasoning",
			mutateConfig: func(cfg *config.Config) {
				cfg.LLM.ReasoningEffort = "max"
			},
		},
		{
			name:   "format",
			format: "yaml",
		},
		{
			name:         "target locale",
			targetLocale: "de",
		},
		{
			name: "source locale",
			mutateConfig: func(cfg *config.Config) {
				cfg.SourceLocale = "en-GB"
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := baseConfig()
			if test.mutateConfig != nil {
				test.mutateConfig(cfg)
			}
			targetLocale := test.targetLocale
			if targetLocale == "" {
				targetLocale = "fr"
			}
			format := test.format
			if format == "" {
				format = "json"
			}
			styleGuide := test.styleGuide
			if styleGuide == "" {
				styleGuide = "Use formal language."
			}

			got, err := policy.Resolve(cfg, targetLocale, format, styleGuide, test.terms)
			if err != nil {
				t.Fatal(err)
			}
			if got.Hash == base.Hash {
				t.Fatalf("hash did not change from %q", base.Hash)
			}
		})
	}
}

func TestResolveUsesLocaleLLMOverride(t *testing.T) {
	cfg := baseConfig()
	cfg.LLM.LocaleOverrides = map[string]config.LLMOverride{
		"ja": {
			Provider:  "gemini",
			Model:     "gemini-3.8-flash",
			APIKeyEnv: "GOOGLE_AI_STUDIO_API_KEY",
		},
	}

	resolved, err := policy.Resolve(cfg, "ja", "json", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	want := cfg.LLMForLocale("ja")
	if !reflect.DeepEqual(resolved.LLM, want) {
		t.Fatalf("effective LLM = %#v, want %#v", resolved.LLM, want)
	}
}

func baseConfig() *config.Config {
	return &config.Config{
		SourceLocale: "en",
		LLM: config.LLM{
			Provider:        "openai",
			Model:           "gpt-5.6-luna",
			ReasoningEffort: "low",
		},
	}
}
