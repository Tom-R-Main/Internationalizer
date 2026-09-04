package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Tom-R-Main/Internationalizer/internal/translate"
	"github.com/Tom-R-Main/Internationalizer/internal/validate"
)

func runJSONCommand(t *testing.T, args ...string) (jsonEnvelope, error) {
	t.Helper()
	root := newRootCmd()
	var out, stderr bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&stderr)
	err := execute(root, args)
	var result jsonEnvelope
	if decodeErr := json.Unmarshal(out.Bytes(), &result); decodeErr != nil {
		t.Fatalf("invalid JSON: %v; stdout=%s; stderr=%s; error=%v", decodeErr, out.String(), stderr.String(), err)
	}
	if result.SchemaVersion != 1 {
		t.Fatalf("schema_version = %d", result.SchemaVersion)
	}
	if stderr.Len() != 0 {
		t.Fatalf("unexpected stderr: %s", stderr.String())
	}
	return result, err
}

func TestJSONEarlyFailuresAreVersioned(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []string
		code string
	}{
		{"unknown flag", []string{"translate", "--json", "--nonexistent"}, "invalid_arguments"},
		{"invalid limit", []string{"validate", "--json", "--limit=-1"}, "invalid_arguments"},
		{"missing config", []string{"validate", "--json", "--config", filepath.Join(t.TempDir(), "missing.yml")}, "config_invalid"},
		{"malformed flag", []string{"translate", "--json", "--limit=abc"}, "invalid_arguments"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result, err := runJSONCommand(t, tc.args...)
			if err == nil || result.Status != "error" || len(result.Errors) != 1 || result.Errors[0].Code != tc.code {
				t.Fatalf("result=%+v error=%v", result, err)
			}
			if len(result.Errors[0].Recovery) == 0 || len(result.Errors[0].Recovery[0].Argv) == 0 {
				t.Fatal("missing structured recovery")
			}
		})
	}
}

func TestJSONConfigurationParseErrorsDoNotEchoValues(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yml")
	const secret = "synthetic-private-value-never-echo"
	if err := os.WriteFile(path, []byte("target_locales: "+secret+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	result, err := runJSONCommand(t, "translate", "--dry-run", "--json", "--config", path)
	data, _ := json.Marshal(result)
	if err == nil || strings.Contains(string(data), secret) || strings.Contains(err.Error(), secret) {
		t.Fatalf("unsafe parse failure: %s", data)
	}
}

func TestJSONRecoveryNeverAutomaticallyReappliesStalePlan(t *testing.T) {
	cmd := newValidateCmd()
	action := errorRecovery(cmd, "STALE_PLAN")
	if strings.Join(action.Argv, " ") != "internationalizer config plan --help" || len(action.RequiredDecisions) == 0 || len(action.SideEffects) != 0 {
		t.Fatalf("unsafe stale-plan recovery: %+v", action)
	}
}

func TestDryRunJSONHasPlannedNotGeneratedState(t *testing.T) {
	configPath := writeValidateProject(t, `{"hello":"Hello","goodbye":"Goodbye"}`, `{}`, "message_syntax: plain\n")
	dir := filepath.Dir(configPath)
	before, err := os.ReadFile(filepath.Join(dir, "fr.json"))
	if err != nil {
		t.Fatal(err)
	}
	result, err := runJSONCommand(t, "translate", "--config", configPath, "--dry-run", "--json", "--limit=1")
	if err != nil || result.Status != "planned" {
		t.Fatalf("result=%+v error=%v", result, err)
	}
	data, _ := json.Marshal(result.Data)
	var run translationJSON
	if err := json.Unmarshal(data, &run); err != nil {
		t.Fatal(err)
	}
	if !run.DryRun || run.ProviderCalled || run.HumanReviewApproved || run.Summary.PlannedKeys != 2 || run.Summary.GeneratedKeys != 0 {
		t.Fatalf("incorrect planning state: %+v", run)
	}
	after, err := os.ReadFile(filepath.Join(dir, "fr.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatal("dry-run changed target")
	}
}

func TestValidationFilteringNeverHidesFailureStatus(t *testing.T) {
	configPath := writeValidateProject(t, `{"hello":"Hello"}`, `{}`, "")
	result, err := runJSONCommand(t, "validate", "--config", configPath, "--json", "--finding-code=source_identical", "--limit=1")
	if !errors.Is(err, errValidationFailed) || result.Status != "validation_failed" {
		t.Fatalf("result=%+v error=%v", result, err)
	}
	data, _ := json.Marshal(result.Data)
	var report validationJSON
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatal(err)
	}
	if report.MatchingFindings != 0 || report.ReportCount != 1 || report.HumanReviewApproved {
		t.Fatalf("unexpected report %+v", report)
	}
}

func TestBoundedValidationCountsAllFindings(t *testing.T) {
	reports := []validate.Report{{Bundle: "web", Locale: "fr", Findings: []validate.Finding{{Code: validate.CodeMissingKey}, {Code: validate.CodeMissingKey}, {Code: validate.CodeProtectedStructureMismatch}}}}
	data := boundedValidation(reports, nil, 1)
	if data.FindingCount != 3 || data.MatchingFindings != 3 || data.FindingCounts[validate.CodeMissingKey] != 2 || !data.Truncated || len(data.Reports[0].Findings) != 1 {
		t.Fatalf("unexpected bounded report %+v", data)
	}
}

func TestProviderCallsAreNotPlannedBatches(t *testing.T) {
	if summaryHasProviderCalls([]translate.Result{{Batches: 3, DryRun: true}}) {
		t.Fatal("planned batch treated as executed provider call")
	}
	if !summaryHasProviderCalls([]translate.Result{{ProviderCalls: 1}}) {
		t.Fatal("actual provider call omitted")
	}
}

func TestJSONAutoDiagnosticIsActionable(t *testing.T) {
	configPath := writeValidateProject(t, `{"docs.tui.skills.desc":"Read <code>{.sift,.claude}/skills</code>"}`, `{}`, "message_syntax: auto\n")
	for _, command := range []string{"validate", "translate"} {
		args := []string{command, "--config", configPath, "--json"}
		if command == "translate" {
			args = append(args, "--dry-run")
		}
		result, err := runJSONCommand(t, args...)
		data, _ := json.Marshal(result.Data)
		if err == nil || !strings.Contains(string(data), "brace syntax inside HTML code") || !strings.Contains(string(data), "plain, i18next, or icu") {
			t.Fatalf("missing actionable ambiguity: %s error=%v", data, err)
		}
	}
}
