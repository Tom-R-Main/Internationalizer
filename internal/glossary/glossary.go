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
	Source      string   `json:"source"`
	Target      string   `json:"target"`
	Variants    []string `json:"variants,omitempty"`
	Enforcement string   `json:"enforcement,omitempty"`
	IgnoreCase  bool     `json:"ignore_case,omitempty"`
	WholeWord   bool     `json:"whole_word,omitempty"`
}

const (
	EnforcementError   = "error"
	EnforcementWarning = "warning"
)

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
	if err := validateTerms(terms); err != nil {
		return nil, fmt.Errorf("validating glossary %s: %w", path, err)
	}
	return terms, nil
}

// Save writes glossary terms to the locale file.
func Save(dir, locale string, terms []Term) error {
	path := filepath.Join(dir, locale+".json")
	if err := validateTerms(terms); err != nil {
		return fmt.Errorf("validating glossary %s: %w", path, err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating glossary directory: %w", err)
	}
	data, err := json.MarshalIndent(terms, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func validateTerms(terms []Term) error {
	for i, term := range terms {
		switch term.Enforcement {
		case "", EnforcementError, EnforcementWarning:
		default:
			return fmt.Errorf("term %d (%q) has invalid enforcement %q", i, term.Source, term.Enforcement)
		}
	}
	return nil
}

// SourceIdenticalExempt reports whether an otherwise source-identical complete
// value is explicitly approved by a glossary term with the same source and
// target. It never treats a term embedded within a larger value as an exemption.
func SourceIdenticalExempt(terms []Term, source, target string) bool {
	for _, term := range terms {
		equal := func(left, right string) bool {
			if term.IgnoreCase {
				return strings.EqualFold(left, right)
			}
			return left == right
		}
		if equal(source, target) && equal(term.Source, source) && equal(term.Target, target) {
			return true
		}
	}
	return false
}

// ApprovedTargets returns the primary target followed by its approved variants.
func ApprovedTargets(term Term) []string {
	targets := make([]string, 1, 1+len(term.Variants))
	targets[0] = term.Target
	return append(targets, term.Variants...)
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
		_, _ = fmt.Fprintf(&b, "| %s | %s | %s |\n", t.Source, t.Target, notes)
	}
	return b.String()
}
