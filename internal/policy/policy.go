// Package policy resolves the effective translation policy for a locale.
package policy

import (
	"github.com/Tom-R-Main/Internationalizer/internal/config"
	"github.com/Tom-R-Main/Internationalizer/internal/glossary"
	"github.com/Tom-R-Main/Internationalizer/internal/llm"
	"github.com/Tom-R-Main/Internationalizer/internal/locale"
	"github.com/Tom-R-Main/Internationalizer/internal/message"
	"github.com/Tom-R-Main/Internationalizer/internal/state"
)

const PromptContractVersion = 3

// Resolved is the resolved prompt, provider settings, and stable hash for
// one target locale and source format.
type Resolved struct {
	Prompt        string
	LLM           config.LLM
	Hash          string
	GuideHash     string
	GlossaryHash  string
	PromptVersion int
}

// Resolve builds the effective translation policy. Callers provide
// already-loaded style-guide and glossary content so policy resolution remains
// independent of filesystem layout.
func Resolve(cfg *config.Config, targetLocale, format, styleGuide string, terms []glossary.Term, bundleSyntax ...message.Syntax) (Resolved, error) {
	syntax := cfg.MessageSyntax
	if len(bundleSyntax) > 0 {
		syntax = bundleSyntax[0]
	}
	if syntax == "" {
		syntax = message.Auto
	}
	if err := message.ValidateSyntax(syntax); err != nil {
		return Resolved{}, err
	}
	sourceLocale := cfg.SourceLocale
	if sourceLocale == "" {
		sourceLocale = "en"
	}
	canonicalSource, err := locale.Canonical(sourceLocale)
	if err != nil {
		return Resolved{}, err
	}
	canonicalTarget, err := locale.Canonical(targetLocale)
	if err != nil {
		return Resolved{}, err
	}

	effectiveLLM := cfg.LLMForLocale(canonicalTarget)
	prompt := llm.BuildSystemPrompt(canonicalSource, canonicalTarget, styleGuide, terms, syntax)
	switch format {
	case "markdown":
		prompt = llm.BuildDocumentPrompt(canonicalSource, canonicalTarget, styleGuide, terms, syntax)
	case "fluent":
		prompt = llm.BuildFluentPrompt(canonicalSource, canonicalTarget, styleGuide, terms)
	}
	prompt += "\nPreserve the complete content of HTML code elements and Markdown code spans exactly."

	guideHash, err := state.HashValue(styleGuide)
	if err != nil {
		return Resolved{}, err
	}
	glossaryHash, err := state.HashValue(terms)
	if err != nil {
		return Resolved{}, err
	}
	hash, err := state.HashValue(struct {
		Syntax       message.Syntax `json:"message_syntax"`
		Version      int            `json:"version"`
		SourceLocale string         `json:"source_locale"`
		TargetLocale string         `json:"target_locale"`
		Format       string         `json:"format"`
		Provider     string         `json:"provider"`
		Model        string         `json:"model"`
		Reasoning    string         `json:"reasoning_effort"`
		GuideHash    string         `json:"guide_hash"`
		GlossaryHash string         `json:"glossary_hash"`
	}{syntax, PromptContractVersion, canonicalSource, canonicalTarget, format, effectiveLLM.Provider, effectiveLLM.Model, llm.EffectiveReasoningEffort(effectiveLLM), guideHash, glossaryHash})
	if err != nil {
		return Resolved{}, err
	}

	return Resolved{
		Prompt:        prompt,
		LLM:           effectiveLLM,
		Hash:          hash,
		GuideHash:     guideHash,
		GlossaryHash:  glossaryHash,
		PromptVersion: PromptContractVersion,
	}, nil
}
