package formats

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func FuzzJSONParseSerializeRoundTrip(f *testing.F) {
	for _, seed := range [][]byte{
		[]byte(`{"save":"Save"}`),
		[]byte(`{"items":["One","Two"],"enabled":true,"count":2,"empty":null}`),
		[]byte(`{"numeric":{"0":"Zero"}}`),
		[]byte(`{"nested":{"welcome":"Hello, {name}"}}`),
		[]byte(`not JSON`),
	} {
		f.Add(seed)
	}

	format := &JSONFormat{}
	f.Fuzz(func(t *testing.T, input []byte) {
		entries, err := format.Parse(input)
		if err != nil {
			return
		}
		output, err := format.Serialize(entries, input)
		if err != nil {
			t.Fatalf("serializing successfully parsed JSON: %v", err)
		}
		roundTripped, err := format.Parse(output)
		if err != nil {
			t.Fatalf("parsing serialized JSON: %v", err)
		}
		if !reflect.DeepEqual(roundTripped, entries) {
			t.Fatalf("round trip = %#v, want %#v", roundTripped, entries)
		}
	})
}

func BenchmarkJSONCatalogRoundTrip(b *testing.B) {
	var source bytes.Buffer
	source.WriteByte('{')
	for i := range 10_000 {
		if i > 0 {
			source.WriteByte(',')
		}
		_, _ = fmt.Fprintf(&source, `"key_%05d":"Value %05d"`, i, i)
	}
	source.WriteString(`,"enabled":true,"retries":3,"nullable":null}`)
	sourceData := source.Bytes()
	format := &JSONFormat{}

	b.ReportAllocs()
	b.SetBytes(int64(len(sourceData)))
	b.ResetTimer()
	for range b.N {
		entries, err := format.Parse(sourceData)
		if err != nil {
			b.Fatal(err)
		}
		if _, err := format.Serialize(entries, sourceData); err != nil {
			b.Fatal(err)
		}
	}
}
