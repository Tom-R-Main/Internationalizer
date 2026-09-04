package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
	"github.com/Tom-R-Main/Internationalizer/internal/llm"
	"github.com/Tom-R-Main/Internationalizer/internal/translate"
	"github.com/spf13/cobra"
)

func newTranslateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "translate",
		Args:  cobra.NoArgs,
		Short: "Translate missing keys using an LLM",
		Long:  "Detect missing translation keys and generate translations via an LLM provider.",
		RunE: func(cmd *cobra.Command, args []string) error {
			limit, _ := cmd.Flags().GetInt("limit")
			if limit < 0 {
				return codedError("invalid_arguments", fmt.Errorf("limit must be non-negative"))
			}
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return safeConfigLoadError(err)
			}

			dryRun, _ := cmd.Flags().GetBool("dry-run")
			adoptExisting, _ := cmd.Flags().GetBool("adopt-existing")
			refreshPolicy, _ := cmd.Flags().GetBool("refresh-policy")
			locales, _ := cmd.Flags().GetStringSlice("locale")
			batchSize, _ := cmd.Flags().GetInt("batch-size")
			concurrency, _ := cmd.Flags().GetInt("concurrency")
			bundles, _ := cmd.Flags().GetStringSlice("bundle")
			if err := selectConfig(cfg, bundles, nil); err != nil {
				return err
			}

			if err := cfg.ValidateProject(); err != nil {
				return codedError("config_invalid", err)
			}
			// Local inspection and explicit adoption do not need provider credentials.
			if !dryRun && !adoptExisting {
				if err := cfg.ValidateCredentialsForLocales(locales); err != nil {
					return codedError("credentials_missing", err)
				}
			}

			var provider llm.Provider
			localeProviders := make(map[string]llm.Provider)
			if !dryRun && !adoptExisting {
				selectedLocales := locales
				if len(selectedLocales) == 0 {
					selectedLocales = cfg.TargetLocales
				}
				needsDefaultProvider := false
				for _, requested := range selectedLocales {
					locale, ok := cfg.ConfiguredTargetLocale(requested)
					if !ok {
						return fmt.Errorf("locale %q is not in target_locales", requested)
					}
					if !cfg.HasLLMOverrideForLocale(locale) {
						needsDefaultProvider = true
						break
					}
				}
				if needsDefaultProvider {
					provider, err = llm.NewProvider(cfg.LLM, cfg.APIKey())
					if err != nil {
						return err
					}
				}
				for _, requested := range selectedLocales {
					locale, ok := cfg.ConfiguredTargetLocale(requested)
					if !ok {
						return fmt.Errorf("locale %q is not in target_locales", requested)
					}
					if !cfg.HasLLMOverrideForLocale(locale) {
						continue
					}
					localeLLM := cfg.LLMForLocale(locale)
					localeProviders[locale], err = llm.NewProvider(localeLLM, cfg.APIKeyForLocale(locale))
					if err != nil {
						return fmt.Errorf("creating LLM provider for locale %s: %w", locale, err)
					}
				}
			}

			start := time.Now()
			results, err := translate.Run(context.Background(), cfg, provider, translate.Options{
				DryRun:          dryRun,
				AdoptExisting:   adoptExisting,
				RefreshPolicy:   refreshPolicy,
				Locales:         locales,
				BatchSize:       batchSize,
				Concurrency:     concurrency,
				LocaleProviders: localeProviders,
			})
			asJSON, _ := cmd.Flags().GetBool("json")
			if asJSON {
				status := "generated"
				if dryRun {
					status = "planned"
				} else if adoptExisting {
					status = "adopted"
				}
				if err != nil {
					status = "blocked"
					err = codedError("translation_failed", err)
				}
				summary := translationSummary{Jobs: len(results)}
				for _, result := range results {
					summary.ErrorCount += len(result.Errors)
					if result.BlockedBySource || len(result.Errors) > 0 {
						summary.BlockedJobs++
					}
					if dryRun {
						summary.PlannedKeys += result.KeysSkipped
					}
					summary.GeneratedKeys += result.KeysTranslated
					if result.CatalogWritten || result.ManifestUpdated {
						summary.PersistedJobs++
					}
				}
				if err != nil && summary.PersistedJobs > 0 {
					status = "partial_failure"
				}
				providerCalled := summaryHasProviderCalls(results)
				truncated := false
				if limit > 0 && len(results) > limit {
					results = results[:limit]
					truncated = true
				}
				remaining := limit
				for i := range results {
					if limit > 0 {
						n := min(max(remaining, 0), len(results[i].Errors))
						truncated = truncated || n < len(results[i].Errors)
						results[i].Errors = results[i].Errors[:n]
						remaining -= n
					}
				}
				return emitJSON(cmd, status, translationJSON{Summary: summary, Jobs: results, Returned: len(results), Truncated: truncated, DryRun: dryRun, ProviderCalled: providerCalled, HumanReviewApproved: false}, err)
			}
			if len(results) > 0 {
				if _, outputErr := fmt.Fprint(cmd.OutOrStdout(), translate.FormatResults(results, time.Since(start))); outputErr != nil {
					return outputErr
				}
			}
			return err
		},
	}

	cmd.Flags().StringP("config", "c", "", "path to config file (default: .internationalizer.yml)")
	cmd.Flags().StringSliceP("locale", "l", nil, "target locale(s) to translate (default: all)")
	cmd.Flags().StringSlice("bundle", nil, "bundle ID(s) to translate (default: all)")
	cmd.Flags().Bool("json", false, "output a versioned JSON result, including failures")
	cmd.Flags().Int("limit", 100, "maximum JSON jobs and error details returned; 0 returns all (does not limit execution)")
	cmd.Flags().Bool("dry-run", false, "show what would be translated without calling the LLM")
	cmd.Flags().Bool("adopt-existing", false, "record existing translations as the provenance baseline without calling the LLM")
	cmd.Flags().Bool("refresh-policy", false, "retranslate entries made stale by prompt, style-guide, glossary, provider, or model changes")
	cmd.Flags().Int("batch-size", 0, "keys per LLM call (overrides config)")
	cmd.Flags().Int("concurrency", 0, "parallel LLM calls (overrides config)")
	cmd.MarkFlagsMutuallyExclusive("dry-run", "adopt-existing")
	cmd.MarkFlagsMutuallyExclusive("adopt-existing", "refresh-policy")

	return cmd
}

type translationSummary struct {
	Jobs          int `json:"jobs"`
	BlockedJobs   int `json:"blocked_jobs"`
	PlannedKeys   int `json:"planned_keys"`
	GeneratedKeys int `json:"generated_keys"`
	PersistedJobs int `json:"persisted_jobs"`
	ErrorCount    int `json:"error_count"`
}

type translationJSON struct {
	Summary             translationSummary `json:"summary"`
	Jobs                []translate.Result `json:"jobs"`
	Returned            int                `json:"returned"`
	Truncated           bool               `json:"truncated"`
	DryRun              bool               `json:"dry_run"`
	ProviderCalled      bool               `json:"provider_called"`
	HumanReviewApproved bool               `json:"human_review_approved"`
}

func summaryHasProviderCalls(results []translate.Result) bool {
	for _, result := range results {
		if result.ProviderCalls > 0 {
			return true
		}
	}
	return false
}
