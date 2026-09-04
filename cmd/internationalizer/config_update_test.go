package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestConfigPlanRepeatedUpdatesPreserveUntouchedBundle(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	original := []byte(`# preserve project comment
source_locale: en
target_locales: [fr, ja]
message_syntax: plain
llm:
  provider: openai
  locale_overrides:
    ja:
      model: preserved-model
glossary_dir: custom/glossary
future_setting: retained
bundles:
  - id: marketing
    source: marketing.json
    target: old-marketing/{locale}.json
    format: json
    future_setting: marketing-value
  - id: web
    source: web.json
    target: old-web/{locale}.json
    message_syntax: i18next
  - id: mobile # untouched bundle comment
    source: mobile.json
    target: mobile/{locale}.json
    message_syntax: plain
    future_setting: mobile-value
`)
	if err := os.WriteFile(filepath.Join(root, ".internationalizer.yml"), original, 0600); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"marketing.json", "web.json", "mobile.json"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(`{"hello":"Hello"}`), 0600); err != nil {
			t.Fatal(err)
		}
	}
	command := newRootCmd()
	var out bytes.Buffer
	command.SetOut(&out)
	command.SetErr(&out)
	if err := execute(command, []string{"config", "plan", "--update-bundle", "web", "--update-bundle", "marketing", "--target", "marketing=new-marketing/{locale}.json", "--target", "web=new-web/{locale}.json", "--json"}); err != nil {
		t.Fatalf("repeated retarget: %v: %s", err, out.String())
	}
	var result struct {
		Status string `json:"status"`
		Data   struct {
			ProposedYAML string `json:"proposed_yaml"`
		} `json:"data"`
	}
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Status != "planned" {
		t.Fatalf("updates were not planned: %s", out.String())
	}
	var expected, actual map[string]any
	if err := yaml.Unmarshal(original, &expected); err != nil {
		t.Fatal(err)
	}
	if err := yaml.Unmarshal([]byte(result.Data.ProposedYAML), &actual); err != nil {
		t.Fatal(err)
	}
	for _, value := range expected["bundles"].([]any) {
		bundle := value.(map[string]any)
		switch bundle["id"] {
		case "marketing":
			bundle["target"] = "new-marketing/{locale}.json"
		case "web":
			bundle["target"] = "new-web/{locale}.json"
		}
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("updates changed settings beyond selected targets:\n%s", result.Data.ProposedYAML)
	}
	for _, comment := range []string{"# preserve project comment", "# untouched bundle comment"} {
		if !strings.Contains(result.Data.ProposedYAML, comment) {
			t.Errorf("lost comment %q", comment)
		}
	}
	current, err := os.ReadFile(filepath.Join(root, ".internationalizer.yml"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(original, current) {
		t.Fatal("planning changed the original config")
	}
}

func TestConfigPlanRetargetExistingBundle(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	for name, contents := range map[string]string{
		".internationalizer.yml": "source_locale: en\ntarget_locales: [fr]\nsource_path: en.json\nmessage_syntax: plain\n",
		"en.json":                `{"hello":"Hello"}`,
	} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(contents), 0600); err != nil {
			t.Fatal(err)
		}
	}
	command := newRootCmd()
	var out bytes.Buffer
	command.SetOut(&out)
	command.SetErr(&out)
	if err := execute(command, []string{"config", "plan", "--update-bundle", "default", "--target", "default=translations/{locale}.json", "--json"}); err != nil {
		t.Fatalf("retarget existing bundle: %v: %s", err, out.String())
	}
	if !strings.Contains(out.String(), "translations/{locale}.json") || !strings.Contains(out.String(), `"status": "planned"`) {
		t.Fatalf("missing proposed target: %s", out.String())
	}
}

func TestConfigPlanRejectsInvalidUpdateFlags(t *testing.T) {
	for _, args := range [][]string{
		{"--target", "default=safe/{locale}.json"},
		{"--update-bundle", "default"},
		{"--update-bundle", "default", "--update-bundle", "default", "--target", "default=safe/{locale}.json"},
		{"--update-bundle", "default", "--add-bundle", "default=en.json", "--target", "default=safe/{locale}.json"},
		{"--update-bundle", "", "--target", "default=safe/{locale}.json"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			command := newRootCmd()
			var out bytes.Buffer
			command.SetOut(&out)
			command.SetErr(&out)
			if err := execute(command, append([]string{"config", "plan", "--json"}, args...)); err == nil {
				t.Fatalf("invalid update accepted: %s", out.String())
			}
			if !strings.Contains(out.String(), `"status": "error"`) {
				t.Fatalf("missing JSON failure: %s", out.String())
			}
		})
	}
}
