package llm

import (
	"strings"
	"testing"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
)

func TestBuildDocumentPromptKeepsProviderJSONContract(t *testing.T) {
	prompt := BuildDocumentPrompt("en", "fr", "", nil)
	if !strings.Contains(prompt, "every original input key") || !strings.Contains(prompt, "JSON object") {
		t.Fatalf("document prompt does not preserve the provider JSON contract: %q", prompt)
	}
}

func TestEffectiveReasoningEffort(t *testing.T) {
	tests := []struct {
		name string
		cfg  config.LLM
		want string
	}{
		{name: "official GPT-5.6 default", cfg: config.LLM{Provider: "openai", Model: "gpt-5.6-luna"}, want: "max"},
		{name: "official GPT-5 explicit", cfg: config.LLM{Provider: "openai", Model: "gpt-5", ReasoningEffort: "low"}, want: "low"},
		{name: "non-GPT-5", cfg: config.LLM{Provider: "openai", Model: "gpt-4.1", ReasoningEffort: "low"}},
		{name: "compatible endpoint", cfg: config.LLM{Provider: "openai", Model: "gpt-5", BaseURL: "https://example.com/v1", ReasoningEffort: "low"}},
		{name: "OpenRouter", cfg: config.LLM{Provider: "openrouter", Model: "gpt-5", ReasoningEffort: "low"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := EffectiveReasoningEffort(test.cfg); got != test.want {
				t.Fatalf("EffectiveReasoningEffort() = %q, want %q", got, test.want)
			}
		})
	}
}
