package llm

import (
	"encoding/json"
	"reflect"
	"testing"
)

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
		"boolean": `{"label":true}`,
		"number":  `{"label":42}`,
		"array":   `{"label":["Save"]}`,
		"null":    `{"label":null}`,
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
