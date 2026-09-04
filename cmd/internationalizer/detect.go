package main

import (
	"fmt"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
	"github.com/Tom-R-Main/Internationalizer/internal/onboarding"
	"github.com/spf13/cobra"
)

type inspectionFlags struct {
	config, bundle, locale, code string
	json                         bool
	limit                        int
}

type inspectionJSON struct {
	Total            map[string]int         `json:"total"`
	Matched          map[string]int         `json:"matched"`
	Offline          bool                   `json:"offline"`
	ProviderVerified bool                   `json:"provider_verified"`
	Inspection       *onboarding.Inspection `json:"inspection"`
}

func (f *inspectionFlags) bind(cmd *cobra.Command) {
	cmd.Flags().StringVar(&f.config, "config", "", "Configuration path")
	cmd.Flags().BoolVar(&f.json, "json", false, "Emit versioned JSON")
	cmd.Flags().StringVar(&f.bundle, "bundle", "", "Filter by bundle ID")
	cmd.Flags().StringVar(&f.locale, "locale", "", "Filter resolved targets by locale")
	cmd.Flags().StringVar(&f.code, "finding-code", "", "Filter diagnostics by stable code")
	cmd.Flags().IntVar(&f.limit, "limit", 50, "Maximum entries per result section (0 for all)")
}

func newDetectCmd() *cobra.Command {
	f := &inspectionFlags{}
	cmd := &cobra.Command{Use: "detect", Short: "Discover configured and uncovered catalogs with runtime evidence", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error { return runInspection(cmd, f, false) }}
	f.bind(cmd)
	return cmd
}

func runInspection(cmd *cobra.Command, f *inspectionFlags, check bool) error {
	if f.limit < 0 {
		return fmt.Errorf("limit must be nonnegative")
	}
	report, err := onboarding.Scan(".", f.config)
	if err != nil {
		return err
	}
	if f.bundle != "" {
		found := false
		for _, b := range report.Bundles {
			if b.ID == f.bundle {
				found = true
				break
			}
		}
		if !found {
			return codedError("unknown_bundle", fmt.Errorf("bundle %q is not configured", f.bundle))
		}
	}
	if f.locale != "" {
		cfg := config.Config{TargetLocales: report.TargetLocales}
		locale, ok := cfg.ConfiguredTargetLocale(f.locale)
		if !ok {
			return codedError("unknown_locale", fmt.Errorf("locale %q is not configured", f.locale))
		}
		f.locale = locale
	}
	totals := map[string]int{"candidates": len(report.Candidates), "bundles": len(report.Bundles), "diagnostics": len(report.Diagnostics), "credentials": len(report.Credentials)}
	failed := !report.ConfigExists
	for _, d := range report.Diagnostics {
		if d.Severity == "error" {
			failed = true
		}
	}
	if f.bundle != "" {
		bundles := report.Bundles[:0]
		for _, b := range report.Bundles {
			if b.ID == f.bundle {
				bundles = append(bundles, b)
			}
		}
		report.Bundles = bundles
		candidates := report.Candidates[:0]
		for _, c := range report.Candidates {
			for _, id := range c.ConfiguredBundles {
				if id == f.bundle {
					candidates = append(candidates, c)
					break
				}
			}
		}
		report.Candidates = candidates
	}
	if f.locale != "" {
		for i := range report.Bundles {
			b := &report.Bundles[i]
			target, ok := b.Targets[f.locale]
			b.Targets = map[string]string{}
			b.Locales = nil
			if ok {
				b.Targets[f.locale] = target
				b.Locales = []string{f.locale}
			}
		}
		credentials := report.Credentials[:0]
		for _, c := range report.Credentials {
			if c.Locale == f.locale {
				credentials = append(credentials, c)
			}
		}
		report.Credentials = credentials
	}
	diagnostics := report.Diagnostics[:0]
	for _, d := range report.Diagnostics {
		if (f.bundle == "" || d.Bundle == "" || d.Bundle == f.bundle) && (f.code == "" || d.Code == f.code) {
			diagnostics = append(diagnostics, d)
		}
	}
	report.Diagnostics = diagnostics
	matched := map[string]int{"candidates": len(report.Candidates), "bundles": len(report.Bundles), "diagnostics": len(report.Diagnostics), "credentials": len(report.Credentials)}
	if f.limit > 0 {
		if len(report.Bundles) > f.limit {
			report.Bundles = report.Bundles[:f.limit]
			report.Truncated = true
		}
		if len(report.Credentials) > f.limit {
			report.Credentials = report.Credentials[:f.limit]
			report.Truncated = true
		}
		if len(report.Candidates) > f.limit {
			report.Candidates = report.Candidates[:f.limit]
			report.Truncated = true
		}
		if len(report.Diagnostics) > f.limit {
			report.Diagnostics = report.Diagnostics[:f.limit]
			report.Truncated = true
		}
	}
	status := "ok"
	var resultErr error
	if check && failed {
		status = "error"
		resultErr = codedError("config_invalid", fmt.Errorf("configuration check failed"))
	}
	if f.json {
		return emitJSON(cmd, status, inspectionJSON{Total: totals, Matched: matched, Offline: true, ProviderVerified: false, Inspection: report}, resultErr)
	}
	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Config: %s (exists: %t)\n", report.ConfigPath, report.ConfigExists); err != nil {
		return err
	}
	for _, b := range report.Bundles {
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Bundle %s: %s -> %s; syntax=%s; framework=%s\n", b.ID, b.Source, b.Target, b.MessageSyntax, b.Framework); err != nil {
			return err
		}
	}
	for _, c := range report.Candidates {
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Catalog %s: framework=%s; suggested syntax=%s; configured=%v\n", c.Source, c.Framework, c.SuggestedSyntax, c.ConfiguredBundles); err != nil {
			return err
		}
	}
	for _, d := range report.Diagnostics {
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s [%s]: %s\n", d.Severity, d.Code, d.Message); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintln(cmd.OutOrStdout(), "Offline inspection only; no provider call or translation. Next: config plan --help"); err != nil {
		return err
	}
	return resultErr
}
