package tm

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTMRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tm.jsonl")

	memory, err := Load(path)
	if err != nil {
		t.Fatalf("Load empty: %v", err)
	}

	// Empty TM should have no records.
	s := memory.Stats()
	if s.TotalRecords != 0 {
		t.Errorf("empty TM has %d records", s.TotalRecords)
	}

	// Add a record.
	rec := Record{
		Bundle:     "app",
		Key:        "common.save",
		Source:     "Save",
		Target:     "Enregistrer",
		Locale:     "fr",
		Hash:       HashSource("Save"),
		PolicyHash: "policy-v1",
		Timestamp:  time.Now(),
	}
	if err := memory.Add(rec); err != nil {
		t.Fatalf("Add: %v", err)
	}

	// Lookup should hit.
	lookup, ok := memory.Lookup("fr", "app", "common.save", HashSource("Save"), "policy-v1")
	if !ok {
		t.Fatal("Lookup miss after Add")
	}
	if lookup.Target != "Enregistrer" {
		t.Errorf("got %q, want %q", lookup.Target, "Enregistrer")
	}

	// Lookup with different hash should miss.
	_, ok = memory.Lookup("fr", "app", "common.save", HashSource("Save changed"), "policy-v1")
	if ok {
		t.Error("Lookup hit with wrong hash")
	}

	// Lookup different locale should miss.
	_, ok = memory.Lookup("de", "app", "common.save", HashSource("Save"), "policy-v1")
	if ok {
		t.Error("Lookup hit for wrong locale")
	}

	// Reload from disk.
	memory2, err := Load(path)
	if err != nil {
		t.Fatalf("Reload: %v", err)
	}
	lookup, ok = memory2.Lookup("fr", "app", "common.save", HashSource("Save"), "policy-v1")
	if !ok || lookup.Target != "Enregistrer" {
		t.Error("Lookup miss after reload from disk")
	}
}

func TestTMLookupRequiresMatchingPolicy(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tm.jsonl")
	memory, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := memory.Add(Record{Bundle: "app", Key: "save", Source: "Save", Target: "Enregistrer", Locale: "fr", Hash: HashSource("Save"), PolicyHash: "policy-v1", Timestamp: time.Now()}); err != nil {
		t.Fatal(err)
	}
	if _, ok := memory.Lookup("fr", "app", "save", HashSource("Save"), "policy-v2"); ok {
		t.Fatal("Lookup reused a translation produced under a different policy")
	}
}

func TestLoadRejectsMalformedRecord(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tm.jsonl")
	if err := os.WriteFile(path, []byte("{\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("Load silently accepted malformed translation memory")
	}
}

func TestTMClear(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tm.jsonl")

	memory, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if err := memory.Add(Record{Key: "a", Source: "A", Target: "A-fr", Locale: "fr", Hash: HashSource("A"), Timestamp: time.Now()}); err != nil {
		t.Fatalf("Add: %v", err)
	}

	if err := memory.Clear(); err != nil {
		t.Fatalf("Clear: %v", err)
	}

	s := memory.Stats()
	if s.TotalRecords != 0 {
		t.Errorf("after clear: %d records", s.TotalRecords)
	}

	// File should be empty.
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if len(data) != 0 {
		t.Errorf("file not empty after clear: %d bytes", len(data))
	}
}

func TestTMAddBatch(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tm.jsonl")

	memory, _ := Load(path)
	records := []Record{
		{Key: "a", Source: "A", Target: "A-fr", Locale: "fr", Hash: HashSource("A"), Timestamp: time.Now()},
		{Key: "b", Source: "B", Target: "B-fr", Locale: "fr", Hash: HashSource("B"), Timestamp: time.Now()},
		{Key: "a", Source: "A", Target: "A-de", Locale: "de", Hash: HashSource("A"), Timestamp: time.Now()},
	}
	if err := memory.AddBatch(records); err != nil {
		t.Fatalf("AddBatch: %v", err)
	}

	s := memory.Stats()
	if s.TotalRecords != 3 {
		t.Errorf("got %d records, want 3", s.TotalRecords)
	}
	if s.ByLocale["fr"] != 2 {
		t.Errorf("fr has %d records, want 2", s.ByLocale["fr"])
	}
}
