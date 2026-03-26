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
			locales, _ := cmd.Flags().GetStringSlice("locale")
			batchSize, _ := cmd.Flags().GetInt("batch-size")
			concurrency, _ := cmd.Flags().GetInt("concurrency")

			// For dry-run, skip API key validation.
			if !dryRun {
				if err := cfg.Validate(); err != nil {
					return err
				}
			}

			var provider llm.Provider
			if !dryRun {
				provider, err = llm.NewProvider(cfg.LLM, cfg.APIKey())
				if err != nil {
					return err
				}
			}

			start := time.Now()
			results, err := translate.Run(context.Background(), cfg, provider, translate.Options{
				DryRun:      dryRun,
				Locales:     locales,
				BatchSize:   batchSize,
				Concurrency: concurrency,
			})
			if err != nil {
				return err
			}

			fmt.Print(translate.FormatResults(results, time.Since(start)))
			return nil
		},
	}

	cmd.Flags().StringP("config", "c", "", "path to config file (default: .internationalizer.yml)")
	cmd.Flags().StringSliceP("locale", "l", nil, "target locale(s) to translate (default: all)")
	cmd.Flags().Bool("dry-run", false, "show what would be translated without calling the LLM")
	cmd.Flags().Int("batch-size", 0, "keys per LLM call (overrides config)")
	cmd.Flags().Int("concurrency", 0, "parallel LLM calls (overrides config)")

	return cmd
}
