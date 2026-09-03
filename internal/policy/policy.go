// Package policy resolves the effective translation policy for a locale.
package policy

import (
	"github.com/Tom-R-Main/Internationalizer/internal/config"
	"github.com/Tom-R-Main/Internationalizer/internal/glossary"
	"github.com/Tom-R-Main/Internationalizer/internal/llm"
	"github.com/Tom-R-Main/Internationalizer/internal/state"
)

const promptPolicyVersion = 1

// Resolved is the resolved prompt, provider settings, and stable hash for
// one target locale and source format.
type Resolved struct {
	Prompt string
	LLM    config.LLM
	Hash   string
}

// Resolve builds the effective translation policy. Callers provide
// already-loaded style-guide and glossary content so policy resolution remains
// independent of filesystem layout.
func Resolve(cfg *config.Config, targetLocale, format, styleGuide string, terms []glossary.Term) (Resolved, error) {
	effectiveLLM := cfg.LLMForLocale(targetLocale)
	prompt := llm.BuildSystemPrompt(cfg.SourceLocale, targetLocale, styleGuide, terms)
	if format == "markdown" {
		prompt = llm.BuildDocumentPrompt(cfg.SourceLocale, targetLocale, styleGuide, terms)
	}

	hash, err := state.HashValue(struct {
		Version      int    `json:"version"`
		SourceLocale string `json:"source_locale"`
		TargetLocale string `json:"target_locale"`
		Format       string `json:"format"`
		Provider     string `json:"provider"`
		Model        string `json:"model"`
		Reasoning    string `json:"reasoning_effort"`
		Prompt       string `json:"prompt"`
	}{promptPolicyVersion, cfg.SourceLocale, targetLocale, format, effectiveLLM.Provider, effectiveLLM.Model, llm.EffectiveReasoningEffort(effectiveLLM), prompt})
	if err != nil {
		return Resolved{}, err
	}

	return Resolved{Prompt: prompt, LLM: effectiveLLM, Hash: hash}, nil
}
