package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestInspectionRejectsUnknownScope(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	if err := os.WriteFile(filepath.Join(root, ".internationalizer.yml"), []byte("source_locale: en\ntarget_locales: [fr]\nsource_path: en.json\nmessage_syntax: plain\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "en.json"), []byte(`{"hello":"Hello"}`), 0600); err != nil {
		t.Fatal(err)
	}
	for _, flag := range []string{"--bundle", "--locale"} {
		t.Run(flag, func(t *testing.T) {
			command := newRootCmd()
			var out bytes.Buffer
			command.SetOut(&out)
			command.SetErr(&out)
			if err := execute(command, []string{"config", "check", flag, "typo", "--json"}); err == nil {
				t.Fatal("unknown selection passed")
			}
			var envelope jsonEnvelope
			if err := json.Unmarshal(out.Bytes(), &envelope); err != nil {
				t.Fatal(err)
			}
			if len(envelope.Errors) != 1 || envelope.Status != "error" {
				t.Fatalf("unexpected failure: %s", out.String())
			}
		})
	}
}

func TestCommandsContractFilter(t *testing.T) {
	command := newRootCmd()
	var out bytes.Buffer
	command.SetOut(&out)
	if err := execute(command, []string{"commands", "--command", "config apply", "--json"}); err != nil {
		t.Fatal(err)
	}
	var envelope struct {
		Data struct {
			Commands []commandContract `json:"commands"`
			Matched  int               `json:"matched"`
		}
	}
	if err := json.Unmarshal(out.Bytes(), &envelope); err != nil {
		t.Fatal(err)
	}
	if envelope.Data.Matched != 1 || len(envelope.Data.Commands) != 1 {
		t.Fatalf("bad filter: %s", out.String())
	}
	contract := envelope.Data.Commands[0]
	if contract.InputSchema == nil || contract.OutputSchema == nil || len(contract.SideEffects) == 0 {
		t.Fatalf("incomplete contract: %+v", contract)
	}
}
