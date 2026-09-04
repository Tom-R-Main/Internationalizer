package llm

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
)

func TestBuildDocumentPromptKeepsProviderJSONContract(t *testing.T) {
	prompt := BuildDocumentPrompt("en", "fr", "", nil)
	if !strings.Contains(prompt, "_content") || !strings.Contains(prompt, "JSON object") {
		t.Fatalf("document prompt does not preserve the provider JSON contract: %q", prompt)
	}
}

func TestBuildFluentPromptDefinesSemanticRuntimeContract(t *testing.T) {
	prompt := BuildFluentPrompt("en", "fr", "", nil)
	for _, rule := range []string{"developer context", "Fluent variables", "target-locale plural variants"} {
		if !strings.Contains(prompt, rule) {
			t.Fatalf("Fluent prompt lacks %q: %s", rule, prompt)
		}
	}
}

func TestMarshalEntriesIncludesDeveloperContextWithoutChangingOutputKey(t *testing.T) {
	data, err := marshalEntries([]Entry{
		{Key: "plain", Value: "Save"},
		{Key: "guided", Value: "Open", Context: "Verb used on a button."},
	})
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded["plain"] != "Save" {
		t.Fatalf("plain entry = %#v", decoded["plain"])
	}
	guided, ok := decoded["guided"].(map[string]interface{})
	if !ok || guided["value"] != "Open" || guided["context"] != "Verb used on a button." {
		t.Fatalf("guided entry = %#v", decoded["guided"])
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
