package llm

import (
	"strings"
	"testing"
)

func TestBuildDocumentPromptKeepsProviderJSONContract(t *testing.T) {
	prompt := BuildDocumentPrompt("en", "fr", "", nil)
	if !strings.Contains(prompt, "_content") || !strings.Contains(prompt, "JSON object") {
		t.Fatalf("document prompt does not preserve the provider JSON contract: %q", prompt)
	}
}
