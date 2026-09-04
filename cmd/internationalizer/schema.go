package main

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/Tom-R-Main/Internationalizer/internal/onboarding"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const jsonSchemaDialect = "https://json-schema.org/draft/2020-12/schema"

// workflowOutputSchema describes JSON mode, not terminal text. Null data is
// permitted because argument/configuration failures happen before a result
// exists. Unknown object fields remain allowed for additive CLI evolution.
func workflowOutputSchema(path string) map[string]any {
	path = strings.TrimPrefix(path, "internationalizer ")
	var data map[string]any
	switch path {
	case "commands":
		data = schemaObject(map[string]any{
			"cli_version": schemaType("string"), "commands": schemaArray(schemaForType(reflect.TypeFor[commandContract]())),
			"total": schemaType("integer"), "matched": schemaType("integer"), "truncated": schemaType("boolean"),
			"exit_codes": schemaForType(reflect.TypeFor[map[string]string]()), "states": schemaArray(schemaType("string")),
			"authorization": schemaType("string"),
		}, "cli_version", "commands", "total", "matched", "truncated", "exit_codes", "states", "authorization")
	case "detect", "config check":
		data = schemaObject(map[string]any{
			"inspection": schemaForType(reflect.TypeFor[onboarding.Inspection]()),
			"total":      schemaForType(reflect.TypeFor[map[string]int]()), "matched": schemaForType(reflect.TypeFor[map[string]int]()),
			"offline": map[string]any{"type": "boolean", "const": true}, "provider_verified": map[string]any{"type": "boolean", "const": false},
		}, "inspection", "total", "matched", "offline", "provider_verified")
		data["description"] = "Offline discovery and configuration evidence; credential presence is not a successful provider call."
	case "config plan":
		data = schemaForType(reflect.TypeFor[onboarding.ConfigPlan]())
		data["description"] = "Reviewable proposal only. required_decisions must be resolved before application; a plan does not grant authorization."
	case "config apply":
		data = schemaForType(reflect.TypeFor[onboarding.ApplyReceipt]())
		data["description"] = "Receipt tied to a plan and verified configuration fingerprint; application does not translate catalogs."
	case "translate":
		data = schemaForType(reflect.TypeFor[translationJSON]())
		data["description"] = "Planning, generation, adoption, and partial failure are distinct. Dry-run makes no provider call or file change; generation is not human approval."
	case "validate":
		data = schemaForType(reflect.TypeFor[validationJSON]())
		data["description"] = "Counts cover the full selected scope before presentation limits. Structural validation and checked human approval are distinct."
	default:
		return map[string]any{"$schema": jsonSchemaDialect, "description": "This command does not declare a versioned workflow JSON result; inspect its command-specific help."}
	}
	envelope := schemaObject(map[string]any{
		"schema_version": map[string]any{"type": "integer", "const": 1},
		"status":         map[string]any{"type": "string", "description": "Workflow state; inspect errors and command-specific data. An exit code of zero does not imply configuration choices or human review are complete."},
		"data":           map[string]any{"anyOf": []any{data, schemaType("null")}},
		"errors":         schemaArray(schemaForType(reflect.TypeFor[jsonFailure]())),
	}, "schema_version", "status", "data", "errors")
	envelope["$schema"] = jsonSchemaDialect
	envelope["title"] = "Internationalizer v1: " + path
	return envelope
}

