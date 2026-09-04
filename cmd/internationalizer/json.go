package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
	localeid "github.com/Tom-R-Main/Internationalizer/internal/locale"
	"github.com/spf13/cobra"
)

type recoveryAction struct {
	Argv              []string `json:"argv"`
	SideEffects       []string `json:"side_effects"`
	RequiredDecisions []string `json:"required_decisions"`
}

type jsonFailure struct {
	Code     string           `json:"code"`
	Message  string           `json:"message"`
	Recovery []recoveryAction `json:"recovery"`
}

type jsonEnvelope struct {
	SchemaVersion int           `json:"schema_version"`
	Status        string        `json:"status"`
	Data          any           `json:"data"`
	Errors        []jsonFailure `json:"errors"`
}

type reportedError struct{ error }

func (e reportedError) Unwrap() error { return e.error }

type commandError struct {
	code  string
	cause error
}

func (e commandError) Error() string    { return e.cause.Error() }
func (e commandError) Unwrap() error    { return e.cause }
func (e commandError) JSONCode() string { return e.code }

func codedError(code string, err error) error { return commandError{code: code, cause: err} }

// yaml's type errors may quote source values. Configuration may contain
// credentials or other private values in unsupported fields; never echo them.
func safeErrorMessage(err error) string {
	message := err.Error()
	if strings.Contains(message, "yaml:") {
		return "configuration YAML could not be parsed; check YAML syntax and field types"
	}
	return message
}

func safeConfigLoadError(err error) error {
	return codedError("config_invalid", errors.New(safeErrorMessage(err)))
}

// emitJSON writes exactly one versioned result, including failures. The returned
// marker keeps the root boundary from writing a second failure envelope.
func emitJSON(cmd *cobra.Command, status string, data any, err error) error {
	envelope := jsonEnvelope{SchemaVersion: 1, Status: status, Data: data, Errors: []jsonFailure{}}
	if err != nil {
		if status == "ok" || status == "applied" || status == "planned" || status == "generated" || status == "adopted" {
			envelope.Status = "error"
		}
		code := "command_failed"
		var coded interface{ JSONCode() string }
		if errors.As(err, &coded) {
			code = coded.JSONCode()
		}
		if errors.Is(err, errValidationFailed) {
			code = "validation_failed"
		}
		envelope.Errors = append(envelope.Errors, jsonFailure{Code: code, Message: safeErrorMessage(err), Recovery: []recoveryAction{errorRecovery(cmd, code)}})
	}
	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "  ")
	if encodeErr := enc.Encode(envelope); encodeErr != nil {
		return reportedError{encodeErr}
	}
	if err != nil {
		return reportedError{err}
	}
	return nil
}

func errorRecovery(cmd *cobra.Command, code string) recoveryAction {
	action := recoveryAction{Argv: []string{"internationalizer", "config", "check", "--json"}, SideEffects: []string{}, RequiredDecisions: []string{}}
	withConfig := true
	switch strings.ToLower(code) {
	case "invalid_arguments", "invalid_plan":
		action.Argv = []string{"internationalizer", "commands", "--json"}
		withConfig = false
	case "stale_plan", "plan_stale", "decisions_required", "plan_tampered":
		action.Argv = []string{"internationalizer", "config", "plan", "--help"}
		withConfig = false
		action.RequiredDecisions = []string{"review_current_configuration_and_create_a_new_plan"}
	case "validation_failed":
		action.Argv = []string{"internationalizer", "validate", "--json", "--limit", "100"}
		action.RequiredDecisions = []string{"review_catalog_findings"}
		for _, name := range []string{"bundle", "locale", "finding-code"} {
			if cmd.Flags().Changed(name) {
				values, _ := cmd.Flags().GetStringSlice(name)
				for _, value := range values {
					action.Argv = append(action.Argv, "--"+name, value)
				}
			}
		}
		for _, name := range []string{"strict", "require-state", "require-approved"} {
			value, _ := cmd.Flags().GetBool(name)
			if value {
				action.Argv = append(action.Argv, "--"+name)
			}
		}
	case "credentials_missing":
		action.RequiredDecisions = []string{"configure_the_required_provider_credential_in_your_environment"}
	case "translation_failed":
		action.RequiredDecisions = []string{"inspect_failed_jobs_and_current_state_before_retrying_provider_requests"}
	}
	if withConfig {
		if flag := cmd.Flags().Lookup("config"); flag != nil && flag.Value.String() != "" {
			action.Argv = append(action.Argv, "--config", flag.Value.String())
		}
	}
	return action
}

func jsonRequested(args []string) bool {
	for _, arg := range args {
		if arg == "--" {
			break
		}
		if arg == "--json" || arg == "--json=true" {
			return true
		}
	}
	return false
}

func execute(root *cobra.Command, args []string) error {
	root.SetArgs(args)
	cmd, err := root.ExecuteC()
	if err == nil {
		return nil
	}
	var reported reportedError
	if errors.As(err, &reported) {
		return err
	}
	if jsonRequested(args) {
		if cmd == nil {
			cmd = root
		}
		if strings.Contains(err.Error(), "unknown flag") || strings.Contains(err.Error(), "invalid argument") || strings.Contains(err.Error(), "unknown command") {
			err = codedError("invalid_arguments", err)
		}
		return emitJSON(cmd, "error", nil, err)
	}
	return err
}

func selectConfig(cfg *config.Config, bundles, locales []string) error {
	if len(bundles) > 0 {
		selected := make([]config.Bundle, 0, len(bundles))
		for _, id := range bundles {
			found := false
			for _, bundle := range cfg.EffectiveBundles() {
				if bundle.ID == id {
					selected = append(selected, bundle)
					found = true
					break
				}
			}
			if !found {
				return codedError("unknown_bundle", fmt.Errorf("bundle %q is not configured", id))
			}
		}
		cfg.Bundles = selected
	}
	if len(locales) > 0 {
		selected := make([]string, 0, len(locales))
		for _, locale := range locales {
			canonical, ok := cfg.ConfiguredTargetLocale(locale)
			if !ok {
				return codedError("unknown_locale", fmt.Errorf("locale %q is not in target_locales", locale))
			}
			selected = append(selected, canonical)
		}
		cfg.TargetLocales = selected
		// Filter only the in-memory execution configuration; on-disk overrides
		// remain untouched. Validation rejects overrides for unselected locales.
		filtered := make(map[string]config.LLMOverride)
		for locale, override := range cfg.LLM.LocaleOverrides {
			canonical, _ := localeid.Canonical(locale)
			for _, selection := range selected {
				selectedCanonical, _ := localeid.Canonical(selection)
				if canonical == selectedCanonical {
					filtered[locale] = override
				}
			}
		}
		cfg.LLM.LocaleOverrides = filtered
	}
	return nil
}
