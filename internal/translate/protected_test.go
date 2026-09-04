package translate

import "testing"

func TestValidateTranslationValuePreservesProtectedMarkdown(t *testing.T) {
	source := "Read [the guide](https://example.com/guide) and run `go test`.\n\n```go\nfmt.Println(\"ok\")\n```\n"
	valid := "Lisez [le guide](https://example.com/guide) et lancez `go test`.\n\n```go\nfmt.Println(\"ok\")\n```\n"
	if err := validateTranslationValue("_content", source, valid, "fr"); err != nil {
		t.Fatalf("valid protected Markdown was rejected: %v", err)
	}

	tests := map[string]string{
		"changed link":   "Lisez [le guide](https://example.com/autre) et lancez `go test`.\n\n```go\nfmt.Println(\"ok\")\n```\n",
		"changed inline": "Lisez [le guide](https://example.com/guide) et lancez `go test ./...`.\n\n```go\nfmt.Println(\"ok\")\n```\n",
		"changed fence":  "Lisez [le guide](https://example.com/guide) et lancez `go test`.\n\n```go\nfmt.Println(\"non\")\n```\n",
	}
	for name, target := range tests {
		t.Run(name, func(t *testing.T) {
			if err := validateTranslationValue("_content", source, target, "fr"); err == nil {
				t.Fatal("invalid protected Markdown was accepted")
			}
		})
	}
}

func TestValidateTranslationValueAllowsTargetOnlyICUPluralBranches(t *testing.T) {
	source := `{count, plural, one {{name} has one item} other {{name} has # items}}`
	target := `{count, plural, one {{name} имеет один предмет} few {{name} имеет # предмета} many {{name} имеет # предметов} other {{name} имеет # предмета}}`
	if err := validateTranslationValue("items", source, target, "ru"); err != nil {
		t.Fatalf("valid target-only ICU branches were rejected: %v", err)
	}
}
