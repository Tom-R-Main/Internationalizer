package llm

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func TestProviderConstructorsUseCurrentDefaults(t *testing.T) {
	if got := NewAnthropic("key", "").model; got != config.DefaultAnthropicModel {
		t.Fatalf("Anthropic model = %q, want %q", got, config.DefaultAnthropicModel)
	}
	openAI := NewOpenAI("key", "", "", "")
	if got := openAI.model; got != config.DefaultOpenAIModel {
		t.Fatalf("OpenAI model = %q, want %q", got, config.DefaultOpenAIModel)
	}
	if got := openAI.reasoningEffort; got != config.DefaultOpenAIReasoning {
		t.Fatalf("OpenAI reasoning effort = %q, want %q", got, config.DefaultOpenAIReasoning)
	}
	if got := NewGemini("key", "").model; got != config.DefaultGeminiModel {
		t.Fatalf("Gemini model = %q, want %q", got, config.DefaultGeminiModel)
	}
	openRouter := NewOpenRouter("key", "")
	if openRouter.model != config.DefaultOpenRouterModel {
		t.Fatalf("OpenRouter model = %q, want %q", openRouter.model, config.DefaultOpenRouterModel)
	}
	if openRouter.Name() != "openrouter" {
		t.Fatalf("OpenRouter provider name = %q, want openrouter", openRouter.Name())
	}
}

func TestOpenRouterUsesCanonicalChatCompletionsEndpoint(t *testing.T) {
	provider := NewOpenRouter("key", "")
	provider.client.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if got, want := req.URL.String(), "https://openrouter.ai/api/v1/chat/completions"; got != want {
			t.Fatalf("request URL = %q, want %q", got, want)
		}
		return jsonResponse(`{"choices":[{"message":{"content":"{\"hello\":\"bonjour\"}"}}]}`), nil
	})

	response, err := provider.Translate(context.Background(), TranslateRequest{
		Entries:      []Entry{{Key: "hello", Value: "hello"}},
		SystemPrompt: "translate",
	})
	if err != nil {
		t.Fatal(err)
	}
	if response.Translations["hello"] != "bonjour" {
		t.Fatalf("translation = %q, want bonjour", response.Translations["hello"])
	}
}

func TestAnthropicOpusRequestOmitsDeprecatedSamplingParameters(t *testing.T) {
	provider := NewAnthropic("key", "")
	provider.client.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		var body map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["model"] != config.DefaultAnthropicModel {
			t.Fatalf("model = %v, want %s", body["model"], config.DefaultAnthropicModel)
		}
		if _, ok := body["temperature"]; ok {
			t.Fatal("Claude Opus 5 request must not include temperature")
		}
		return jsonResponse(`{"content":[{"type":"text","text":"{\"hello\":\"bonjour\"}"}]}`), nil
	})

	if _, err := provider.Translate(context.Background(), TranslateRequest{
		Entries:      []Entry{{Key: "hello", Value: "hello"}},
		SystemPrompt: "translate",
		Temperature:  0.1,
	}); err != nil {
		t.Fatal(err)
	}
}

func jsonResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}
