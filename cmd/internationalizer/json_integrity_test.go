package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Tom-R-Main/Internationalizer/internal/jsonintegrity"
)

func TestIntegrityErrorJSONIncludesLocations(t *testing.T) {
	_, err := jsonintegrity.Decode([]byte(`{"a.b":"PRIVATE","a":{"b":"OTHER PRIVATE"}}`))
	if err == nil {
		t.Fatal("ambiguous input accepted")
	}
	cmd := newValidateCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	if emitJSON(cmd, "error", nil, fmt.Errorf("parsing source catalog.json: %w", err)) == nil {
		t.Fatal("failure swallowed")
	}
	var envelope struct {
		Errors []struct {
			Code    string
			Details map[string]string
		}
	}
	if err := json.Unmarshal(out.Bytes(), &envelope); err != nil {
		t.Fatal(err)
	}
	if len(envelope.Errors) != 1 || envelope.Errors[0].Code != "json_flattened_key_collision" {
		t.Fatalf("output = %s", out.String())
	}
	details := envelope.Errors[0].Details
	if details["path"] != "/a/b" || details["other_path"] != "/a.b" || details["key"] != "a.b" {
		t.Fatalf("details = %+v", details)
	}
	if bytes.Contains(out.Bytes(), []byte("PRIVATE")) {
		t.Fatal("catalog value leaked")
	}
}
