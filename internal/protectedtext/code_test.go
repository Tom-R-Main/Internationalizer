package protectedtext

import (
	"reflect"
	"testing"
)

func TestHTMLCodeSpans(t *testing.T) {
	for _, test := range []struct {
		source string
		want   []string
	}{
		{`Read <CODE title="a>b"><b>{foo,bar}</b></CODE> now`, []string{`<CODE title="a>b"><b>{foo,bar}</b></CODE>`}},
		{`<!-- <code>ignore</code> --> Text <code>keep</code>`, []string{`<code>keep</code>`}},
		{`<code>outer<code>inner</code>end</code>`, []string{`<code>outer<code>inner</code>end</code>`}},
		{`before <code>unclosed`, []string{`<code>unclosed`}},
		{`<codec>text</codec><code/>`, nil},
	} {
		if got := HTMLCode(test.source); !reflect.DeepEqual(got, test.want) {
			t.Errorf("HTMLCode(%q)=%q want %q", test.source, got, test.want)
		}
	}
}