// workflowInputSchema is a normalized invocation, not a claim that the CLI
// accepts JSON stdin. Each flags property maps to --<name>; array values are
// repeated flags. The separate command argv is supplied by commands --json.
func workflowInputSchema(cmd *cobra.Command) map[string]any {
	properties := map[string]any{}
	required := []string{}
	visit := func(flag *pflag.Flag) {
		if flag.Hidden {
			return
		}
		property := flagInputSchema(flag)
		properties[flag.Name] = property
		if len(flag.Annotations[cobra.BashCompOneRequiredFlag]) > 0 {
			required = append(required, flag.Name)
		}
	}
	cmd.InheritedFlags().VisitAll(visit)
	cmd.Flags().VisitAll(visit)
	path := strings.TrimPrefix(cmd.CommandPath(), cmd.Root().Name()+" ")
	if path == "config apply" && cmd.Flags().Lookup("plan") != nil {
		found := false
		for _, name := range required {
			found = found || name == "plan"
		}
		if !found {
			required = append(required, "plan")
		}
	}
	flags := schemaObject(properties, required...)
	flags["additionalProperties"] = false
	schema := schemaObject(map[string]any{
		"flags": flags,
		"args":  map[string]any{"type": "array", "items": schemaType("string")},
	}, "flags")
	schema["$schema"] = jsonSchemaDialect
	schema["description"] = "Normalized invocation: map flag names to --name arguments; repeat array-valued flags. This is not a JSON-stdin API. --no-input disables prompts, not authorization."
	switch path {
	case "commands", "detect", "config check", "config plan", "config apply", "translate", "validate":
		schema["properties"].(map[string]any)["args"].(map[string]any)["maxItems"] = 0
	}
	return schema
}

func flagInputSchema(flag *pflag.Flag) map[string]any {
	property := schemaType("string")
	switch flag.Value.Type() {
	case "bool":
		property = schemaType("boolean")
		if value, err := strconv.ParseBool(flag.DefValue); err == nil {
			property["default"] = value
		}
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "count":
		property = schemaType("integer")
		if value, err := strconv.ParseInt(flag.DefValue, 10, 64); err == nil {
			property["default"] = value
		}
	case "float32", "float64":
		property = schemaType("number")
		if value, err := strconv.ParseFloat(flag.DefValue, 64); err == nil {
			property["default"] = value
		}
	case "stringArray", "stringSlice":
		property = schemaArray(schemaType("string"))
	case "string":
		property["default"] = flag.DefValue
	}
	property["description"] = flag.Usage
	property["x-cli-flag"] = "--" + flag.Name
	property["x-cli-type"] = flag.Value.Type()
	if flag.Name == "limit" {
		property["minimum"] = 0
	}
	return property
}

func schemaType(kind string) map[string]any { return map[string]any{"type": kind} }
func schemaArray(items map[string]any) map[string]any {
	return map[string]any{"type": "array", "items": items}
}
func schemaObject(properties map[string]any, required ...string) map[string]any {
	if required == nil {
		required = []string{}
	}
	return map[string]any{"type": "object", "properties": properties, "required": required, "additionalProperties": true}
}

// Struct JSON tags are the authority for required keys. Nil slices/maps encode
// as null in the CLI today; schemas acknowledge that instead of promising [].
func schemaForType(t reflect.Type) map[string]any {
	switch t.Kind() {
	case reflect.Pointer:
		return map[string]any{"anyOf": []any{schemaForType(t.Elem()), schemaType("null")}}
	case reflect.Struct:
		properties := map[string]any{}
		required := []string{}
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if !field.IsExported() {
				continue
			}
			tag := strings.Split(field.Tag.Get("json"), ",")
			if tag[0] == "-" {
				continue
			}
			name := tag[0]
			if name == "" {
				name = field.Name
			}
			properties[name] = schemaForType(field.Type)
			optional := false
			for _, option := range tag[1:] {
				optional = optional || option == "omitempty" || option == "omitzero"
			}
			if !optional {
				required = append(required, name)
			}
		}
		return schemaObject(properties, required...)
	case reflect.Slice:
		return map[string]any{"type": []string{"array", "null"}, "items": schemaForType(t.Elem())}
	case reflect.Array:
		return schemaArray(schemaForType(t.Elem()))
	case reflect.Map:
		return map[string]any{"type": []string{"object", "null"}, "additionalProperties": schemaForType(t.Elem())}
	case reflect.String:
		return schemaType("string")
	case reflect.Bool:
		return schemaType("boolean")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return schemaType("integer")
	case reflect.Float32, reflect.Float64:
		return schemaType("number")
	default:
		return map[string]any{}
	}
}
