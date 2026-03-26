package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
	"github.com/Tom-R-Main/Internationalizer/internal/glossary"
)

// Entry is a source key-value pair to be translated.
type Entry struct {
	Key   string
	Value string
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
		return NewOpenAI(apiKey, cfg.Model, baseURL), nil
	case "gemini":
		return NewGemini(apiKey, cfg.Model), nil
	case "openrouter":
		return NewOpenAI(apiKey, cfg.Model, "https://openrouter.ai/api/v1"), nil
	default:
		return nil, fmt.Errorf("unknown LLM provider: %s (supported: anthropic, openai, gemini, openrouter)", cfg.Provider)
	}
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
			b.WriteString(fmt.Sprintf("| %s | %s | %s |\n", t.Source, t.Target, notes))
		}
	}

	if styleGuide != "" {
		b.WriteString("\n## Style Guide\n")
		b.WriteString(styleGuide)
		b.WriteString("\n")
	}

	return b.String()
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
	b.WriteString("- Do not add commentary — output only the translated document.\n")
	b.WriteString("- Keep brand names and technical terms in English unless the glossary specifies otherwise.\n")

	if len(terms) > 0 {
		b.WriteString("\n## Glossary\n")
		b.WriteString("| Source | Translation |\n")
		b.WriteString("|--------|-------------|\n")
		for _, t := range terms {
			b.WriteString(fmt.Sprintf("| %s | %s |\n", t.Source, t.Target))
		}
	}

	if styleGuide != "" {
		b.WriteString("\n## Style Guide\n")
		b.WriteString(styleGuide)
		b.WriteString("\n")
	}

	return b.String()
}
