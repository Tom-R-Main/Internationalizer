package styleguide

import (
	"os"
	"path/filepath"
	"strings"
)

// Load reads the shared conventions and per-locale style guide,
// concatenating them into a single string for LLM prompt injection.
// Returns empty string (not error) if files don't exist.
func Load(dir, locale string) (string, error) {
	var parts []string

	// Read shared conventions.
	conventions, err := readIfExists(filepath.Join(dir, "_conventions.md"))
	if err != nil {
		return "", err
	}
	if conventions != "" {
		parts = append(parts, conventions)
	}

	// Read per-locale guide.
	localeGuide, err := readIfExists(filepath.Join(dir, locale+".md"))
	if err != nil {
		return "", err
	}
	if localeGuide != "" {
		parts = append(parts, localeGuide)
	}

	return strings.Join(parts, "\n\n---\n\n"), nil
}

func readIfExists(path string) (string, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}
