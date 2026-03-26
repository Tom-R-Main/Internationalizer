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

type Gemini struct {
	apiKey string
	model  string
	client *http.Client
}

func NewGemini(apiKey, model string) *Gemini {
	if model == "" {
		model = "gemini-3.1-pro-preview"
	}
	return &Gemini{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{Timeout: 120 * time.Second},
	}
}

func (g *Gemini) Name() string { return "gemini" }

func (g *Gemini) Translate(ctx context.Context, req TranslateRequest) (*TranslateResponse, error) {
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

	body := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"role": "user",
				"parts": []map[string]string{
					{"text": string(inputJSON)},
				},
			},
		},
		"systemInstruction": map[string]interface{}{
			"parts": []map[string]string{
				{"text": req.SystemPrompt},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature":      temp,
			"maxOutputTokens":  max(4096, len(inputJSON)*4),
			"responseMimeType": "application/json",
		},
	}

	var resp *TranslateResponse
	err = g.withRetry(ctx, func() error {
		var retryErr error
		resp, retryErr = g.doRequest(ctx, body)
		return retryErr
	})
	return resp, err
}

func (g *Gemini) doRequest(ctx context.Context, body map[string]interface{}) (*TranslateResponse, error) {
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		g.model, g.apiKey)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(bodyJSON))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := g.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("gemini request: %w", err)
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
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
		UsageMetadata struct {
			PromptTokenCount     int `json:"promptTokenCount"`
			CandidatesTokenCount int `json:"candidatesTokenCount"`
		} `json:"usageMetadata"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no content in gemini response")
	}

	text := result.Candidates[0].Content.Parts[0].Text
	translations, err := ParseTranslationResponse(text)
	if err != nil {
		return nil, fmt.Errorf("parsing translations: %w", err)
	}

	return &TranslateResponse{
		Translations: translations,
		Usage: TokenUsage{
			InputTokens:  result.UsageMetadata.PromptTokenCount,
			OutputTokens: result.UsageMetadata.CandidatesTokenCount,
		},
	}, nil
}

func (g *Gemini) withRetry(ctx context.Context, fn func() error) error {
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
