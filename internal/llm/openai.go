package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"
)

// OpenAI implements the Provider interface using the OpenAI Responses API.
// Also used for OpenRouter and other compatible endpoints via base_url.
type OpenAI struct {
	apiKey  string
	model   string
	baseURL string
	client  *http.Client
	// useCompat forces Chat Completions API for non-OpenAI endpoints (e.g. OpenRouter).
	useCompat bool
}

func NewOpenAI(apiKey, model, baseURL string) *OpenAI {
	if model == "" {
		model = "gpt-5.4"
	}
	if baseURL == "" {
		baseURL = "https://api.openai.com"
	}
	baseURL = strings.TrimRight(baseURL, "/")
	// Use Chat Completions for non-OpenAI endpoints (OpenRouter, Ollama, etc.)
	useCompat := !strings.Contains(baseURL, "api.openai.com")
	return &OpenAI{
		apiKey:    apiKey,
		model:     model,
		baseURL:   baseURL,
		client:    &http.Client{Timeout: 120 * time.Second},
		useCompat: useCompat,
	}
}

func (o *OpenAI) Name() string { return "openai" }

func (o *OpenAI) Translate(ctx context.Context, req TranslateRequest) (*TranslateResponse, error) {
	input := make(map[string]string, len(req.Entries))
	for _, e := range req.Entries {
		input[e.Key] = e.Value
	}
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("marshal input: %w", err)
	}

	temp := req.Temperature
	if temp == 0 {
		temp = 0.1
	}

	var body map[string]interface{}
	if o.useCompat {
		body = o.buildChatCompletionsBody(req.SystemPrompt, string(inputJSON), temp)
	} else {
		body = o.buildResponsesBody(req.SystemPrompt, string(inputJSON), temp)
	}

	var resp *TranslateResponse
	err = o.withRetry(ctx, func() error {
		var retryErr error
		resp, retryErr = o.doRequest(ctx, body)
		return retryErr
	})
	return resp, err
}

// buildResponsesBody creates a request body for the OpenAI Responses API.
func (o *OpenAI) buildResponsesBody(systemPrompt, userInput string, temp float64) map[string]interface{} {
	return map[string]interface{}{
		"model":        o.model,
		"instructions": systemPrompt,
		"input":        userInput,
		"temperature":  temp,
	}
}

// buildChatCompletionsBody creates a request body for the Chat Completions API
// (used by OpenRouter and other OpenAI-compatible endpoints).
func (o *OpenAI) buildChatCompletionsBody(systemPrompt, userInput string, temp float64) map[string]interface{} {
	return map[string]interface{}{
		"model":       o.model,
		"temperature": temp,
		"response_format": map[string]string{
			"type": "json_object",
		},
		"messages": []map[string]interface{}{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userInput},
		},
	}
}

func (o *OpenAI) doRequest(ctx context.Context, body map[string]interface{}) (*TranslateResponse, error) {
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	endpoint := "/v1/responses"
	if o.useCompat {
		endpoint = "/v1/chat/completions"
	}
	url := o.baseURL + endpoint

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(bodyJSON))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+o.apiKey)

	httpResp, err := o.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("openai request: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if httpResp.StatusCode != 200 {
		return nil, &APIError{
			StatusCode: httpResp.StatusCode,
			Body:       string(respBody),
		}
	}

	if o.useCompat {
		return o.parseChatCompletionsResponse(respBody)
	}
	return o.parseResponsesAPIResponse(respBody)
}

// parseResponsesAPIResponse parses the OpenAI Responses API format.
func (o *OpenAI) parseResponsesAPIResponse(respBody []byte) (*TranslateResponse, error) {
	var result struct {
		Output []struct {
			Type    string `json:"type"`
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"output"`
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	var text string
	for _, item := range result.Output {
		if item.Type == "message" {
			for _, c := range item.Content {
				if c.Type == "output_text" {
					text = c.Text
					break
				}
			}
			if text != "" {
				break
			}
		}
	}
	if text == "" {
		return nil, fmt.Errorf("no text output in response")
	}

	translations, err := ParseTranslationResponse(text)
	if err != nil {
		return nil, fmt.Errorf("parsing translations: %w", err)
	}

	return &TranslateResponse{
		Translations: translations,
		Usage: TokenUsage{
			InputTokens:  result.Usage.InputTokens,
			OutputTokens: result.Usage.OutputTokens,
		},
	}, nil
}

// parseChatCompletionsResponse parses the Chat Completions API format.
func (o *OpenAI) parseChatCompletionsResponse(respBody []byte) (*TranslateResponse, error) {
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	translations, err := ParseTranslationResponse(result.Choices[0].Message.Content)
	if err != nil {
		return nil, fmt.Errorf("parsing translations: %w", err)
	}

	return &TranslateResponse{
		Translations: translations,
		Usage: TokenUsage{
			InputTokens:  result.Usage.PromptTokens,
			OutputTokens: result.Usage.CompletionTokens,
		},
	}, nil
}

func (o *OpenAI) withRetry(ctx context.Context, fn func() error) error {
	const maxRetries = 3
	for attempt := range maxRetries {
		err := fn()
		if err == nil {
			return nil
		}
		apiErr, ok := err.(*APIError)
		if !ok || !apiErr.Retryable() {
			return err
		}
		if attempt == maxRetries-1 {
			return err
		}
		backoff := time.Duration(math.Pow(2, float64(attempt))) * time.Second
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
		}
	}
	return nil
}
