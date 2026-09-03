package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateJSONReturnsFailureAfterWritingReport(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	if err := os.WriteFile(sourcePath, []byte(`{"a":"A"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(dir, ".internationalizer.yml")
	configData := []byte("target_locales: [fr]\nsource_path: " + sourcePath + "\n")
	if err := os.WriteFile(configPath, configData, 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := newValidateCmd()
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetArgs([]string{"--config", configPath, "--json"})
	err := cmd.Execute()
	if !errors.Is(err, errValidationFailed) {
		t.Fatalf("Execute error = %v, want %v", err, errValidationFailed)
	}
	if !strings.Contains(stdout.String(), `"missing": [`) {
		t.Fatalf("JSON report was not written before failure: %q", stdout.String())
	}
}

func TestValidateDefaultAllowsExtraKey(t *testing.T) {
	configPath := writeValidateProject(t, `{"a":"A"}`, `{"a":"Un A","extra":"Supplémentaire"}`, "")

	cmd := newValidateCmd()
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetArgs([]string{"--config", configPath})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error for default extra key: %v", err)
	}
}

func TestValidateStrictRejectsExtraKey(t *testing.T) {
	configPath := writeValidateProject(t, `{"a":"A"}`, `{"a":"Un A","extra":"Supplémentaire"}`, "")

	cmd := newValidateCmd()
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetArgs([]string{"--config", configPath, "--strict"})
	if err := cmd.Execute(); !errors.Is(err, errValidationFailed) {
		t.Fatalf("Execute error = %v, want %v", err, errValidationFailed)
	}
}

func TestValidateRequireStateRejectsMissingManifest(t *testing.T) {
	dir := t.TempDir()
	manifestPath := filepath.Join(dir, "missing.lock")
	configPath := writeValidateProject(t, `{"a":"A"}`, `{"a":"Un A"}`, "manifest_path: "+manifestPath+"\n")

	cmd := newValidateCmd()
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetArgs([]string{"--config", configPath, "--require-state"})
	if err := cmd.Execute(); !errors.Is(err, errValidationFailed) {
		t.Fatalf("Execute error = %v, want %v", err, errValidationFailed)
	}
	if !strings.Contains(stdout.String(), "untracked") {
		t.Fatalf("require-state report lacks untracked finding: %q", stdout.String())
	}
}

func TestValidateStrictQuietEmitsNothing(t *testing.T) {
	configPath := writeValidateProject(t, `{"a":"A"}`, `{"a":"Un A","extra":"Supplémentaire"}`, "")

	cmd := newValidateCmd()
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"--config", configPath, "--strict", "--quiet"})
	if err := cmd.Execute(); !errors.Is(err, errValidationFailed) {
		t.Fatalf("Execute error = %v, want %v", err, errValidationFailed)
	}
	if stdout.Len() != 0 || stderr.Len() != 0 {
		t.Fatalf("quiet output: stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
}

func writeValidateProject(t *testing.T, source, target, extraConfig string) string {
	t.Helper()
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	targetPath := filepath.Join(dir, "fr.json")
	if err := os.WriteFile(sourcePath, []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(targetPath, []byte(target), 0o644); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(dir, ".internationalizer.yml")
	configData := "target_locales: [fr]\nsource_path: " + sourcePath + "\n" + extraConfig
	if err := os.WriteFile(configPath, []byte(configData), 0o644); err != nil {
		t.Fatal(err)
	}
	return configPath
}
