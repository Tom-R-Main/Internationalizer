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
		Short: "Translate missing keys using an LLM",
		Long:  "Detect missing translation keys and generate translations via an LLM provider.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return err
			}

			dryRun, _ := cmd.Flags().GetBool("dry-run")
			adoptExisting, _ := cmd.Flags().GetBool("adopt-existing")
			refreshPolicy, _ := cmd.Flags().GetBool("refresh-policy")
			locales, _ := cmd.Flags().GetStringSlice("locale")
			batchSize, _ := cmd.Flags().GetInt("batch-size")
			concurrency, _ := cmd.Flags().GetInt("concurrency")

			if err := cfg.ValidateProject(); err != nil {
				return err
			}
			// Local inspection and explicit adoption do not need provider credentials.
			if !dryRun && !adoptExisting {
				if err := cfg.ValidateCredentialsForLocales(locales); err != nil {
					return err
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
	cmd.Flags().Bool("dry-run", false, "show what would be translated without calling the LLM")
	cmd.Flags().Bool("adopt-existing", false, "record existing translations as the provenance baseline without calling the LLM")
	cmd.Flags().Bool("refresh-policy", false, "retranslate entries made stale by prompt, style-guide, glossary, provider, or model changes")
	cmd.Flags().Int("batch-size", 0, "keys per LLM call (overrides config)")
	cmd.Flags().Int("concurrency", 0, "parallel LLM calls (overrides config)")
	cmd.MarkFlagsMutuallyExclusive("dry-run", "adopt-existing")
	cmd.MarkFlagsMutuallyExclusive("adopt-existing", "refresh-policy")

	return cmd
}
