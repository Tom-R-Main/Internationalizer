package main

import (
	"encoding/json"
	"fmt"
	"github.com/Tom-R-Main/Internationalizer/internal/config"
	"github.com/Tom-R-Main/Internationalizer/internal/message"
	"github.com/Tom-R-Main/Internationalizer/internal/onboarding"
	"github.com/spf13/cobra"
	"os"
	"sort"
	"strings"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "config", Short: "Inspect, plan, and explicitly apply project configuration"}
	f := &inspectionFlags{}
	check := &cobra.Command{Use: "check", Short: "Check resolved configuration offline without changing files", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error { return runInspection(cmd, f, true) }}
	f.bind(check)
	cmd.AddCommand(check, newConfigPlanCmd(), newConfigApplyCmd())
	return cmd
}

func assignments(values []string) (map[string]string, error) {
	out := map[string]string{}
	for _, value := range values {
		key, v, ok := strings.Cut(value, "=")
		if !ok || key == "" || v == "" {
			return nil, fmt.Errorf("expected ID=value, got %q", value)
		}
		if _, exists := out[key]; exists {
			return nil, fmt.Errorf("duplicate decision for %q", key)
		}
		out[key] = v
	}
	return out, nil
}

func newConfigPlanCmd() *cobra.Command {
	var path, out, sourceLocale string
	var additions, syntaxes, targets, confirm, locales []string
	var asJSON bool
	cmd := &cobra.Command{Use: "plan", Short: "Propose a reviewable config change; never apply it", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		adds, err := assignments(additions)
		if err != nil {
			return err
		}
		modes, err := assignments(syntaxes)
		if err != nil {
			return err
		}
		targetMap, err := assignments(targets)
		if err != nil {
			return err
		}
		opts := onboarding.PlanOptions{Syntax: map[string]message.Syntax{}, ConfirmSources: confirm, SourceLocale: sourceLocale, TargetLocales: locales}
		for id, mode := range modes {
			opts.Syntax[id] = message.Syntax(mode)
		}
		for id := range targetMap {
			if _, ok := adds[id]; !ok {
				return fmt.Errorf("--target %s requires --add-bundle %s=source", id, id)
			}
		}
		if len(adds) > 0 {
			report, scanErr := onboarding.Scan(".", path)
			if scanErr != nil {
				return scanErr
			}
			ids := make([]string, 0, len(adds))
			for id := range adds {
				ids = append(ids, id)
			}
			sort.Strings(ids)
			for _, id := range ids {
				source := adds[id]
				b := config.Bundle{ID: id, Source: source, Target: targetMap[id], MessageSyntax: opts.Syntax[id]}
				for _, c := range report.Candidates {
					if c.ID == source || c.Source == source {
						b.Source = c.Source
						b.Format = c.Format
						if b.Target == "" {
							b.Target = c.Target
						}
						break
					}
				}
				if b.Target == "" {
					return fmt.Errorf("source %q is not a discovered catalog; supply --target %s=path/{locale}.json", source, id)
				}
				opts.AddBundles = append(opts.AddBundles, b)
			}
		}
		plan, err := onboarding.BuildPlan(".", path, opts)
		if err != nil {
			return err
		}
		if out != "" {
			data, marshalErr := json.MarshalIndent(plan, "", "  ")
			if marshalErr != nil {
				return marshalErr
			}
			file, openErr := os.OpenFile(out, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
			if openErr != nil {
				return fmt.Errorf("save plan (will not overwrite): %w", openErr)
			}
			_, writeErr := file.Write(append(data, '\n'))
			closeErr := file.Close()
			if writeErr != nil {
				return writeErr
			}
			if closeErr != nil {
				return closeErr
			}
		}
		status := "planned"
		if len(plan.RequiredDecisions) > 0 {
			status = "needs_decision"
		}
		if asJSON {
			return emitJSON(cmd, status, plan, nil)
		}
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), plan.Diff); err != nil {
			return err
		}
		for _, d := range plan.RequiredDecisions {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", d.Code, d.Message); err != nil {
				return err
			}
		}
		if out != "" {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Saved plan to %s. Review it, then: config apply --plan %s --no-input\n", out, out); err != nil {
				return err
			}
		}
		return nil
	}}
	cmd.Flags().StringVar(&path, "config", "", "Configuration path")
	cmd.Flags().StringVar(&out, "out", "", "Save plan to a new file (never overwrites)")
	cmd.Flags().StringArrayVar(&additions, "add-bundle", nil, "Explicit bundle ID=discovered-source-path (repeatable)")
	cmd.Flags().StringArrayVar(&syntaxes, "syntax", nil, "Explicit bundle ID=plain|i18next|icu|auto (repeatable)")
	cmd.Flags().StringArrayVar(&targets, "target", nil, "Target override for added bundle ID=path/{locale}.json")
	cmd.Flags().StringArrayVar(&confirm, "confirm-source", nil, "Confirm authoritative source path, including tmp/ (repeatable)")
	cmd.Flags().StringVar(&sourceLocale, "source-locale", "", "Explicit source locale (existing setting preserved when omitted)")
	cmd.Flags().StringArrayVar(&locales, "locale", nil, "Explicit target locale set (repeatable; existing set preserved when omitted)")
	cmd.Flags().BoolVar(&asJSON, "json", false, "Emit versioned JSON")
	return cmd
}

func newConfigApplyCmd() *cobra.Command {
	var path string
	var asJSON, noInput bool
	cmd := &cobra.Command{Use: "apply", Short: "Apply an explicitly selected saved plan after drift checks", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		if path == "" {
			return fmt.Errorf("--plan is required; --no-input only disables prompts")
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var plan onboarding.ConfigPlan
		if err = json.Unmarshal(data, &plan); err != nil {
			return fmt.Errorf("invalid plan: %w", err)
		}
		receipt, err := onboarding.ApplyPlan(&plan)
		if asJSON {
			status := "error"
			if err == nil {
				status = receipt.Status
			}
			return emitJSON(cmd, status, receipt, err)
		}
		if err != nil {
			return err
		}
		return json.NewEncoder(cmd.OutOrStdout()).Encode(receipt)
	}}
	cmd.Flags().StringVar(&path, "plan", "", "Saved config plan to apply (explicit mutation request)")
	cmd.Flags().BoolVar(&noInput, "no-input", false, "Disable prompting; does not authorize any additional action")
	cmd.Flags().BoolVar(&asJSON, "json", false, "Emit versioned JSON receipt")
	return cmd
}
