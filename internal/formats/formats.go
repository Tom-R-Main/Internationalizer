package formats

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Format defines how to parse and serialize a locale file format.
type Format interface {
	Name() string
	Extensions() []string
	// Parse reads raw file data and returns a flat key->value map.
	// Nested structures are flattened with dot notation (e.g., "common.save").
	Parse(data []byte) (map[string]string, error)
	// Serialize writes the flat key->value map back to the format,
	// using the original data to preserve structure and ordering.
	Serialize(entries map[string]string, original []byte) ([]byte, error)
}

var registry = []Format{
	&JSONFormat{},
	&YAMLFormat{},
	&MarkdownFormat{},
}

// FormatForFile returns the appropriate format handler for a file path.
func FormatForFile(filename string) (Format, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, f := range registry {
		for _, e := range f.Extensions() {
			if ext == e {
				return f, nil
			}
		}
	}
	return nil, fmt.Errorf("unsupported file format: %s", ext)
}

// AllFormats returns all registered format handlers.
func AllFormats() []Format {
	return registry
}
