package llm

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

func TestParseTranslationResponseRejectsIntegrityLoss(t *testing.T) {
	for _, input := range []string{
		`{"a":"{{name}}","a":"lost"}`,
		`{"a.b":"one","a":{"b":"two"}}`,
		"```json\n{\"a.b\":\"one\",\"a\":{\"b\":\"two\"}}\n```",
		`{"a":"one","a":"two","extra":"` + "```json\n{\"safe\":\"value\"}\n```" + `"}`,
	} {
		_, err := ParseTranslationResponse(input)
		var coded interface{ JSONCode() string }
		if !errors.As(err, &coded) {
			t.Fatalf("input accepted or lost integrity error: %v", err)
		}
	}
}

func TestParseTranslationResponse_RawJSON(t *testing.T) {
	input := `{"common.save": "Enregistrer", "common.cancel": "Annuler"}`
	result, err := ParseTranslationResponse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["common.save"] != "Enregistrer" {
		t.Errorf("got %q, want %q", result["common.save"], "Enregistrer")
	}
	if result["common.cancel"] != "Annuler" {
		t.Errorf("got %q, want %q", result["common.cancel"], "Annuler")
	}
}

func TestParseTranslationResponse_CodeBlock(t *testing.T) {
	input := "Here are the translations:\n```json\n{\"hello\": \"Bonjour\"}\n```\n"
	result, err := ParseTranslationResponse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["hello"] != "Bonjour" {
		t.Errorf("got %q, want %q", result["hello"], "Bonjour")
	}
}

func TestParseTranslationResponse_NestedJSON(t *testing.T) {
	input := `{"common": {"save": "Guardar", "cancel": "Cancelar"}}`
	result, err := ParseTranslationResponse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["common.save"] != "Guardar" {
		t.Errorf("got %q, want %q", result["common.save"], "Guardar")
	}
}

func TestParseTranslationResponse_SurroundingText(t *testing.T) {
	input := `Sure, here are the translations: {"key": "value"} hope that helps!`
	result, err := ParseTranslationResponse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["key"] != "value" {
		t.Errorf("got %q, want %q", result["key"], "value")
	}
}

func TestParseTranslationResponse_Invalid(t *testing.T) {
	_, err := ParseTranslationResponse("this is not json at all")
	if err == nil {
		t.Error("expected error for non-JSON input")
	}
}

func TestParseTranslationResponse_RejectsNonStringLeaves(t *testing.T) {
	tests := map[string]string{
		"boolean":           `{"label":true}`,
		"number":            `{"label":42}`,
		"array":             `{"label":["Save"]}`,
		"null":              `{"label":null}`,
		"array of objects":  `{"label":[{"hidden":"value"}]}`,
		"root array":        `[{"hidden":"value"}]`,
		"root null":         `null`,
		"trailing document": `{"label":"one"} {"label":"two"}`,
		"trailing scalar":   `{"label":"one"} true`,
	}
	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := ParseTranslationResponse(input); err == nil {
				t.Fatalf("ParseTranslationResponse(%s) accepted a non-string translation", input)
			}
		})
	}
}

func FuzzParseTranslationResponseRoundTrip(f *testing.F) {
	for _, seed := range []string{
		`{"save":"Save"}`,
		`{"common":{"welcome":"Hello, {name}"}}`,
		"```json\n{\"save\":\"Save\"}\n```",
		`not JSON`,
		`{"value":null}`,
		`{"a.b":{"c":"X"},"a":{"b":{"d":"Y"}}}`,
		`{"":{"":"X"}}`,
		`{"a":"x","a":"y"}`,
		`{"a.b":"x","a":{"b":"y"}}`,
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		parsed, err := ParseTranslationResponse(input)
		if err != nil {
			return
		}
		canonical, err := json.Marshal(parsed)
		if err != nil {
			t.Fatalf("marshalling successful parse: %v", err)
		}
		roundTripped, err := ParseTranslationResponse(string(canonical))
		if err != nil {
			t.Fatalf("parsing canonical response: %v", err)
		}
		if !reflect.DeepEqual(roundTripped, parsed) {
			t.Fatalf("round trip = %#v, want %#v", roundTripped, parsed)
		}
	})
}
