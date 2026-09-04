package main

import (
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
		Args:  cobra.NoArgs,
		Short: "Validate locale files against the source locale",
		Long: `Check target locale structure and interpolation against the source locale.

Use --strict to require translated values and enforce extra-key, protected
structure, glossary, and configured plural rules. Use --require-state to verify
that source, policy, and target content still match the translation manifest.
Use --require-approved to additionally require explicit human approval.`,
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
			if err := cfg.ValidateProject(); err != nil {
				return codedError("config_invalid", err)
			}
			bundles, _ := cmd.Flags().GetStringSlice("bundle")
			locales, _ := cmd.Flags().GetStringSlice("locale")
			if err := selectConfig(cfg, bundles, locales); err != nil {
				return err
			}

			strict, _ := cmd.Flags().GetBool("strict")
			requireState, _ := cmd.Flags().GetBool("require-state")
			requireApproved, _ := cmd.Flags().GetBool("require-approved")
			reports, err := validate.ValidateWithOptions(cfg, validate.Options{
				Strict:          strict,
				RequireState:    requireState,
				RequireApproved: requireApproved,
			})
			if err != nil {
				return err
			}

			asJSON, _ := cmd.Flags().GetBool("json")
			quiet, _ := cmd.Flags().GetBool("quiet")

			if asJSON {
				codes, _ := cmd.Flags().GetStringSlice("finding-code")
				data := boundedValidation(reports, codes, limit)
				status := "structural_validation_passed"
				var failure error
				if validate.HasFailures(reports) {
					status = "validation_failed"
					failure = errValidationFailed
				}
				data.HumanReviewChecked = requireApproved
				data.HumanReviewApproved = requireApproved && failure == nil
				return emitJSON(cmd, status, data, failure)
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
	cmd.Flags().StringSlice("bundle", nil, "bundle ID(s) to validate (default: all)")
	cmd.Flags().StringSliceP("locale", "l", nil, "target locale(s) to validate (default: all)")
	cmd.Flags().StringSlice("finding-code", nil, "filter JSON findings by stable code; does not change failure status")
	cmd.Flags().Int("limit", 100, "maximum JSON reports and detail items returned; 0 returns all")
	cmd.Flags().BoolP("quiet", "q", false, "exit code only, no output")
	cmd.Flags().Bool("strict", false, "fail on untranslated values and strict policy findings")
	cmd.Flags().Bool("require-state", false, "fail when translation manifest state is missing or stale")
	cmd.Flags().Bool("require-approved", false, "fail unless current translations have explicit human approval")

	return cmd
}

type validationJSON struct {
	ReportCount         int                          `json:"report_count"`
	FindingCount        int                          `json:"finding_count"`
	MatchingFindings    int                          `json:"matching_findings"`
	FindingCounts       map[validate.FindingCode]int `json:"finding_counts"`
	Reports             []validate.Report            `json:"reports"`
	Truncated           bool                         `json:"truncated"`
	HumanReviewApproved bool                         `json:"human_review_approved"`
	HumanReviewChecked  bool                         `json:"human_review_checked"`
}

// Limits are presentation-only: aggregate status always reflects every selected
// bundle and locale, including findings hidden by a filter or output bound.
func boundedValidation(reports []validate.Report, codes []string, limit int) validationJSON {
	data := validationJSON{ReportCount: len(reports), FindingCounts: map[validate.FindingCode]int{}, Reports: []validate.Report{}}
	remaining := limit
	for _, report := range reports {
		selected := []validate.Finding{}
		for _, finding := range report.Findings {
			data.FindingCount++
			data.FindingCounts[finding.Code]++
			match := len(codes) == 0
			for _, code := range codes {
				if code == string(finding.Code) {
					match = true
				}
			}
			if !match {
				continue
			}
			data.MatchingFindings++
			if limit == 0 || remaining > 0 {
				selected = append(selected, finding)
				remaining--
			} else {
				data.Truncated = true
			}
		}
		report.Findings = selected
		// Legacy detail arrays are bounded too, and omitted under a code filter
		// because they do not carry codes. Counts above retain the full findings.
		if len(codes) > 0 {
			report.Missing = nil
			report.Extra = nil
			report.Mismatches = nil
			report.Errors = nil
		} else {
			if limit > 0 {
				n := min(max(remaining, 0), len(report.Missing))
				data.Truncated = data.Truncated || n < len(report.Missing)
				report.Missing = report.Missing[:n]
				remaining -= n
				n = min(max(remaining, 0), len(report.Extra))
				data.Truncated = data.Truncated || n < len(report.Extra)
				report.Extra = report.Extra[:n]
				remaining -= n
				n = min(max(remaining, 0), len(report.Mismatches))
				data.Truncated = data.Truncated || n < len(report.Mismatches)
				report.Mismatches = report.Mismatches[:n]
				remaining -= n
				n = min(max(remaining, 0), len(report.Errors))
				data.Truncated = data.Truncated || n < len(report.Errors)
				report.Errors = report.Errors[:n]
				remaining -= n
			}
		}
		if limit == 0 || len(data.Reports) < limit {
			data.Reports = append(data.Reports, report)
		} else {
			data.Truncated = true
		}
	}
	return data
}
