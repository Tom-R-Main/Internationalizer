package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
	"github.com/Tom-R-Main/Internationalizer/internal/glossary"
)

// Entry is a source key-value pair to be translated.
type Entry struct {
	Key     string
	Value   string
	Context string
}

// TranslateRequest is the input to a translation call.
type TranslateRequest struct {
	SourceLocale string
	TargetLocale string
	Entries      []Entry
	SystemPrompt string
	Temperature  float64
}

// TranslateResponse is the result of a translation call.
type TranslateResponse struct {
	Translations map[string]string
	Usage        TokenUsage
}

// TokenUsage tracks LLM token consumption.
type TokenUsage struct {
	InputTokens  int
	OutputTokens int
}

// Provider is the interface all LLM backends implement.
type Provider interface {
	Name() string
	Translate(ctx context.Context, req TranslateRequest) (*TranslateResponse, error)
}

// NewProvider creates an LLM provider from config.
func NewProvider(cfg config.LLM, apiKey string) (Provider, error) {
	switch cfg.Provider {
	case "anthropic":
		return NewAnthropic(apiKey, cfg.Model), nil
	case "openai":
		baseURL := cfg.BaseURL
		if baseURL == "" {
			baseURL = "https://api.openai.com"
		}
		return NewOpenAI(apiKey, cfg.Model, baseURL, EffectiveReasoningEffort(cfg)), nil
	case "gemini":
		return NewGemini(apiKey, cfg.Model), nil
	case "openrouter":
		return NewOpenRouter(apiKey, cfg.Model), nil
	default:
		return nil, fmt.Errorf("unknown LLM provider: %s (supported: anthropic, openai, gemini, openrouter)", cfg.Provider)
	}
}

// EffectiveReasoningEffort returns the reasoning setting that reaches the
// provider request. Compatible endpoints and non-GPT-5 models ignore it.
func EffectiveReasoningEffort(cfg config.LLM) string {
	if cfg.Provider != "openai" || !strings.HasPrefix(cfg.Model, "gpt-5") {
		return ""
	}
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com"
	}
	if !strings.Contains(strings.TrimRight(baseURL, "/"), "api.openai.com") {
		return ""
	}
	if cfg.ReasoningEffort == "" && strings.HasPrefix(cfg.Model, "gpt-5.6") {
		return config.DefaultOpenAIReasoning
	}
	return cfg.ReasoningEffort
}

// BuildSystemPrompt constructs the system prompt for translation,
// combining rules, glossary terms, and style guide content.
func BuildSystemPrompt(sourceLocale, targetLocale, styleGuide string, terms []glossary.Term) string {
	var b strings.Builder

	b.WriteString("You are a professional software localizer. Translate the following ")
	b.WriteString("key-value pairs from ")
	b.WriteString(sourceLocale)
	b.WriteString(" to ")
	b.WriteString(targetLocale)
	b.WriteString(".\n\n")

	b.WriteString("## Rules\n")
	b.WriteString("- Return a JSON object mapping each key to its translated value.\n")
	b.WriteString("- Preserve interpolation variables exactly: {{variable}}, {variable}, %{variable}.\n")
	b.WriteString("- Preserve all HTML tags exactly (<strong>, <br/>, etc.).\n")
	b.WriteString("- Do not translate the JSON keys, only the values.\n")
	b.WriteString("- Do not add any explanation or commentary, only output valid JSON.\n")
	b.WriteString("- Keep brand names and technical terms in English unless the glossary specifies otherwise.\n")
	b.WriteString("- Follow the target language's CLDR plural rules when translating plural forms.\n")
	b.WriteString("- Match the tone and formality level described in the style guide.\n")
	b.WriteString("- Keep translations concise — UI strings should not be more than 150% the length of the source.\n")

	if len(terms) > 0 {
		b.WriteString("\n## Glossary\n")
		b.WriteString("Use these approved translations for the following terms:\n\n")
		b.WriteString("| Source | Translation | Notes |\n")
		b.WriteString("|--------|-------------|-------|\n")
		for _, t := range terms {
			notes := ""
			if t.WholeWord {
				notes = "whole word match"
			}
			_, _ = fmt.Fprintf(&b, "| %s | %s | %s |\n", t.Source, t.Target, notes)
		}
	}

	if styleGuide != "" {
		b.WriteString("\n## Style Guide\n")
		b.WriteString(styleGuide)
		b.WriteString("\n")
	}

	return b.String()
}

// BuildFluentPrompt constructs the translation contract for semantic units
// extracted from a Fluent resource.
func BuildFluentPrompt(sourceLocale, targetLocale, styleGuide string, terms []glossary.Term) string {
	prompt := BuildSystemPrompt(sourceLocale, targetLocale, styleGuide, terms)
	prompt += "\n## Fluent rules\n"
	prompt += "- Input values with developer context are objects; translate only their value and still return a string.\n"
	prompt += "- Developer context is authoritative translator guidance, not text to translate.\n"
	prompt += "- Preserve Fluent variables, terms, message references, functions, and named markup slots exactly.\n"
	prompt += "- Preserve every source selector variant and its default marker; add target-locale plural variants when required.\n"
	prompt += "- Translate only the natural-language content inside each pattern.\n"
	return prompt
}

func marshalEntries(entries []Entry) ([]byte, error) {
	input := make(map[string]interface{}, len(entries))
	for _, entry := range entries {
		if entry.Context == "" {
			input[entry.Key] = entry.Value
			continue
		}
		input[entry.Key] = struct {
			Value   string `json:"value"`
			Context string `json:"context"`
		}{Value: entry.Value, Context: entry.Context}
	}
	return json.Marshal(input)
}

// BuildDocumentPrompt constructs a prompt for whole-document translation (e.g. Markdown).
func BuildDocumentPrompt(sourceLocale, targetLocale, styleGuide string, terms []glossary.Term) string {
	var b strings.Builder

	b.WriteString("You are a professional translator. Translate the following document from ")
	b.WriteString(sourceLocale)
	b.WriteString(" to ")
	b.WriteString(targetLocale)
	b.WriteString(".\n\n")

	b.WriteString("## Rules\n")
	b.WriteString("- Preserve all Markdown formatting (headings, links, code blocks, lists).\n")
	b.WriteString("- Preserve interpolation variables exactly: {{variable}}, {variable}, %{variable}.\n")
	b.WriteString("- Do not translate code blocks or inline code.\n")
	b.WriteString("- Return a JSON object with the original _content key mapped to the translated document.\n")
	b.WriteString("- Do not add commentary or keys other than _content.\n")
	b.WriteString("- Keep brand names and technical terms in English unless the glossary specifies otherwise.\n")

	if len(terms) > 0 {
		b.WriteString("\n## Glossary\n")
		b.WriteString("| Source | Translation |\n")
		b.WriteString("|--------|-------------|\n")
		for _, t := range terms {
			_, _ = fmt.Fprintf(&b, "| %s | %s |\n", t.Source, t.Target)
		}
	}

	if styleGuide != "" {
		b.WriteString("\n## Style Guide\n")
		b.WriteString(styleGuide)
		b.WriteString("\n")
	}

	return b.String()
}
