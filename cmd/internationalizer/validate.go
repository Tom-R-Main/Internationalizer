package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
	"github.com/Tom-R-Main/Internationalizer/internal/validate"
	"github.com/spf13/cobra"
)

var errValidationFailed = errors.New("validation failed")

func newValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate locale files against the source locale",
		Long: `Check target locale structure and interpolation against the source locale.

Use --strict to require translated values and enforce extra-key, protected
structure, glossary, and configured plural rules. Use --require-state to verify
that source, policy, and target content still match the translation manifest.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return err
			}

			strict, _ := cmd.Flags().GetBool("strict")
			requireState, _ := cmd.Flags().GetBool("require-state")
			reports, err := validate.ValidateWithOptions(cfg, validate.Options{
				Strict:       strict,
				RequireState: requireState,
			})
			if err != nil {
				return err
			}

			asJSON, _ := cmd.Flags().GetBool("json")
			quiet, _ := cmd.Flags().GetBool("quiet")

			if asJSON {
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				if err := enc.Encode(reports); err != nil {
					return err
				}
			} else if !quiet {
				if _, err := fmt.Fprint(cmd.OutOrStdout(), validate.FormatHuman(reports)); err != nil {
					return err
				}
			}

			if validate.HasFailures(reports) {
				return errValidationFailed
			}
			return nil
		},
	}

	cmd.Flags().StringP("config", "c", "", "path to config file (default: .internationalizer.yml)")
	cmd.Flags().Bool("json", false, "output report as JSON")
	cmd.Flags().BoolP("quiet", "q", false, "exit code only, no output")
	cmd.Flags().Bool("strict", false, "fail on untranslated values and strict policy findings")
	cmd.Flags().Bool("require-state", false, "fail when translation manifest state is missing or stale")

	return cmd
}
