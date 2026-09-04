package jsonintegrity

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestDecodeIntegrityLocations(t *testing.T) {
	cases := []struct{ input, code, key, first, second string }{
		{`{"a.b":"one","a":{"b":"two"}}`, "json_flattened_key_collision", "a.b", "/a.b", "/a/b"},
		{`{"a":{"b":"two"},"a.b":"one"}`, "json_flattened_key_collision", "a.b", "/a/b", "/a.b"},
		{`{"a.0":{},"a":[{}]}`, "json_flattened_key_collision", "a.0", "/a.0", "/a/0"},
		{`{"a.b":null,"a":{"b":false}}`, "json_flattened_key_collision", "a.b", "/a.b", "/a/b"},
		{`{"a/b~":{"x":"one","\u0078":"two"}}`, "json_duplicate_member", "a/b~.x", "/a~1b~0/x", "/a~1b~0/x"},
		{`{"section":{"hello":"one","hello":"two"}}`, "json_duplicate_member", "section.hello", "/section/hello", "/section/hello"},
		{`{"items":[{"hello":"one","hello":"two"}]}`, "json_duplicate_member", "items.0.hello", "/items/0/hello", "/items/0/hello"},
		{`{"":"one","":"one"}`, "json_duplicate_member", "", "/", "/"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			_, err := Decode([]byte(tc.input))
			var got *Error
			if !errors.As(err, &got) {
				t.Fatalf("got %v, want integrity error", err)
			}
			if got.Code != tc.code || got.Key != tc.key || got.OtherPath != tc.first || got.Path != tc.second {
				t.Fatalf("got %#v, want %s %q %q %q", got, tc.code, tc.key, tc.first, tc.second)
			}
			if got.Offset <= got.OtherOffset {
				t.Fatalf("offsets do not distinguish source locations: %#v", got)
			}
			if strings.Contains(err.Error(), "one") || strings.Contains(err.Error(), "two") {
				t.Fatal("error leaked a value")
			}
		})
	}
}

func TestDecodeValidNeighboringIdentities(t *testing.T) {
	for _, input := range []string{
		`{"a.b":"dotted","a":{"c":"nested"}}`,
		`{"a":"flat","a.b":"also flat"}`,
		`{"":"empty",".x":"leading dot"}`,
		`{"":{"x":"under empty key"}}`,
		`{"numeric":{"0":"object"},"array":["array"],"metadata":{"count":9007199254740993,"ok":true,"absent":null}}`,
		`["root array",{"nested":"value"}]`, `"root string"`, `null`,
	} {
		if _, err := Decode([]byte(input)); err != nil {
			t.Fatalf("rejected valid input %s: %v", input, err)
		}
	}
}

func TestDecodeMalformedAndDepthBoundary(t *testing.T) {
	for _, input := range []string{"", `{`, `[`, `{"x":}`, `{"x":1,}`, `[1,]`, `true false`, `{} garbage`} {
		if _, err := Decode([]byte(input)); err == nil {
			t.Fatalf("accepted %q", input)
		}
	}
	input := strings.Repeat("[", MaxDepth) + "0" + strings.Repeat("]", MaxDepth)
	if _, err := Decode([]byte(input)); err != nil {
		t.Fatalf("rejected depth limit: %v", err)
	}
	if _, err := Decode([]byte("[" + input + "]")); err == nil {
		t.Fatal("accepted depth beyond limit")
	}
}

func FuzzDecodeDeterministicRoundTrip(f *testing.F) {
	for _, input := range []string{`{"a":"x"}`, `{"a":"x","a":"y"}`, `{"a.b":"x","a":{"b":"y"}}`, `{"a":9007199254740993,"items":["x",null]}`, `{"":"x"}`, `"root"`, `{`} {
		f.Add(input)
	}
	f.Fuzz(func(t *testing.T, input string) {
		if len(input) > 65536 {
			t.Skip()
		}
		value, err := Decode([]byte(input))
		second, secondErr := Decode([]byte(input))
		if err != nil {
			if secondErr == nil || secondErr.Error() != err.Error() {
				t.Fatalf("nondeterministic rejection: %v / %v", err, secondErr)
			}
			return
		}
		if secondErr != nil || !reflect.DeepEqual(value, second) {
			t.Fatalf("nondeterministic successful parse")
		}
		encoded, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		roundTrip, err := Decode(encoded)
		if err != nil || !reflect.DeepEqual(value, roundTrip) {
			t.Fatalf("roundtrip changed decoded content: %v", err)
		}
	})
}
