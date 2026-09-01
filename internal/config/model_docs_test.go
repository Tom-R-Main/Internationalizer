package config

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"
)

func TestREADMEUsesCurrentProviderDefaults(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locating config test source")
	}
	repositoryRoot := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", ".."))
	readme, err := os.ReadFile(filepath.Join(repositoryRoot, "README.md"))
	if err != nil {
		t.Fatal(err)
	}

	providers := map[string]string{
		"anthropic":  DefaultAnthropicModel,
		"openai":     DefaultOpenAIModel,
		"gemini":     DefaultGeminiModel,
		"openrouter": DefaultOpenRouterModel,
	}
	for provider, want := range providers {
		pattern := regexp.MustCompile(`(?m)#\s+` + regexp.QuoteMeta(provider) + `:\s+(\S+)`)
		match := pattern.FindSubmatch(readme)
		if len(match) != 2 {
			t.Errorf("README does not document the %s default", provider)
			continue
		}
		if got := string(match[1]); got != want {
			t.Errorf("README documents %s model %q, want %q", provider, got, want)
		}
	}

	if !regexp.MustCompile(`(?m)^\s*reasoning_effort:\s*` + regexp.QuoteMeta(DefaultOpenAIReasoning) + `\s*$`).Match(readme) {
		t.Fatalf("README does not document OpenAI reasoning_effort: %s", DefaultOpenAIReasoning)
	}
}
