// Package review manages explicit human approval of current translations.
package review

import (
	"fmt"
	"sort"
	"time"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
	localeid "github.com/Tom-R-Main/Internationalizer/internal/locale"
	"github.com/Tom-R-Main/Internationalizer/internal/state"
	"github.com/Tom-R-Main/Internationalizer/internal/validate"
)

// Filter selects manifest entries for listing or approval.
type Filter struct {
	Locale string
	Bundle string
	Keys   []string
	Status state.ReviewStatus
	All    bool
}

// List returns deterministically ordered manifest entries matching filter.
func List(manifest *state.Manifest, filter Filter) ([]state.Entry, error) {
	canonicalLocale := ""
	var err error
	if filter.Locale != "" {
		canonicalLocale, err = localeid.Canonical(filter.Locale)
		if err != nil {
			return nil, err
		}
	}
	if filter.Status != "" && filter.Status != state.ReviewNeedsReview && filter.Status != state.ReviewApproved {
		return nil, fmt.Errorf("unsupported review status %q", filter.Status)
	}
	keySet := make(map[string]struct{}, len(filter.Keys))
	for _, key := range filter.Keys {
		keySet[key] = struct{}{}
	}

	entries := make([]state.Entry, 0, len(manifest.Translations))
	for _, entry := range manifest.Translations {
		if canonicalLocale != "" {
			entryLocale, err := localeid.Canonical(entry.Locale)
			if err != nil || entryLocale != canonicalLocale {
				continue
			}
		}
		if filter.Bundle != "" && entry.Bundle != filter.Bundle {
			continue
		}
		if filter.Status != "" && entry.ReviewStatus != filter.Status {
			continue
		}
		if len(keySet) > 0 {
			if _, ok := keySet[entry.Key]; !ok {
				continue
			}
		}
		entries = append(entries, entry)
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Locale != entries[j].Locale {
			return entries[i].Locale < entries[j].Locale
		}
		if entries[i].Bundle != entries[j].Bundle {
			return entries[i].Bundle < entries[j].Bundle
		}
		return entries[i].Key < entries[j].Key
	})
	return entries, nil
}

// Approve verifies selected entries against source, policy, and target state,
// then records explicit approval. It never approves a stale or invalid value.
func Approve(cfg *config.Config, filter Filter, reviewedAt time.Time) ([]state.Entry, error) {
	effectiveConfig := *cfg
	effectiveConfig.ApplyDefaults()
	cfg = &effectiveConfig
	if err := cfg.ValidateProject(); err != nil {
		return nil, err
	}
	if filter.Locale == "" {
		return nil, fmt.Errorf("locale is required")
	}
	configuredLocale, ok := cfg.ConfiguredTargetLocale(filter.Locale)
	if !ok {
		return nil, fmt.Errorf("locale %q is not in target_locales", filter.Locale)
	}
	if filter.All == (len(filter.Keys) > 0) {
		return nil, fmt.Errorf("choose exactly one of all entries or one or more keys")
	}
	if len(filter.Keys) > 0 && filter.Bundle == "" {
		return nil, fmt.Errorf("bundle is required when approving individual keys")
	}

	manifest, err := state.Load(cfg.ManifestPath)
	if err != nil {
		return nil, err
	}
	selection := filter
	selection.Locale = configuredLocale
	selection.Status = ""
	entries, err := List(manifest, selection)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("no tracked translations match the approval selection")
	}
	if len(filter.Keys) > 0 {
		found := make(map[string]struct{}, len(entries))
		for _, entry := range entries {
			found[entry.Key] = struct{}{}
		}
		for _, key := range filter.Keys {
			if _, ok := found[key]; !ok {
				return nil, fmt.Errorf("translation %s/%s/%s is not tracked", filter.Bundle, key, configuredLocale)
			}
		}
	}

	reports, err := validate.ValidateWithOptions(cfg, validate.Options{RequireState: true})
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if err := entryValidationError(entry, reports); err != nil {
			return nil, err
		}
	}

	approved := make([]state.Entry, 0, len(entries))
	for _, entry := range entries {
		updated, err := manifest.Approve(entry.Bundle, entry.Key, entry.Locale, reviewedAt)
		if err != nil {
			return nil, err
		}
		approved = append(approved, updated)
	}
	if err := manifest.Save(cfg.ManifestPath); err != nil {
		return nil, err
	}
	return approved, nil
}

func entryValidationError(entry state.Entry, reports []validate.Report) error {
	for _, report := range reports {
		reportLocale, reportLocaleErr := localeid.Canonical(report.Locale)
		entryLocale, entryLocaleErr := localeid.Canonical(entry.Locale)
		if report.Bundle != entry.Bundle || reportLocaleErr != nil || entryLocaleErr != nil || reportLocale != entryLocale {
			continue
		}
		if len(report.Errors) > 0 {
			return fmt.Errorf("cannot approve %s/%s/%s: %s", entry.Bundle, entry.Key, entry.Locale, report.Errors[0])
		}
		if report.BlockedBySource {
			return fmt.Errorf("cannot approve %s/%s/%s: source bundle %s is invalid", entry.Bundle, entry.Key, entry.Locale, report.SourcePath)
		}
		for _, missing := range report.Missing {
			if missing == entry.Key {
				return fmt.Errorf("cannot approve %s/%s/%s: target key is missing", entry.Bundle, entry.Key, entry.Locale)
			}
		}
		for _, mismatch := range report.Mismatches {
			if mismatch.Key == entry.Key {
				return fmt.Errorf("cannot approve %s/%s/%s: interpolation mismatch", entry.Bundle, entry.Key, entry.Locale)
			}
		}
		for _, finding := range report.Findings {
			if finding.Key == entry.Key && finding.Severity == validate.SeverityError {
				return fmt.Errorf("cannot approve %s/%s/%s: %s", entry.Bundle, entry.Key, entry.Locale, finding.Message)
			}
		}
		return nil
	}
	return fmt.Errorf("cannot approve %s/%s/%s: no configured bundle/locale report", entry.Bundle, entry.Key, entry.Locale)
}
