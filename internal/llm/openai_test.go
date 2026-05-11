package llm

import "testing"

func TestBuildResponsesBody_GPT5UsesReasoningAndOmitsTemperature(t *testing.T) {
	provider := NewOpenAI("test-key", "gpt-5.5", "https://api.openai.com", "low")

	body := provider.buildResponsesBody("system", `{"hello":"Hello"}`, 0.1)

	if body["model"] != "gpt-5.5" {
		t.Fatalf("model = %v, want gpt-5.5", body["model"])
	}
	if _, ok := body["temperature"]; ok {
		t.Fatal("gpt-5 Responses body must not include temperature")
	}
	reasoning, ok := body["reasoning"].(map[string]string)
	if !ok {
		t.Fatalf("reasoning = %#v, want map[string]string", body["reasoning"])
	}
	if reasoning["effort"] != "low" {
		t.Fatalf("reasoning.effort = %q, want low", reasoning["effort"])
	}
}

func TestBuildResponsesBody_NonGPT5KeepsTemperature(t *testing.T) {
	provider := NewOpenAI("test-key", "gpt-4.1", "https://api.openai.com", "")

	body := provider.buildResponsesBody("system", `{"hello":"Hello"}`, 0.2)

	if body["temperature"] != 0.2 {
		t.Fatalf("temperature = %v, want 0.2", body["temperature"])
	}
	if _, ok := body["reasoning"]; ok {
		t.Fatal("non-GPT-5 body should not include reasoning by default")
	}
}

func TestBuildChatCompletionsBody_CompatKeepsTemperatureForGPT5Names(t *testing.T) {
	provider := NewOpenAI("test-key", "gpt-5.5", "https://openrouter.ai/api/v1", "low")

	body := provider.buildChatCompletionsBody("system", `{"hello":"Hello"}`, 0.3)

	if body["temperature"] != 0.3 {
		t.Fatalf("temperature = %v, want 0.3", body["temperature"])
	}
}
