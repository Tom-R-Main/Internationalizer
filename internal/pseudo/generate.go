package pseudo

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
	"github.com/Tom-R-Main/Internationalizer/internal/formats"
	localeid "github.com/Tom-R-Main/Internationalizer/internal/locale"
	"github.com/Tom-R-Main/Internationalizer/internal/state"
)

// GenerateOptions configures pseudolocale generation.
type GenerateOptions struct {
	Strategy Strategy
	Locale   string
	Force    bool
	DryRun   bool
}

// GenerateResult describes one generated bundle artifact.
type GenerateResult struct {
	Bundle     string
	Locale     string
	TargetPath string
	Units      int
	Written    bool
}

// Generate creates deterministic pseudolocale artifacts without an LLM or TM.
func Generate(cfg *config.Config, opts GenerateOptions) ([]GenerateResult, error) {
	effectiveConfig := *cfg
	effectiveConfig.ApplyDefaults()
	cfg = &effectiveConfig
	if err := cfg.ValidateProject(); err != nil {
		return nil, err
	}
	if _, err := DefaultLocale(opts.Strategy); err != nil {
		return nil, err
	}
	locale := opts.Locale
	if locale == "" {
		var err error
		locale, err = DefaultLocale(opts.Strategy)
		if err != nil {
			return nil, err
		}
	}
	canonicalLocale, err := localeid.Canonical(locale)
	if err != nil {
		return nil, fmt.Errorf("invalid pseudo locale %q: %w", locale, err)
	}

	manifest, err := state.Load(cfg.ManifestPath)
	if err != nil {
		return nil, err
	}
	policyHash, err := state.HashValue(struct {
		Version  int      `json:"version"`
		Strategy Strategy `json:"strategy"`
	}{Version: 1, Strategy: opts.Strategy})
	if err != nil {
		return nil, err
	}

	results := make([]GenerateResult, 0, len(cfg.EffectiveBundles()))
	now := time.Now().UTC()
	for _, bundle := range cfg.EffectiveBundles() {
		format, err := formatForBundle(bundle)
		if err != nil {
			return nil, fmt.Errorf("bundle %q format: %w", bundle.ID, err)
		}
		sourceData, err := os.ReadFile(bundle.Source)
		if err != nil {
			return nil, fmt.Errorf("reading bundle %q source %s: %w", bundle.ID, bundle.Source, err)
		}
		sourceUnits, err := formats.ParseUnits(format, sourceData)
		if err != nil {
			return nil, fmt.Errorf("parsing bundle %q source: %w", bundle.ID, err)
		}
		targetPath, err := bundle.TargetPath(locale)
		if err != nil {
			return nil, err
		}
		if existing, readErr := os.ReadFile(targetPath); readErr == nil {
			if !opts.Force && !pseudoOwnsArtifact(manifest, bundle.ID, canonicalLocale, format, existing) {
				return nil, fmt.Errorf("refusing to overwrite %s without --force because it is not a tracked pseudo artifact", targetPath)
			}
		} else if !os.IsNotExist(readErr) {
			return nil, fmt.Errorf("reading pseudo target %s: %w", targetPath, readErr)
		}

		pseudoUnits := make([]formats.Unit, len(sourceUnits))
		for index, unit := range sourceUnits {
			transformed, err := Transform(unit.Value, opts.Strategy)
			if err != nil {
				return nil, fmt.Errorf("pseudolocalizing bundle %q unit %q: %w", bundle.ID, unit.ID, err)
			}
			unit.Value = transformed
			pseudoUnits[index] = unit
		}
		output, err := formats.SerializeUnits(format, pseudoUnits, sourceData)
		if err != nil {
			return nil, fmt.Errorf("serializing pseudo target %s: %w", targetPath, err)
		}
		output = appendOneNewline(output)
		result := GenerateResult{Bundle: bundle.ID, Locale: canonicalLocale, TargetPath: targetPath, Units: len(pseudoUnits), Written: !opts.DryRun}
		results = append(results, result)
		if opts.DryRun {
			continue
		}
		if err := state.WriteFileAtomic(targetPath, output, 0o644); err != nil {
			return nil, err
		}
		for index, sourceUnit := range sourceUnits {
			manifest.Set(state.Entry{
				Bundle:       bundle.ID,
				Key:          sourceUnit.ID,
				Locale:       canonicalLocale,
				SourceHash:   state.SourceHash(format.Name(), sourceUnit.Value),
				PolicyHash:   policyHash,
				TargetHash:   state.TargetHash(pseudoUnits[index].Value),
				Origin:       "pseudo",
				ReviewStatus: state.ReviewNeedsReview,
				UpdatedAt:    now,
			})
		}
	}
	if !opts.DryRun {
		if err := manifest.Save(cfg.ManifestPath); err != nil {
			return nil, err
		}
	}
	return results, nil
}

func formatForBundle(bundle config.Bundle) (formats.Format, error) {
	if bundle.Format != "" {
		return formats.FormatByName(bundle.Format)
	}
	return formats.FormatForFile(bundle.Source)
}

func pseudoOwnsArtifact(manifest *state.Manifest, bundle, locale string, format formats.Format, data []byte) bool {
	units, err := formats.ParseUnits(format, data)
	if err != nil || len(units) == 0 {
		return false
	}
	for _, unit := range units {
		entry, ok := manifest.Get(bundle, unit.ID, locale)
		if !ok || entry.Origin != "pseudo" || entry.TargetHash != state.TargetHash(unit.Value) {
			return false
		}
	}
	return true
}

func appendOneNewline(data []byte) []byte {
	return append(bytes.TrimRight(data, "\n"), '\n')
}
