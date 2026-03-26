package llm

import "testing"

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
