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
