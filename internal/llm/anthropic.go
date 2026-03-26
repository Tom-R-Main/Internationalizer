package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
)

type Anthropic struct {
	apiKey string
	model  string
	client *http.Client
}

func NewAnthropic(apiKey, model string) *Anthropic {
	if model == "" {
		model = "claude-sonnet-4-6"
	}
	return &Anthropic{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{Timeout: 120 * time.Second},
	}
}

func (a *Anthropic) Name() string { return "anthropic" }

func (a *Anthropic) Translate(ctx context.Context, req TranslateRequest) (*TranslateResponse, error) {
	// Build the user message as a JSON object of keys to translate.
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
	maxTokens := max(4096, len(inputJSON)*4)

	body := map[string]interface{}{
		"model":      a.model,
		"max_tokens": maxTokens,
		"temperature": temp,
		"system":     req.SystemPrompt,
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": string(inputJSON),
			},
		},
	}

	var resp *TranslateResponse
	err = a.withRetry(ctx, func() error {
		var retryErr error
		resp, retryErr = a.doRequest(ctx, body)
		return retryErr
	})
	return resp, err
}

func (a *Anthropic) doRequest(ctx context.Context, body map[string]interface{}) (*TranslateResponse, error) {
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewReader(bodyJSON))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", a.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	httpResp, err := a.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("anthropic request: %w", err)
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

	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	var text string
	for _, c := range result.Content {
		if c.Type == "text" {
			text = c.Text
			break
		}
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

func (a *Anthropic) withRetry(ctx context.Context, fn func() error) error {
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

// APIError represents an HTTP error from an LLM API.
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Body)
}

func (e *APIError) Retryable() bool {
	return e.StatusCode == 429 || e.StatusCode == 529 || e.StatusCode >= 500
}
