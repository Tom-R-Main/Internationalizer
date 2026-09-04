// Package state persists the content-addressed provenance used to decide
// whether an existing translation is current, stale, or manually edited.
package state

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Tom-R-Main/Internationalizer/internal/locale"
)

const (
	SchemaVersion       = 2
	legacySchemaVersion = 1
)

// ReviewStatus records human approval independently from translation origin.
// Provider, TM, adoption, and manual provenance answer how a value arrived;
// review status answers whether a person approved that exact artifact.
type ReviewStatus string

const (
	ReviewNeedsReview ReviewStatus = "needs_review"
	ReviewApproved    ReviewStatus = "approved"
)

// Manifest is the versioned on-disk translation state.
type Manifest struct {
	SchemaVersion int              `json:"schema_version"`
	Translations  map[string]Entry `json:"translations"`
}

// Entry records the inputs and output for one bundle key and target locale.
type Entry struct {
	Bundle        string       `json:"bundle"`
	Key           string       `json:"key"`
	Locale        string       `json:"locale"`
	SourceHash    string       `json:"source_hash"`
	PolicyHash    string       `json:"policy_hash"`
	GuideHash     string       `json:"guide_hash,omitempty"`
	GlossaryHash  string       `json:"glossary_hash,omitempty"`
	PromptVersion int          `json:"prompt_version,omitempty"`
	TargetHash    string       `json:"target_hash"`
	Origin        string       `json:"origin,omitempty"`
	Provider      string       `json:"provider,omitempty"`
	Model         string       `json:"model,omitempty"`
	ReviewStatus  ReviewStatus `json:"review_status"`
	ReviewedAt    *time.Time   `json:"reviewed_at,omitempty"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

// New returns an empty manifest using the current schema.
func New() *Manifest {
	return &Manifest{
		SchemaVersion: SchemaVersion,
		Translations:  make(map[string]Entry),
	}
}

// Load reads a manifest. A missing file is an empty manifest; malformed or
// unsupported state is an error because treating it as absent loses provenance.
func Load(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return New(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading manifest %s: %w", path, err)
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("parsing manifest %s: %w", path, err)
	}
	if manifest.SchemaVersion != SchemaVersion && manifest.SchemaVersion != legacySchemaVersion {
		return nil, fmt.Errorf("manifest %s uses schema version %d; supported version is %d", path, manifest.SchemaVersion, SchemaVersion)
	}
	legacy := manifest.SchemaVersion == legacySchemaVersion
	manifest.SchemaVersion = SchemaVersion
	if manifest.Translations == nil {
		manifest.Translations = make(map[string]Entry)
	}
	canonicalTranslations := make(map[string]Entry, len(manifest.Translations))
	for _, entry := range manifest.Translations {
		entry.Locale = canonicalLocaleOrOriginal(entry.Locale)
		if legacy || entry.ReviewStatus == "" {
			entry.ReviewStatus = ReviewNeedsReview
			entry.ReviewedAt = nil
		}
		if err := validateReviewState(entry); err != nil {
			return nil, fmt.Errorf("manifest %s contains invalid review state for %s/%s/%s: %w", path, entry.Bundle, entry.Key, entry.Locale, err)
		}
		identity := Identity(entry.Bundle, entry.Key, entry.Locale)
		if existing, duplicate := canonicalTranslations[identity]; duplicate && existing != entry {
			return nil, fmt.Errorf("manifest %s contains conflicting entries for %s/%s/%s", path, entry.Bundle, entry.Key, entry.Locale)
		}
		canonicalTranslations[identity] = entry
	}
	manifest.Translations = canonicalTranslations
	return &manifest, nil
}

// Save replaces the manifest only after the complete new file is durable.
func (m *Manifest) Save(path string) error {
	m.SchemaVersion = SchemaVersion
	for _, entry := range m.Translations {
		if err := validateReviewState(entry); err != nil {
			return fmt.Errorf("saving manifest %s with invalid review state for %s/%s/%s: %w", path, entry.Bundle, entry.Key, entry.Locale, err)
		}
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding manifest: %w", err)
	}
	return WriteFileAtomic(path, append(data, '\n'), 0o644)
}

// Get returns the recorded state for one logical translation identity.
func (m *Manifest) Get(bundle, key, locale string) (Entry, bool) {
	entry, ok := m.Translations[Identity(bundle, key, locale)]
	return entry, ok
}

// Set records one logical translation identity.
func (m *Manifest) Set(entry Entry) {
	if m.Translations == nil {
		m.Translations = make(map[string]Entry)
	}
	entry.Locale = canonicalLocaleOrOriginal(entry.Locale)
	if entry.ReviewStatus == "" {
		entry.ReviewStatus = ReviewNeedsReview
		entry.ReviewedAt = nil
	}
	m.Translations[Identity(entry.Bundle, entry.Key, entry.Locale)] = entry
}

// Approve marks the exact recorded artifact approved at reviewedAt.
func (m *Manifest) Approve(bundle, key, locale string, reviewedAt time.Time) (Entry, error) {
	entry, ok := m.Get(bundle, key, locale)
	if !ok {
		return Entry{}, fmt.Errorf("translation %s/%s/%s is not tracked", bundle, key, locale)
	}
	stamp := reviewedAt.UTC()
	entry.ReviewStatus = ReviewApproved
	entry.ReviewedAt = &stamp
	m.Set(entry)
	return entry, nil
}

func validateReviewState(entry Entry) error {
	switch entry.ReviewStatus {
	case ReviewNeedsReview:
		if entry.ReviewedAt != nil {
			return fmt.Errorf("needs_review entry has reviewed_at")
		}
	case ReviewApproved:
		if entry.ReviewedAt == nil {
			return fmt.Errorf("approved entry lacks reviewed_at")
		}
	default:
		return fmt.Errorf("unknown review_status %q", entry.ReviewStatus)
	}
	return nil
}

// Identity returns a stable full SHA-256 identifier for a logical translation.
func Identity(bundle, key, locale string) string {
	return hashBytes([]byte(bundle + "\x00" + key + "\x00" + canonicalLocaleOrOriginal(locale)))
}

func canonicalLocaleOrOriginal(value string) string {
	canonical, err := locale.Canonical(value)
	if err != nil {
		return value
	}
	return canonical
}

// SourceHash binds translation reuse to both content and source format.
func SourceHash(format, source string) string {
	canonical := fmt.Sprintf("1:%d:%s:%d:%s", len(format), format, len(source), source)
	return hashBytes([]byte(canonical))
}

// SourceUnitHash also binds translator context and adapter-owned structure to
// provenance. Flat legacy units retain their existing SourceHash identity.
func SourceUnitHash(format, value, context, structure string) string {
	if context == "" && structure == "" {
		return SourceHash(format, value)
	}
	canonical := fmt.Sprintf("2:%d:%s:%d:%s:%d:%s:%d:%s", len(format), format, len(value), value, len(context), context, len(structure), structure)
	return hashBytes([]byte(canonical))
}

// TargetHash records the exact translated value last applied or adopted.
func TargetHash(target string) string {
	return hashBytes([]byte(target))
}

// HashValue produces a deterministic SHA-256 hash of a JSON-serializable value.
func HashValue(value interface{}) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("encoding hash input: %w", err)
	}
	return hashBytes(data), nil
}

func hashBytes(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

// WriteFileAtomic writes data to a temporary sibling, syncs it, then renames it
// over path. The existing file remains untouched if preparation fails.
func WriteFileAtomic(path string, data []byte, perm os.FileMode) (err error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating directory %s: %w", dir, err)
	}
	tmp, err := os.CreateTemp(dir, ".internationalizer-*")
	if err != nil {
		return fmt.Errorf("creating temporary file for %s: %w", path, err)
	}
	tmpPath := tmp.Name()
	defer func() {
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
	}()

	if err := tmp.Chmod(perm); err != nil {
		return fmt.Errorf("setting permissions for %s: %w", path, err)
	}
	if _, err := tmp.Write(data); err != nil {
		return fmt.Errorf("writing temporary file for %s: %w", path, err)
	}
	if err := tmp.Sync(); err != nil {
		return fmt.Errorf("syncing temporary file for %s: %w", path, err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("closing temporary file for %s: %w", path, err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("replacing %s: %w", path, err)
	}
	return nil
}
