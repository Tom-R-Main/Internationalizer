package validate

import (
	"strings"
	"testing"

	"github.com/Tom-R-Main/Internationalizer/internal/message"

	"gopkg.in/yaml.v3"
)

func TestI18nextLiteralShellBraces(t *testing.T) {
	const code = `<code>&lt;root&gt;/{.sift,.claude,.codex,.agents}/skills</code>`
	cfg := validationConfig(t, map[string]string{"docs": "Read " + code + " for {{user.name}}"}, map[string]string{"docs": "Lire " + code + " pour {{user.name}}"})
	if err := yaml.Unmarshal([]byte("message_syntax: i18next\n"), cfg); err != nil {
		t.Fatal(err)
	}
	reports, err := ValidateWithOptions(cfg, Options{Strict: true})
	if err != nil || HasFailures(reports) {
		t.Fatalf("literal shell braces rejected: err=%v reports=%+v", err, reports)
	}
}

func TestSyntaxValidationContracts(t *testing.T) {
	for _, test := range []struct {
		name, source, target string
		syntax               message.Syntax
		fail                 bool
	}{
		{"i18next literal braces", "Run {one,two}", "Faire {one,two}", message.I18next, false},
		{"plain literal braces", "Read {draft", "Lire {brouillon", message.Plain, false},
		{"plain braces are not variables", "Read {{draft}}", "Lire {{brouillon}}", message.Plain, false},
		{"nested variable lost", "Hello {{user.name}}", "Bonjour {{name}}", message.I18next, true},
		{"escaping changed", "Hello {{- name}}", "Bonjour {{name}}", message.I18next, true},
		{"formatter lost", "Price {{amount, number}}", "Prix {{amount}}", message.I18next, true},
		{"repeated variable lost", "{{name}} / {{name}}", "{{name}}", message.I18next, true},
		{"code changed", "Read <CODE class='x>y'>{foo,bar}</CODE>", "Lire <CODE class='x>y'>{foo,baz}</CODE>", message.I18next, true},
		{"ICU unshaped source", "broken {", "cassé", message.ICU, true},
		{"ICU unshaped target", "Hello", "Bonjour }", message.ICU, true},
		{"ICU valid", "Hi {name}", "Salut {name}", message.ICU, false},
		{"ICU malformed", "{n, plural, one {One}}", "{n, plural, one {Un}}", message.ICU, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			cfg := validationConfig(t, map[string]string{"key": test.source}, map[string]string{"key": test.target})
			cfg.MessageSyntax = test.syntax
			reports, err := Validate(cfg)
			if err != nil {
				t.Fatal(err)
			}
			if HasFailures(reports) != test.fail {
				t.Fatalf("failure=%v want=%v: %+v", HasFailures(reports), test.fail, reports)
			}
		})
	}
}

func TestI18nextSyntaxRequiresLocalePluralKeys(t *testing.T) {
	cfg := validationConfig(t, map[string]string{"items_one": "One item", "items_other": "{{count}} items"}, map[string]string{"items_one": "Un article", "items_other": "{{count}} articles"})
	cfg.MessageSyntax = message.I18next
	cfg.TargetLocales = []string{"ru"}
	reports, err := Validate(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.Join(reports[0].Missing, ","), "items_few") {
		t.Fatalf("missing plural forms not detected: %+v", reports)
	}
}
