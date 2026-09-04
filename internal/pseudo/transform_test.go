package pseudo

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Tom-R-Main/Internationalizer/internal/fluentpattern"
	"github.com/Tom-R-Main/Internationalizer/internal/formats"
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

func TestTransformPreservesFluentSelectorStructure(t *testing.T) {
	source := `{ $count ->
    [one] One item for { -brand-short-name }
   *[other] { $count } items for { $user }
}`
	target, err := Transform(source, Accented)
	if err != nil {
		t.Fatal(err)
	}
	for _, syntax := range []string{"{ $count ->", "[one]", "*[other]", "{ -brand-short-name }", "{ $count }", "{ $user }"} {
		if !strings.Contains(target, syntax) {
			t.Fatalf("pseudo target lost %q: %s", syntax, target)
		}
	}
}

func TestPseudoOptionalFluentCorpus(t *testing.T) {
	root := os.Getenv("INTERNATIONALIZER_FLUENT_CORPUS")
	if root == "" {
		t.Skip("set INTERNATIONALIZER_FLUENT_CORPUS to a directory of .ftl resources")
	}
	files := 0
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		clean := filepath.ToSlash(path)
		if entry.IsDir() || filepath.Ext(path) != ".ftl" || strings.Contains(clean, "/test/") || strings.Contains(clean, "/tests/") {
			return nil
		}
		files++
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		units, err := formats.ParseUnits(&formats.FluentFormat{}, data)
		if err != nil {
			return err
		}
		for _, unit := range units {
			target, err := Transform(unit.Value, Accented)
			if err != nil {
				return err
			}
			_, _, preserved, err := fluentpattern.Compare(unit.Value, target)
			if err != nil || !preserved {
				return &fluentPseudoCorpusError{path: path, unit: unit.ID, err: err}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if files == 0 {
		t.Fatalf("no .ftl files found under %s", root)
	}
}

type fluentPseudoCorpusError struct {
	path string
	unit string
	err  error
}

func (err *fluentPseudoCorpusError) Error() string {
	detail := "structure mismatch"
	if err.err != nil {
		detail = err.err.Error()
	}
	return err.path + ": " + err.unit + ": " + detail
}
