package config

import (
	"testing"

	"github.com/Tom-R-Main/Internationalizer/internal/message"
)

func TestMessageSyntaxInheritanceAndValidation(t *testing.T) {
	cfg := Config{SourceLocale: "en", TargetLocales: []string{"fr"}, MessageSyntax: message.I18next, Bundles: []Bundle{
		{ID: "web", Source: "web/en.json", Target: "web/{locale}.json"},
		{ID: "icu", Source: "icu/en.json", Target: "icu/{locale}.json", MessageSyntax: message.ICU},
		{ID: "fluent", Source: "en.ftl", Target: "{locale}.ftl", MessageSyntax: message.Auto},
	}}
	cfg.ApplyDefaults()
	if err := cfg.ValidateProject(); err != nil {
		t.Fatal(err)
	}
	got := cfg.EffectiveBundles()
	if got[0].MessageSyntax != message.I18next || got[1].MessageSyntax != message.ICU || got[2].MessageSyntax != message.Auto {
		t.Fatalf("inheritance: %+v", got)
	}
	if cfg.Bundles[0].MessageSyntax != "" {
		t.Fatal("resolution mutated input config")
	}
	if got[1].PluralStyle("i18next-v4") != "" {
		t.Fatal("ICU inherited i18next plural keys")
	}
	for _, syntax := range []message.Syntax{"typo", "fluent", message.Legacy} {
		cfg.MessageSyntax = syntax
		if err := cfg.ValidateProject(); err == nil {
			t.Fatalf("accepted invalid default %q", syntax)
		}
	}
	cfg.MessageSyntax = message.Auto
	cfg.Bundles[0].MessageSyntax = "typo"
	if err := cfg.ValidateProject(); err == nil {
		t.Fatal("accepted invalid override")
	}
	cfg.Bundles[0].MessageSyntax = message.Auto
	cfg.Bundles[2].MessageSyntax = message.I18next
	if err := cfg.ValidateProject(); err == nil {
		t.Fatal("accepted non-Fluent resource grammar")
	}
	cfg.Bundles[2] = Bundle{ID: "doc", Source: "README.md", Target: "{locale}.md", MessageSyntax: message.ICU}
	if err := cfg.ValidateProject(); err == nil {
		t.Fatal("accepted ICU document")
	}
	cfg.Bundles[2].Source = "README.MDX"
	if err := cfg.ValidateProject(); err == nil {
		t.Fatal("accepted ICU MDX document")
	}
	cfg.Bundles[2] = Bundle{ID: "fluent", Source: "EN.FTL", Target: "{locale}.ftl", MessageSyntax: message.Plain}
	if err := cfg.ValidateProject(); err == nil {
		t.Fatal("accepted plain Fluent resource with uppercase extension")
	}
}
