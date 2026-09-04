// Package locale owns canonical locale identity and CLDR plural metadata.
package locale

import (
	"fmt"
	"slices"
	"strings"
	"sync"

	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
)

// PluralType selects cardinal or ordinal CLDR plural rules.
type PluralType string

const (
	Cardinal PluralType = "cardinal"
	Ordinal  PluralType = "ordinal"
)

var categoryCache sync.Map

// Canonical validates a BCP 47 language tag and returns its canonical form.
func Canonical(input string) (string, error) {
	tag, err := Parse(input)
	if err != nil {
		return "", err
	}
	return tag.String(), nil
}

// Parse validates a BCP 47 language tag. Underscore aliases are deliberately
// rejected so configuration and persisted state have one unambiguous syntax.
func Parse(input string) (language.Tag, error) {
	if input == "" {
		return language.Tag{}, fmt.Errorf("locale must not be empty")
	}
	if strings.ContainsRune(input, '_') {
		return language.Tag{}, fmt.Errorf("locale must use BCP 47 hyphen separators")
	}
	tag, err := language.Parse(input)
	if err != nil {
		return language.Tag{}, fmt.Errorf("invalid BCP 47 language tag: %w", err)
	}
	return tag, nil
}

// CardinalCategories returns the CLDR cardinal categories used by a locale.
func CardinalCategories(input string) ([]string, error) {
	return categories(input, Cardinal)
}

// OrdinalCategories returns the CLDR ordinal categories used by a locale.
func OrdinalCategories(input string) ([]string, error) {
	return categories(input, Ordinal)
}

func categories(input string, pluralType PluralType) ([]string, error) {
	tag, err := Parse(input)
	if err != nil {
		return nil, err
	}
	key := string(pluralType) + ":" + tag.String()
	if cached, ok := categoryCache.Load(key); ok {
		return slices.Clone(cached.([]string)), nil
	}

	rules := plural.Cardinal
	if pluralType == Ordinal {
		rules = plural.Ordinal
	}
	seen := make(map[plural.Form]struct{}, 6)
	record := func(form plural.Form) { seen[form] = struct{}{} }

	// x/text's generated CLDR evaluator differentiates integer rules through
	// mod 100 plus a few documented large-number exceptions. Cover both the
	// complete modulus and those exception ranges.
	for integer := 0; integer <= 2000; integer++ {
		record(rules.MatchPlural(tag, integer, 0, 0, 0, 0))
	}
	for _, integer := range []int{10_000, 100_000, 1_000_000, 2_000_000, 10_000_000} {
		record(rules.MatchPlural(tag, integer, 0, 0, 0, 0))
	}

	// Decimal rules can distinguish the number of visible digits and fraction
	// operands. Values through 199 cover both sides of x/text's bounded rule
	// representation and all modulo-100 classes.
	for visibleDigits := 1; visibleDigits <= 4; visibleDigits++ {
		maxFraction := 1
		for range visibleDigits {
			maxFraction *= 10
		}
		maxFraction--
		if maxFraction > 199 {
			maxFraction = 199
		}
		for integer := 0; integer <= 200; integer++ {
			for fraction := 0; fraction <= maxFraction; fraction++ {
				record(rules.MatchPlural(tag, integer, visibleDigits, visibleDigits, fraction, fraction))
			}
		}
	}

	// x/text/feature/plural still embeds CLDR 32 rules. Overlay category-set
	// changes needed by the current CLDR 49 contract; the public evaluator also
	// cannot represent the compact-exponent operand used by "many" rules.
	if pluralType == Cardinal && hasCompactExponentMany(tag) {
		seen[plural.Many] = struct{}{}
	}
	if pluralType == Cardinal && hasCurrentCardinalTwo(tag) {
		seen[plural.Two] = struct{}{}
	}
	if pluralType == Cardinal && needsCurrentCardinalOne(tag) {
		seen[plural.One] = struct{}{}
	}
	if pluralType == Cardinal && needsCurrentCardinalZeroFewMany(tag) {
		seen[plural.Zero] = struct{}{}
		seen[plural.Few] = struct{}{}
		seen[plural.Many] = struct{}{}
	}
	if pluralType == Ordinal && needsCurrentOrdinalMany(tag) {
		seen[plural.Many] = struct{}{}
	}
	if pluralType == Ordinal && needsCurrentOrdinalOne(tag) {
		seen[plural.One] = struct{}{}
	}
	if pluralType == Cardinal && hasRetiredCardinalMany(tag) {
		delete(seen, plural.Many)
	}

	ordered := make([]string, 0, len(seen))
	for _, candidate := range []struct {
		form plural.Form
		name string
	}{
		{plural.Zero, "zero"},
		{plural.One, "one"},
		{plural.Two, "two"},
		{plural.Few, "few"},
		{plural.Many, "many"},
		{plural.Other, "other"},
	} {
		if _, ok := seen[candidate.form]; ok {
			ordered = append(ordered, candidate.name)
		}
	}
	if !slices.Contains(ordered, "other") {
		ordered = append(ordered, "other")
	}
	categoryCache.Store(key, slices.Clone(ordered))
	return ordered, nil
}

func hasCompactExponentMany(tag language.Tag) bool {
	base, _ := tag.Base()
	switch base.String() {
	case "ca", "es", "fr", "gl", "it", "lld", "pt", "scn", "vec":
		return true
	default:
		return false
	}
}

func hasRetiredCardinalMany(tag language.Tag) bool {
	base, _ := tag.Base()
	// Modern CLDR removed Hebrew's former 20/30/... "many" category. Keep the
	// current category set even though x/text still evaluates the older rule.
	return base.String() == "he"
}

func hasCurrentCardinalTwo(tag language.Tag) bool {
	base, _ := tag.Base()
	// CLDR 42 split Maltese cardinal "other" to add a distinct "two" form.
	return base.String() == "mt"
}

func needsCurrentCardinalOne(tag language.Tag) bool {
	base, _ := tag.Base()
	// These languages postdate x/text's embedded plural tables.
	switch base.String() {
	case "bal", "lld", "scn", "vec":
		return true
	default:
		return false
	}
}

func needsCurrentCardinalZeroFewMany(tag language.Tag) bool {
	base, _ := tag.Base()
	// Modern Cornish distinguishes all six cardinal categories.
	return base.String() == "kw"
}

func needsCurrentOrdinalMany(tag language.Tag) bool {
	base, _ := tag.Base()
	// Sicilian and Venetian ordinal rules postdate x/text's embedded tables.
	switch base.String() {
	case "kw", "lld", "scn", "vec":
		return true
	default:
		return false
	}
}

func needsCurrentOrdinalOne(tag language.Tag) bool {
	base, _ := tag.Base()
	switch base.String() {
	case "bal", "kw":
		return true
	default:
		return false
	}
}
