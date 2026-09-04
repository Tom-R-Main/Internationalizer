package message

import "testing"

func TestResolveSyntax(t *testing.T) {
	for _, test := range []struct {
		source, format string
		policy, want   Syntax
	}{
		{"{{amount, number}}", "json", Auto, Legacy},
		{"{n, plural, one {[Guide](https://example.com)} other {[Guide](https://example.com)}}", "json", Auto, ICU},
		{"{ $count ->\n [one] One { -brand }\n *[other] { $count }\n}", "fluent", Auto, Fluent},
		{"{ $count ->\n [one] One { -brand }\n *[other] { $count }\n}", "json", Auto, Fluent},
		{"<code>{.sift,.claude}</code>", "json", I18next, I18next},
		{"{", "json", ICU, ICU},
		{"# Hi {name}", "markdown", Auto, Legacy},
		{"{name}", "json", Plain, Plain},
	} {
		if got := ResolveSyntax(test.format, test.policy, test.source); got != test.want {
			t.Errorf("ResolveSyntax(%q,%s,%q)=%s want %s", test.format, test.policy, test.source, got, test.want)
		}
	}
}
