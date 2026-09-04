package llm

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/Tom-R-Main/Internationalizer/internal/jsonintegrity"
)

// ParseTranslationResponse extracts a map[string]string from LLM response text.
// It handles raw JSON, JSON wrapped in markdown code blocks, and nested JSON.
func ParseTranslationResponse(text string) (map[string]string, error) {
	text = strings.TrimSpace(text)

	// Try raw JSON first.
	if result, err := tryParseJSON(text); err == nil {
		return result, nil
	} else if isTerminalParseError(err) {
		return nil, err
	}

	// Try extracting from markdown code block.
	if idx := strings.Index(text, "```"); idx >= 0 {
		start := strings.Index(text[idx+3:], "\n")
		if start >= 0 {
			rest := text[idx+3+start+1:]
			if end := strings.Index(rest, "```"); end >= 0 {
				if result, err := tryParseJSON(strings.TrimSpace(rest[:end])); err == nil {
					return result, nil
				} else if isTerminalParseError(err) {
					return nil, err
				}
			}
		}
	}

	// Try finding the first { to last } span.
	first := strings.Index(text, "{")
	last := strings.LastIndex(text, "}")
	if first >= 0 && last > first {
		if result, err := tryParseJSON(text[first : last+1]); err == nil {
			return result, nil
		} else if isTerminalParseError(err) {
			return nil, err
		}
	}

	return nil, fmt.Errorf("could not parse translation response as a JSON object of strings")
}

func isTerminalParseError(err error) bool {
	var integrity *jsonintegrity.Error
	var shape *responseShapeError
	return errors.As(err, &integrity) || errors.As(err, &shape)
}

type responseShapeError struct{ message string }

func (e *responseShapeError) Error() string { return e.message }

func tryParseJSON(text string) (map[string]string, error) {
	// Decode through interface values so null and other non-string leaves cannot
	// silently become empty strings during unmarshalling.
	raw, err := jsonintegrity.Decode([]byte(text))
	if err != nil {
		return nil, err
	}
	nested, ok := raw.(map[string]interface{})
	if !ok {
		return nil, &responseShapeError{message: "translation response must be a JSON object"}
	}

	result := make(map[string]string)
	if err := flattenResponse("", nested, result); err != nil {
		return nil, &responseShapeError{message: err.Error()}
	}
	return result, nil
}

func flattenResponse(prefix string, val map[string]interface{}, out map[string]string) error {
	keys := make([]string, 0, len(val))
	for key := range val {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		v := val[key]
		p := key
		if prefix != "" {
			p = prefix + "." + key
		}
		switch child := v.(type) {
		case string:
			out[p] = child
		case map[string]interface{}:
			if err := flattenResponse(p, child, out); err != nil {
				return err
			}
		default:
			return fmt.Errorf("translation %q must be a string, got %T", p, v)
		}
	}
	return nil
}
