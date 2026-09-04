package translate

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Tom-R-Main/Internationalizer/internal/llm"
	"github.com/Tom-R-Main/Internationalizer/internal/message"
)

func TestI18nextProviderBoundary(t *testing.T) {
	const code = `<code>&lt;root&gt;/{.sift,.claude,.codex,.agents}/skills</code>`
	const target = "Lire " + code + " pour {{user.name}}"
	for _, candidate := range []string{target, strings.Replace(target, ".claude", ".changed", 1), strings.Replace(target, "{{user.name}}", "{{name}}", 1)} {
		t.Run(candidate, func(t *testing.T) {
			dir := t.TempDir()
			sourcePath := filepath.Join(dir, "en.json")
			data, err := json.Marshal(map[string]string{"docs": "Read " + code + " for {{user.name}}"})
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(sourcePath, data, 0o644); err != nil {
				t.Fatal(err)
			}
			cfg := testConfig(dir, sourcePath)
			cfg.MessageSyntax = message.I18next
			provider := &fakeProvider{response: &llm.TranslateResponse{Translations: map[string]string{"docs": candidate}}}
			results, err := Run(context.Background(), cfg, provider, Options{})
			if (err == nil) != (candidate == target) {
				t.Fatalf("error=%v results=%+v", err, results)
			}
			if candidate != target {
				if _, err := os.Stat(filepath.Join(dir, "fr.json")); !os.IsNotExist(err) {
					t.Fatalf("invalid target written: %v", err)
				}
			} else {
				if !strings.Contains(provider.requests[0].SystemPrompt, "Message syntax: i18next") {
					t.Fatal("missing syntax prompt")
				}
				cfg.MessageSyntax = message.Plain
				results, err := Run(context.Background(), cfg, provider, Options{DryRun: true})
				if err != nil || results[0].KeysPolicyStale != 1 {
					t.Fatalf("syntax change not stale: %v %+v", err, results)
				}
				if err := os.Remove(filepath.Join(dir, "fr.json")); err != nil {
					t.Fatal(err)
				}
				results, err = Run(context.Background(), cfg, provider, Options{})
				if err != nil || provider.calls != 2 || results[0].KeysCached != 0 {
					t.Fatalf("reused incompatible TM: calls=%d results=%+v error=%v", provider.calls, results, err)
				}
			}
		})
	}
}
