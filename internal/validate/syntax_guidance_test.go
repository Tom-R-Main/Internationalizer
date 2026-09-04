package validate

import (
	"strings"
	"testing"

	"github.com/Tom-R-Main/Internationalizer/internal/message"
)

func TestSourceSyntaxGuidanceOnlyForAuto(t *testing.T) {
	source := `Read <code>{.sift,.claude}/skills</code>`
	for _, policy := range []message.Syntax{message.Auto, message.ICU} {
		findings := SyntaxSourceFindings("docs.tui.skills.desc", source, "en", message.ICU, policy)
		if len(findings) != 1 || findings[0].Code != CodeICUMessageSyntax {
			t.Fatalf("policy %s: %+v", policy, findings)
		}
		guidance := strings.Contains(findings[0].Message, "Select the bundle's runtime syntax")
		if guidance != (policy == message.Auto) {
			t.Fatalf("policy %s: %s", policy, findings[0].Message)
		}
	}
}
