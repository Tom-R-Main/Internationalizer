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
		Key:       "common.save",
		Source:    "Save",
		Target:    "Enregistrer",
		Locale:    "fr",
		Hash:      HashSource("Save"),
		Timestamp: time.Now(),
	}
	if err := memory.Add(rec); err != nil {
		t.Fatalf("Add: %v", err)
	}

	// Lookup should hit.
	target, ok := memory.Lookup("fr", "common.save", HashSource("Save"))
	if !ok {
		t.Fatal("Lookup miss after Add")
	}
	if target != "Enregistrer" {
		t.Errorf("got %q, want %q", target, "Enregistrer")
	}

	// Lookup with different hash should miss.
	_, ok = memory.Lookup("fr", "common.save", HashSource("Save changed"))
	if ok {
		t.Error("Lookup hit with wrong hash")
	}

	// Lookup different locale should miss.
	_, ok = memory.Lookup("de", "common.save", HashSource("Save"))
	if ok {
		t.Error("Lookup hit for wrong locale")
	}

	// Reload from disk.
	memory2, err := Load(path)
	if err != nil {
		t.Fatalf("Reload: %v", err)
	}
	target, ok = memory2.Lookup("fr", "common.save", HashSource("Save"))
	if !ok || target != "Enregistrer" {
		t.Error("Lookup miss after reload from disk")
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
