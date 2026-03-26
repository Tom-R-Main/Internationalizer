package glossary

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Term is a glossary entry mapping a source term to its approved translation.
type Term struct {
	Source     string `json:"source"`
	Target     string `json:"target"`
	IgnoreCase bool   `json:"ignore_case,omitempty"`
	WholeWord  bool   `json:"whole_word,omitempty"`
}

// Load reads the glossary file for a locale from the given directory.
func Load(dir, locale string) ([]Term, error) {
	path := filepath.Join(dir, locale+".json")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading glossary %s: %w", path, err)
	}

	var terms []Term
	if err := json.Unmarshal(data, &terms); err != nil {
		return nil, fmt.Errorf("parsing glossary %s: %w", path, err)
	}
	return terms, nil
}

// Save writes glossary terms to the locale file.
func Save(dir, locale string, terms []Term) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating glossary directory: %w", err)
	}
	path := filepath.Join(dir, locale+".json")
	data, err := json.MarshalIndent(terms, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

// Add appends a term to the glossary for a locale.
func Add(dir, locale, source, target string) error {
	terms, err := Load(dir, locale)
	if err != nil {
		return err
	}
	// Check for duplicates.
	for i, t := range terms {
		if strings.EqualFold(t.Source, source) {
			terms[i].Target = target
			return Save(dir, locale, terms)
		}
	}
	terms = append(terms, Term{Source: source, Target: target})
	return Save(dir, locale, terms)
}

// Remove deletes a term from the glossary for a locale.
func Remove(dir, locale, source string) error {
	terms, err := Load(dir, locale)
	if err != nil {
		return err
	}
	filtered := terms[:0]
	for _, t := range terms {
		if !strings.EqualFold(t.Source, source) {
			filtered = append(filtered, t)
		}
	}
	if len(filtered) == len(terms) {
		return fmt.Errorf("term %q not found in glossary for %s", source, locale)
	}
	return Save(dir, locale, filtered)
}

// FormatForPrompt returns a markdown table of glossary terms for LLM injection.
func FormatForPrompt(terms []Term) string {
	if len(terms) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("| Source | Translation | Notes |\n")
	b.WriteString("|--------|-------------|-------|\n")
	for _, t := range terms {
		notes := ""
		if t.IgnoreCase {
			notes += "case-insensitive"
		}
		if t.WholeWord {
			if notes != "" {
				notes += ", "
			}
			notes += "whole word"
		}
		b.WriteString(fmt.Sprintf("| %s | %s | %s |\n", t.Source, t.Target, notes))
	}
	return b.String()
}
