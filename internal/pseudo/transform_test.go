package pseudo

import (
	"strings"
	"testing"

	"github.com/Tom-R-Main/Internationalizer/internal/message"
	"github.com/Tom-R-Main/Internationalizer/internal/validate"
)

func TestAccentedPreservesProtectedSyntax(t *testing.T) {
	source := "Save {{name}} in <strong>Account</strong>; read [guide](https://example.com/a_(b)) and run `go test`."
	target, err := Transform(source, Accented)
	if err != nil {
		t.Fatal(err)
	}
	if target == source || !strings.HasPrefix(target, "[!! ") {
		t.Fatalf("target = %q", target)
	}
	if findings := validate.ProtectedFindings("message", source, target, "en-XA"); len(findings) != 0 {
		t.Fatalf("protected findings = %#v; target = %q", findings, target)
	}
}

func TestTransformsOnlyICULiteralText(t *testing.T) {
	source := `{count, plural, one {One <strong>item</strong> for {name}} other {# items for {name}}}`
	target, err := Transform(source, Accented)
	if err != nil {
		t.Fatal(err)
	}
	if issues := message.Compare(source, target, "en-XA"); len(issues) != 0 {
		t.Fatalf("ICU issues = %#v; target = %q", issues, target)
	}
	if findings := validate.ProtectedFindings("items", source, target, "en-XA"); len(findings) != 0 {
		t.Fatalf("protected findings = %#v; target = %q", findings, target)
	}
}

func TestBidiUsesIsolateAndPreservesInterpolation(t *testing.T) {
	target, err := Transform("Hello {name}", Bidi)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(target, "\u2067") || !strings.HasSuffix(target, "\u2069") || !strings.Contains(target, "{name}") {
		t.Fatalf("target = %q", target)
	}
}

func TestRejectsUnknownStrategy(t *testing.T) {
	if _, err := Transform("Hello", Strategy("unknown")); err == nil {
		t.Fatal("Transform accepted unknown strategy")
	}
}
