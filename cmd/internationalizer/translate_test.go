package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestTranslateAdoptExistingDoesNotRequireProviderCredentials(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	targetPath := filepath.Join(dir, "fr.json")
	manifestPath := filepath.Join(dir, "state", "manifest.json")
	if err := os.WriteFile(sourcePath, []byte(`{"a":"Save"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(targetPath, []byte(`{"a":"Enregistrer"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(dir, ".internationalizer.yml")
	configData := []byte(
		"target_locales: [fr]\n" +
			"source_path: " + sourcePath + "\n" +
			"manifest_path: " + manifestPath + "\n" +
			"tm_path: " + filepath.Join(dir, "state", "tm.jsonl") + "\n" +
			"llm:\n  provider: openai\n  api_key_env: TEST_INTERNATIONALIZER_MISSING_KEY\n",
	)
	if err := os.WriteFile(configPath, configData, 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := newTranslateCmd()
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetArgs([]string{"--config", configPath, "--adopt-existing"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(manifestPath); err != nil {
		t.Fatalf("manifest was not created: %v", err)
	}
}
