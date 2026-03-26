package tm

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Record is a single translation memory entry.
type Record struct {
	Key       string    `json:"key"`
	Source    string    `json:"source"`
	Target    string    `json:"target"`
	Locale    string    `json:"locale"`
	Hash      string    `json:"hash"`
	Timestamp time.Time `json:"timestamp"`
}

// Stats holds translation memory statistics.
type Stats struct {
	TotalRecords int            `json:"total_records"`
	ByLocale     map[string]int `json:"by_locale"`
	FileSize     int64          `json:"file_size_bytes"`
}

// TM is an append-only JSONL translation memory.
type TM struct {
	path  string
	mu    sync.RWMutex
	index map[string]map[string]Record // locale -> key -> record
}

// Load reads a JSONL translation memory file into memory.
func Load(path string) (*TM, error) {
	t := &TM{
		path:  path,
		index: make(map[string]map[string]Record),
	}

	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return t, nil
	}
	if err != nil {
		return nil, fmt.Errorf("opening TM %s: %w", path, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var rec Record
		if err := json.Unmarshal(line, &rec); err != nil {
			continue // skip malformed lines
		}
		if t.index[rec.Locale] == nil {
			t.index[rec.Locale] = make(map[string]Record)
		}
		t.index[rec.Locale][rec.Key] = rec
	}
	return t, scanner.Err()
}

// Lookup checks if a cached translation exists for the given key and source hash.
func (t *TM) Lookup(locale, key, sourceHash string) (string, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	localeMap, ok := t.index[locale]
	if !ok {
		return "", false
	}
	rec, ok := localeMap[key]
	if !ok || rec.Hash != sourceHash {
		return "", false
	}
	return rec.Target, true
}

// Add appends a record to the JSONL file and updates the in-memory index.
func (t *TM) Add(rec Record) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(t.path), 0o755); err != nil {
		return fmt.Errorf("creating TM directory: %w", err)
	}

	f, err := os.OpenFile(t.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("opening TM for write: %w", err)
	}
	defer f.Close()

	line, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	if _, err := f.Write(append(line, '\n')); err != nil {
		return err
	}

	if t.index[rec.Locale] == nil {
		t.index[rec.Locale] = make(map[string]Record)
	}
	t.index[rec.Locale][rec.Key] = rec
	return nil
}

// AddBatch appends multiple records efficiently in a single file operation.
func (t *TM) AddBatch(records []Record) error {
	if len(records) == 0 {
		return nil
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(t.path), 0o755); err != nil {
		return fmt.Errorf("creating TM directory: %w", err)
	}

	f, err := os.OpenFile(t.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("opening TM for write: %w", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, rec := range records {
		line, err := json.Marshal(rec)
		if err != nil {
			return err
		}
		if _, err := w.Write(line); err != nil {
			return err
		}
		if err := w.WriteByte('\n'); err != nil {
			return err
		}

		if t.index[rec.Locale] == nil {
			t.index[rec.Locale] = make(map[string]Record)
		}
		t.index[rec.Locale][rec.Key] = rec
	}
	return w.Flush()
}

// Stats returns translation memory statistics.
func (t *TM) Stats() Stats {
	t.mu.RLock()
	defer t.mu.RUnlock()

	s := Stats{
		ByLocale: make(map[string]int),
	}
	for locale, keys := range t.index {
		s.ByLocale[locale] = len(keys)
		s.TotalRecords += len(keys)
	}

	if info, err := os.Stat(t.path); err == nil {
		s.FileSize = info.Size()
	}
	return s
}

// Clear truncates the TM file and resets the in-memory index.
func (t *TM) Clear() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.index = make(map[string]map[string]Record)
	if err := os.Truncate(t.path, 0); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// Export writes all records as a JSON array to the given writer.
func (t *TM) Export(w io.Writer) error {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var records []Record
	for _, keys := range t.index {
		for _, rec := range keys {
			records = append(records, rec)
		}
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(records)
}

// Compact rewrites the TM file keeping only the latest record per (locale, key).
func (t *TM) Compact() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	f, err := os.Create(t.path)
	if err != nil {
		return fmt.Errorf("creating TM for compact: %w", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, keys := range t.index {
		for _, rec := range keys {
			line, err := json.Marshal(rec)
			if err != nil {
				return err
			}
			if _, err := w.Write(line); err != nil {
				return err
			}
			if err := w.WriteByte('\n'); err != nil {
				return err
			}
		}
	}
	return w.Flush()
}

// HashSource returns the SHA-256 hash of a source string.
func HashSource(source string) string {
	h := sha256.Sum256([]byte(source))
	return fmt.Sprintf("%x", h)
}
