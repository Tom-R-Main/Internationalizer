package main

import (
	"encoding/json"
	"reflect"
	"slices"
	"testing"

	"github.com/spf13/cobra"
)

func TestWorkflowOutputSchemasDescribeActualData(t *testing.T) {
	for _, tc := range []struct {
		path     string
		required []string
	}{
		{"commands", []string{"commands", "cli_version", "exit_codes"}},
		{"detect", []string{"inspection", "total", "matched", "offline", "provider_verified"}},
		{"config check", []string{"inspection", "total", "matched"}},
		{"config plan", []string{"id", "config_path", "proposed_yaml", "required_decisions", "observations"}},
		{"config apply", []string{"plan_id", "status", "changed_paths", "config_sha256", "observations_revalidated"}},
		{"translate", []string{"summary", "jobs", "dry_run", "provider_called", "human_review_approved"}},
		{"validate", []string{"report_count", "finding_count", "reports", "human_review_checked"}},
	} {
		t.Run(tc.path, func(t *testing.T) {
			schema := workflowOutputSchema("internationalizer " + tc.path)
			if schema["$schema"] != jsonSchemaDialect || schema["additionalProperties"] != true {
				t.Fatalf("invalid schema metadata: %+v", schema)
			}
			if _, err := json.Marshal(schema); err != nil {
				t.Fatal(err)
			}
			properties := schema["properties"].(map[string]any)
			if properties["schema_version"].(map[string]any)["const"] != 1 {
				t.Fatal("missing schema version constraint")
			}
			alternatives := properties["data"].(map[string]any)["anyOf"].([]any)
			if alternatives[1].(map[string]any)["type"] != "null" {
				t.Fatal("early failures cannot return null data")
			}
			data := alternatives[0].(map[string]any)
			for _, name := range tc.required {
				if !slices.Contains(data["required"].([]string), name) {
					t.Errorf("missing required key %s", name)
				}
				if _, exists := data["properties"].(map[string]any)[name]; !exists {
					t.Errorf("missing property %s", name)
				}
			}
		})
	}
}

func TestWorkflowInputSchemaHasTypedFlags(t *testing.T) {
	root := newRootCmd()
	for _, path := range [][]string{{"translate"}, {"config", "plan"}, {"config", "apply"}} {
		cmd, _, err := root.Find(path)
		if err != nil {
			t.Fatal(err)
		}
		schema := workflowInputSchema(cmd)
		properties := schema["properties"].(map[string]any)
		flags := properties["flags"].(map[string]any)
		flagProps := flags["properties"].(map[string]any)
		if flagProps["json"].(map[string]any)["type"] != "boolean" || properties["args"].(map[string]any)["maxItems"] != 0 {
			t.Fatalf("incorrect flag/args shapes: %+v", schema)
		}
		switch cmd.Name() {
		case "translate":
			if flagProps["locale"].(map[string]any)["type"] != "array" || flagProps["limit"].(map[string]any)["type"] != "integer" {
				t.Fatal("wrong translation flag types")
			}
		case "plan":
			if flagProps["add-bundle"].(map[string]any)["type"] != "array" {
				t.Fatal("repeatable add-bundle flag is not an array")
			}
		case "apply":
			if !slices.Contains(flags["required"].([]string), "plan") {
				t.Fatal("apply does not require an explicit plan")
			}
		}
	}
}

func TestTranslationSchemaDistinguishesPersistence(t *testing.T) {
	schema := workflowOutputSchema("translate")
	data := schema["properties"].(map[string]any)["data"].(map[string]any)["anyOf"].([]any)[0].(map[string]any)
	properties := data["properties"].(map[string]any)
	summary := properties["summary"].(map[string]any)["properties"].(map[string]any)
	if _, ok := summary["persisted_jobs"]; !ok {
		t.Fatal("summary schema omits persisted jobs")
	}
	job := properties["jobs"].(map[string]any)["items"].(map[string]any)
	for _, name := range []string{"catalog_written", "manifest_updated"} {
		if !slices.Contains(job["required"].([]string), name) {
			t.Errorf("job schema omits required persistence flag %s", name)
		}
	}
}

func TestSchemaTracksOptionalAndNullableJSONFields(t *testing.T) {
	type payload struct {
		Name     string   `json:"name"`
		Optional string   `json:"optional,omitempty"`
		Items    []string `json:"items"`
	}
	schema := schemaForType(reflect.TypeFor[payload]())
	if !reflect.DeepEqual(schema["required"], []string{"name", "items"}) {
		t.Fatalf("required = %v", schema["required"])
	}
	items := schema["properties"].(map[string]any)["items"].(map[string]any)
	if !reflect.DeepEqual(items["type"], []string{"array", "null"}) {
		t.Fatalf("slice type = %v", items["type"])
	}
}

func TestInputSchemaUsesFlagDefaultsNotCurrentValues(t *testing.T) {
	cmd := &cobra.Command{Use: "example"}
	cmd.Flags().Int("count", 3, "count")
	if err := cmd.Flags().Set("count", "9"); err != nil {
		t.Fatal(err)
	}
	schema := workflowInputSchema(cmd)
	flags := schema["properties"].(map[string]any)["flags"].(map[string]any)["properties"].(map[string]any)
	if flags["count"].(map[string]any)["default"] != int64(3) {
		t.Fatal("schema default changed with invocation state")
	}
}
